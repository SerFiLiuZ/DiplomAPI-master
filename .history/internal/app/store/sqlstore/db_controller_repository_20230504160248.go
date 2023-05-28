package sqlstore

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gopherschool/http-rest-api/internal/app/model"
)

type DBControllerRepository struct {
	store *Store
}

func (r *DBControllerRepository) UpdateDoneByTask(IDTask int) error {
	_, err := r.store.db.Exec(`
		UPDATE cards 
			SET done = 1 
			WHERE idCard = (
				SELECT idCard 
				FROM tasks 
				WHERE idTask = ?
			) AND (
				SELECT COUNT(*) 
				FROM tasks 
				WHERE idCard = (
					SELECT idCard 
					FROM tasks 
					WHERE idTask = ?
				) AND done = 0
			) = 0;
	`, IDTask, IDTask)
	log.Println(err)
	return err
}

func (r *DBControllerRepository) GetAppications(IDManager int) ([]*model.Appication, error) {
	rows, err := r.store.db.Query(`
		SELECT a.passwordCorp, a.FIO, a.phoneNumber, a.userName, a.dataApplication, a.chatID    
			FROM applications AS a
			JOIN users AS u 
			JOIN legalentity AS l ON l.idLegalEntity = u.idLegalEntity
			WHERE u.idUsers = ?;
    `, IDManager)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	applications := []*model.Appication{}

	for rows.Next() {
		application := &model.Appication{}
		err = rows.Scan(&application.PasswordCorp, &application.FIO, &application.PhoneNumber, &application.UserName, &application.DataApplication, &application.ChatID)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		applications = append(applications, application)
	}

	err = rows.Err()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return applications, nil

}

func (r *DBControllerRepository) AcceptApplication(chatIDUser int, IDManager int) error {
	rows, err := r.store.db.Query(`
		SELECT passwordCorp, FIO, phoneNumber, chatID, userName 
			FROM applications 
			WHERE chatID =?`,
		chatIDUser,
	)
	if err != nil {
		log.Println(err)
		return err
	}
	defer rows.Close()

	u := &model.User{}

	for rows.Next() {
		if err := rows.Scan(&u.LegalEntity.PasswordCorp, &u.FIO, &u.PhoneNumber, &u.TgData.ChatID, &u.TgData.UserNameTg); err != nil {
			log.Println(err)
			return err
		}
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return err
	}

	u_manager, err := r.store.User().Find(IDManager)
	if err != nil {
		log.Println(err)
		return err
	}

	var idAutorization int
	_, err = r.store.db.Exec(`
		INSERT INTO autorization (login, password) VALUES (?,?)`,
		"login."+u_manager.LegalEntity.NameCorp+"."+strconv.Itoa(*u.TgData.ChatID),
		"password."+u_manager.LegalEntity.NameCorp+"."+strconv.Itoa(*u.TgData.ChatID),
	)
	if err != nil {
		log.Println(err)
		return err
	}

	err = r.store.db.QueryRow(`SELECT idAutorization FROM autorization WHERE login=? AND password=?`,
		"login."+u_manager.LegalEntity.NameCorp+"."+strconv.Itoa(*u.TgData.ChatID),
		"password."+u_manager.LegalEntity.NameCorp+"."+strconv.Itoa(*u.TgData.ChatID),
	).Scan(&idAutorization)
	if err != nil {
		log.Println(err)
		return err
	}

	u.Autorization.ID = idAutorization

	var idtgdata int
	_, err = r.store.db.Exec("INSERT INTO tgdata (chatID, userName) VALUES (?,?)",
		u.TgData.ChatID,
		u.TgData.UserNameTg)
	if err != nil {
		log.Println(err)
		return err
	}

	err = r.store.db.QueryRow(`SELECT idtgdata FROM tgdata WHERE chatID=? AND userName=?`,
		u.TgData.ChatID,
		u.TgData.UserNameTg,
	).Scan(&idtgdata)
	if err != nil {
		log.Println(err)
		return err
	}

	u.TgData.ID = &idtgdata

	_, err = r.store.db.Exec(
		"INSERT INTO users (FIO, status, idAutorization, idLegalEntity, idTgData, assigned_to, phoneNumber) VALUES (?,?,?,?,?,?,?)",
		u.FIO,
		"worker",
		u.Autorization.ID,
		u_manager.LegalEntity.ID,
		u.TgData.ID,
		IDManager,
		u.PhoneNumber,
	)
	if err != nil {
		log.Println(err)
		return err
	}

	err = r.RejectApplication(chatIDUser)
	log.Println(err)
	return err
}

func (r *DBControllerRepository) RejectApplication(chatIDUser int) error {
	_, err := r.store.db.Exec(
		"DELETE FROM applications WHERE chatID=?",
		chatIDUser,
	)
	log.Println(err)
	return err
}

func (r *DBControllerRepository) SendMesage(TaskTitle string, TaskDes string, TaskDueDate string, TaskSelectedWorkers []*model.User) error {
	type_message := "sendMessage"
	bot_token := "6182612460:AAGQwtZa8TPijoa6YoiC5nkoWD3jhfpeRUI"
	url := "https://api.telegram.org/bot" + bot_token + "/" + type_message
	dateParts := strings.Split(TaskDueDate, "T")
	date := strings.TrimSuffix(dateParts[0], "Z")
	content_text := `
		У Вас новое задание.
		Задача: ` + TaskTitle + `
		Описание: ` + TaskDes + `
		Дата окончания задания: ` + date + `
	`
	for _, worker := range TaskSelectedWorkers {
		message := map[string]string{
			"chat_id": strconv.Itoa(*worker.TgData.ChatID),
			"text":    content_text,
		}
		payload, err := json.Marshal(message)
		if err != nil {
			log.Println(err)
		}

		// Send the HTTP POST request
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
		if err != nil {
			log.Println(err)
		}

		// Parse the response
		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			log.Println(err)
		}
	}

	return nil
}

func (r *DBControllerRepository) GetBoardDone(BoardId int) (int, error) {
	log.Print("GetBoardDone")
	log.Print(BoardId)
	rows, err := r.store.db.Query(`
    SELECT idCard from cards WHERE idBoard = ?
    `, BoardId)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	defer rows.Close()

	var IDCards []int
	for rows.Next() {
		var IDCard int
		err = rows.Scan(&IDCard)
		if err != nil {
			log.Println(err)
			return 0, err
		}
		IDCards = append(IDCards, IDCard)
	}
	log.Print(IDCards)
	done := 0

	if len(IDCards) > 0 {
		for _, card := range IDCards {
			buff, err := r.GetCardDone(card)
			if err != nil {
				log.Println(err)
				return 1, err
			}
			done += buff
		}
	} else {
		done = 1
	}
	return done, err
}

func (r *DBControllerRepository) GetCardDone(CardId int) (int, error) {
	rows, err := r.store.db.Query(`SELECT IF(COUNT(*) > 0, COUNT(*), 1) as count FROM tasks WHERE idCard = ? AND done = 0;`, CardId)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	defer rows.Close()
	var count int
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			log.Println(err)
			return 1, err
		}
	}
	err = rows.Err()
	if err != nil {
		log.Println(err)
		return 1, err
	}
	return count, nil
}
