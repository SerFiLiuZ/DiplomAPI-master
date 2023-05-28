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

// Create ...
func (r *UserRepository) Create(u *model.User) error {
	if err := u.Validate(); err != nil {
		return err
	}

	if err := u.BeforeCreate(); err != nil {
		return err
	}

	return r.store.db.QueryRow(
		"INSERT INTO users (email, encrypted_password) VALUES ($1, $2) RETURNING id",
		u.Email,
		u.EncryptedPassword,
	).Scan(&u.ID)
}

// Find ...
func (r *UserRepository) Find(id int) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.QueryRow(
		"SELECT id, email, password, FIO FROM users WHERE id = ?",
		id,
	).Scan(
		&u.ID,
		&u.Email,
		&u.Password,
		&u.FIO,
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
	log.Print(splitAignedIdUser)

	AssignedFIOUser := ""

	for _, IdString := range splitAignedIdUser {
		IdInt, err := strconv.Atoi(IdString)
		log.Print(IdInt)
		if err != nil {
			return "", err
		}
		u, err1 := r.store.userRepository.Find(IdInt)
		if err1 != nil {
			return "", err1
		}
		log.Print(u.FIO)
		AssignedFIOUser = AssignedFIOUser + ", " + u.FIO
	}

	return AssignedFIOUser, store.ErrRecordNotFound
}
