package sqlstore

import (
	"github.com/gopherschool/http-rest-api/internal/app/model"
)

type BoardRepository struct {
	store *Store
}

// getAllBoards ...
func (r *BoardRepository) GetAllBoards(IDUser int, UserStarus string) ([]*model.Board, error) {
	if UserStarus == "manager" {
		rows, err := r.store.db.Query(
			"SELECT idBoard, title, done FROM taskdb.boards Where idUser = ?",
			IDUser,
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

	return nil, nil
}

func (r *BoardRepository) CreateBoard(Title string, IDUser int) error {
	_, err := r.store.db.Exec(
		"INSERT INTO taskdb.boards (title, idUser, done) VALUES (?,?,?)",
		Title,
		IDUser,
		false,
	)
	return err
}

func (r *BoardRepository) DeleteBoard(IDBoard int) error {
	_, err := r.store.db.Exec(
		"DELETE FROM taskdb.boards WHERE idBoard=?",
		IDBoard,
	)
	return err
}
