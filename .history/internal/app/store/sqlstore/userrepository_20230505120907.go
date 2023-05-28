package sqlstore

import (
	"database/sql"
	"log"
	"strconv"
	"strings"

	"github.com/gopherschool/http-rest-api/internal/app/model"
	"github.com/gopherschool/http-rest-api/internal/app/store"
)

// UserRepository ...
type UserRepository struct {
	store *Store
}

// Find ...
func (r *UserRepository) Find(id int) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.QueryRow(`
		SELECT u.idUsers, u.fio, u.status, u.assigned_to, u.phoneNumber, a.idAutorization, a.login, a.password, l.idLegalEntity, l.nameCorp, l.passwordCorp, t.idtgdata, t.chatID, t.userName
			FROM users u
			LEFT JOIN autorization a ON u.idAutorization = a.idAutorization
			LEFT JOIN legalentity l ON u.idLegalEntity = l.idLegalEntity
			LEFT JOIN tgdata t ON u.idTgData = t.idtgdata
			WHERE u.idUsers = ?`,
		id,
	).Scan(
		&u.ID,
		&u.FIO,
		&u.Status,
		&u.AssignedTo,
		&u.PhoneNumber,
		&u.Autorization.ID,
		&u.Autorization.Email,
		&u.Autorization.Password,
		&u.LegalEntity.ID,
		&u.LegalEntity.NameCorp,
		&u.LegalEntity.PasswordCorp,
		&u.TgData.ID,
		&u.TgData.ChatID,
		&u.TgData.UserNameTg,
	); err != nil {
		if err == sql.ErrNoRows {
			log.Println(err)
			return nil, store.ErrRecordNotFound
		}
		log.Println(err)
		return nil, err
	}

	return u, nil
}

// FindByEmail ...
func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.QueryRow(`
	SELECT u.idUsers, u.status, a.login, a.password
		FROM autorization a
		LEFT JOIN users u ON u.idAutorization = a.idAutorization
		WHERE a.login = ?`,
		email,
	).Scan(
		&u.ID,
		&u.Status,
		&u.Autorization.Email,
		&u.Autorization.Password,
	); err != nil {
		if err == sql.ErrNoRows {
			log.Println(err)
			return nil, store.ErrRecordNotFound
		}
		log.Println(err)
		return nil, err
	}

	return u, nil
}

func (r *UserRepository) FindAllById(AssignedIdUser string) (string, error) {
	splitAignedIdUser := strings.Split(AssignedIdUser, ",")

	AssignedFIOUser := ""

	for _, IdString := range splitAignedIdUser {
		IdInt, err := strconv.Atoi(IdString)

		if err != nil {
			log.Println(err)
			return "", err
		}

		u, err1 := r.store.userRepository.Find(IdInt)
		if err1 != nil {
			log.Println(err)
			return "", err1
		}

		AssignedFIOUser = AssignedFIOUser + ", " + u.FIO
	}

	if len(AssignedFIOUser) < 2 {
		return "", store.ErrRecordNotFound
	}

	return AssignedFIOUser[2:], nil
}

func (r *UserRepository) FindAllByIdManager(IDManager int) ([]*model.User, error) {
	rows, err := r.store.db.Query(`
		SELECT u.idUsers, u.FIO, u.phoneNumber, t.chatID, t.userName
			FROM users u 
			LEFT JOIN tgdata t ON u.idTgData = t.idtgdata
			WHERE assigned_to = ?
		`,
		IDManager,
	)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()
	var users []*model.User
	for rows.Next() {
		u := &model.User{}
		if err := rows.Scan(
			&u.ID,
			&u.FIO,
			&u.PhoneNumber,
			&u.TgData.ChatID,
			&u.TgData.UserNameTg,
		); err != nil {
			log.Println(err)
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *UserRepository) InsertUserToDeleteUser(ID int) error {
	u, err := r.Find(ID)
}

func (r *UserRepository) DeliteUser(ID int) error {

	_, err = r.store.db.Exec(
		"DELETE FROM users WHERE idUsers =?",
		ID,
	)
	log.Println(err)
	return err
}
