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

// Create ...
func (r *UserRepository) Create(u *model.User) error {
	if err := u.Validate(); err != nil {
		return err
	}

	if err := u.BeforeCreate(); err != nil {
		return err
	}

	return r.store.db.QueryRow(
		"INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id",
		u.Email,
		u.Password,
	).Scan(&u.ID)
}

// Find ...
func (r *UserRepository) Find(id int) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.QueryRow(
		"SELECT idUsers, FIO, status, idLegal_entity, email, password FROM users WHERE idUsers = ?",
		id,
	).Scan(
		&u.ID,
		&u.FIO,
		&u.Status,
		&u.IdLegalEntity,
		&u.Email,
		&u.Password,
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
