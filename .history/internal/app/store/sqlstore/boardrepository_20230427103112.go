package sqlstore

import (
	"github.com/gopherschool/http-rest-api/internal/app/model"
)

type BoardRepository struct {
	store *Store
}

// getAllBoards ...
func (r *BoardRepository) GetAllBoards(IDUser int) ([]*model.Board, error) {
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

func (r *BoardRepository) Find(IDUser int, BoardId int) ([]*model.Card, error) {
	// Execute the SQL query
	rows, err := r.store.db.Query("SELECT idCards, name, description, done FROM cards Where idBoard = ?", BoardId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Loop through the rows of data and store it in a slice of Board structs
	boardCards := []*model.Card{}

	for rows.Next() {
		card := &model.Card{}
		err = rows.Scan(&card.ID, &card.Name, &card.Description, &card.Done)
		if err != nil {
			return nil, err
		}
		boardCards = append(boardCards, card)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return boardCards, nil
}
