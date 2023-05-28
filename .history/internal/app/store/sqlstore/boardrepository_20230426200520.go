package sqlstore

import (
	"github.com/gopherschool/http-rest-api/internal/app/model"
)

type BoardRepository struct {
	store *Store
}

// FindByEmail ...
func (r *BoardRepository) getAllBoards() ([]*model.Board, error) {
	rows, err := r.store.db.Query(
		"SELECT id, title, done FROM boards",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	boards := []*model.Board{}

	for rows.Next() {
		b := &model.Board{}
		if err := rows.Scan(&b.ID, &b.Title, &b.Done); err != nil {
			return nil, err
		}
		boards = append(boards, b)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return boards, nil
}
