package sqlstore

import (
	"log"

	"github.com/gopherschool/http-rest-api/internal/app/model"
)

type BoardRepository struct {
	store *Store
}

// getAllBoards ...
func (r *BoardRepository) GetAllBoards(IDUser int, UserStarus string) ([]*model.Board, error) {

	var response string = ""

	if UserStarus == "manager" {
		response = "SELECT idBoard, title, done FROM boards Where idUser = ?"
	}
	if UserStarus == "worker" {
		response = `
		SELECT idBoard, title, done FROM boards WHERE idBoard IN (
			SELECT DISTINCT idBoard FROM cards WHERE idCard IN (
				SELECT DISTINCT idCard FROM tasks WHERE FIND_IN_SET(?, assigned) > 0
			)
		);
		`
	}

	rows, err := r.store.db.Query(
		response,
		IDUser,
	)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	boards := []*model.Board{}

	for rows.Next() {
		b := &model.Board{}
		if err := rows.Scan(&b.ID, &b.Title, &b.Done); err != nil {
			log.Println(err)
			return nil, err
		}
		boards = append(boards, b)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return boards, nil
}

func (r *BoardRepository) GetBoard(IDBoard int) (*model.Board, error) {
	rows, err := r.store.db.Query(`
		SELECT idBoard, Title, idUser, done FROM boards WHERE idBoard =?`,
		IDBoard)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer rows.Close()

	b := &model.Board{}
	if err := rows.Scan(&b.ID, &b.Title, &b.IDUser, &b.Done); err != nil {
		log.Println(err)
		return nil, err
	}
	return b, nil

}

func (r *BoardRepository) CreateBoard(Title string, IDUser int) error {
	_, err := r.store.db.Exec(
		"INSERT INTO boards (title, idUser, done) VALUES (?,?,?)",
		Title,
		IDUser,
		false,
	)
	log.Println(err)
	return err
}

func (r *BoardRepository) DeleteBoard(IDBoard int) error {
	_, err := r.store.db.Exec(
		"INSERT INTO deleted_boards (last_idBoard, last_Title, last_idUser, last_done) VALUES (?,?,?)",
	)
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = r.store.db.Exec(
		"DELETE FROM boards WHERE idBoard=?",
		IDBoard,
	)
	log.Println(err)
	return err
}
