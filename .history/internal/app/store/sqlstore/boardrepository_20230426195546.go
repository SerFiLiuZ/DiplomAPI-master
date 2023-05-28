package sqlstore

import (
	"database/sql"

	"github.com/gopherschool/http-rest-api/internal/app/model"
	"github.com/gopherschool/http-rest-api/internal/app/store"
)

type BoardRepository struct {
	store *Store
}

// FindByEmail ...
func (r *UserRepository) getAllBoards() (*[]model.Board, error) {
	u := &model.Board{}
	if err := r.store.db.QueryRow(
		"SELECT idUsers, email, password FROM users WHERE email = ?",
		email,
	).Scan(
		&u.ID,
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
