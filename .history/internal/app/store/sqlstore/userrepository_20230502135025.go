package sqlstore

import (
	"database/sql"
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
	if err := r.store.db.QueryRow(
		"SELECT idUsers, FIO, status, idLegal_entity, email, password, assigned_to, phoneNumder, chatIDInTg, usernameTg, passwordCorp, nameCorp FROM users WHERE idUsers = ?",
		id,
	).Scan(
		&u.ID,
		&u.FIO,
		&u.Status,
		&u.IdLegalEntity,
		&u.Email,
		&u.Password,
		&u.AssignedTo,
		&u.PhoneNumber,
		&u.ChatID,
		&u.UserNameTg,
		&u.PasswordCorp,
		&u.NameCorp,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return u, nil
}

// FindByEmail ...
func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.QueryRow(
		"SELECT idUsers, email, password, status FROM users WHERE email = ?",
		email,
	).Scan(
		&u.ID,
		&u.Email,
		&u.Password,
		&u.Status,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

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
			return "", err
		}

		u, err1 := r.store.userRepository.Find(IdInt)
		if err1 != nil {
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

	rows, err := r.store.db.Query(
		"SELECT idUsers, FIO, phoneNumder, chatIDInTg, usernameTg FROM users WHERE assigned_to =?",
		IDManager,
	)
	if err != nil {
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
			&u.ChatID,
			&u.UserNameTg,
		); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *UserRepository) DeliteUser(FIO string) error {
	_, err := r.store.db.Exec(
		"DELETE FROM taskdb.users WHERE FIO =?",
		FIO,
	)
	return err
}
