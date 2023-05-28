package sqlstore

import (
	"github.com/gopherschool/http-rest-api/internal/app/model"
)

type BoardRepository struct {
	store *Store
}

// getAllBoards ...
func (r *BoardRepository) GetAllBoards(IDUser int, UserStarus string) ([]*model.Board, error) {

	var response string = ""

	if UserStarus == "manager" {
		response = "SELECT idBoard, title, done FROM taskdb.boards Where idUser = ?"
	}
	if UserStarus == "worker" {
		response = `
		SELECT idBoard, title, done FROM taskdb.boards WHERE idBoard IN (
			SELECT DISTINCT idBoard FROM taskdb.cards WHERE idCard IN (
				SELECT DISTINCT idCard FROM taskdb.tasks WHERE FIND_IN_SET(?, assigned) > 0
			)
		);
		`
	}

	rows, err := r.store.db.Query(
		response,
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
