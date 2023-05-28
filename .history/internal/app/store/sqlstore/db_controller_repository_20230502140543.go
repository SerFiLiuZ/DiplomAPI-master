package sqlstore

import (
	"log"
	"strconv"

	"github.com/gopherschool/http-rest-api/internal/app/model"
)

type DBControllerRepository struct {
	store *Store
}

func (r *DBControllerRepository) UpdateDoneByTask(IDTask int) error {
	_, err := r.store.db.Exec(`
		UPDATE taskdb.cards 
		SET done = 1 
		WHERE idCard = (
			SELECT idCard 
			FROM taskdb.tasks 
			WHERE idTask = ?
		) AND (
			SELECT COUNT(*) 
			FROM taskdb.tasks 
			WHERE idCard = (
				SELECT idCard 
				FROM taskdb.tasks 
				WHERE idTask = ?
			) AND done = 0
		) = 0;
	`, IDTask, IDTask)
	return err
}

func (r *DBControllerRepository) GetAppications(IDManager int) ([]*model.Appication, error) {
	rows, err := r.store.db.Query(`
	SELECT applications.passwordCorp, applications.FIO, applications.phoneNumber, applications.userName, applications.dataApplication, applications.chatID    
	FROM taskdb.applications WHERE taskdb.applications.passwordCorp = (SELECT users.passwordCorp from taskdb.users where users.idUsers = 1 );
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applications := []*model.Appication{}

	for rows.Next() {
		application := &model.Appication{}
		err = rows.Scan(&application.PasswordCorp, &application.FIO, &application.PhoneNumber, &application.UserName, &application.DataApplication, &application.ChatID)
		if err != nil {
			return nil, err
		}
		applications = append(applications, application)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return applications, nil

}

func (r *DBControllerRepository) AcceptApplication(chatIDUser int, IDManager int) error {
	rows, err := r.store.db.Query(
		"SELECT passwordCorp, FIO, phoneNumber, chatID, userName FROM taskdb.applications WHERE chatID =?",
		chatIDUser,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	u := &model.User{}

	for rows.Next() {
		if err := rows.Scan(&u.PasswordCorp, &u.FIO, &u.PhoneNumber, &u.ChatID, &u.UserNameTg); err != nil {
			return err
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}

	u_manager, err := r.store.User().Find(IDManager)
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = r.store.db.Exec(
		"INSERT INTO taskdb.users (FIO, status, idLegal_entity, email, password, assigned_to, phoneNumder, chatIDInTg, usernameTg, passwordCorp, nameCorp) VALUES (?,?,?,?,?,?,?,?,?,?,?)",
		u.FIO,
		1,
		u_manager.IdLegalEntity,
		"login."+u_manager.NameCorp+"."+strconv.Itoa(*u.ChatID),
		"password."+u_manager.NameCorp+"."+strconv.Itoa(*u.ChatID),
		u_manager.ID,
		u.PhoneNumber,
		u.ChatID,
		u.UserNameTg,
		u_manager.PasswordCorp,
		u_manager.NameCorp,
	)
	if err != nil {
		return err
	}

	err = r.RejectApplication(chatIDUser)
	return err
}

func (r *DBControllerRepository) RejectApplication(chatIDUser int) error {
	_, err := r.store.db.Exec(
		"DELETE FROM taskdb.applications WHERE chatID=?",
		chatIDUser,
	)
	return err
}

func (r *DBControllerRepository) SendMesage(TaskTitle string, TaskDueDate string, TaskSelectedWorkers []*model.User) error {

	return nil
}
