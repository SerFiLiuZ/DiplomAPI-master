package sqlstore

import (
	"github.com/gopherschool/http-rest-api/internal/app/model"
)

type CardRepository struct {
	store *Store
}

func (r *CardRepository) FindCardsByBoardID(IDUser int, UserStatus string, BoardID int) ([]*model.Card, error) {
	// Execute the SQL query
	rows, err := r.store.db.Query("SELECT idCard, name, description, done, idBoard FROM taskdb.cards WHERE idBoard = ? AND idBoard IN (SELECT idBoard FROM taskdb.boards WHERE idUser = ?)", BoardID, IDUser)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Loop through the rows of data and store it in a slice of Board structs
	boardCards := []*model.Card{}

	for rows.Next() {
		card := &model.Card{}
		err = rows.Scan(&card.ID, &card.Name, &card.Description, &card.Done, &card.IDBoard)
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

func (r *CardRepository) DeleteCard(IDCard int) error {
	_, err := r.store.db.Exec(
		"DELETE FROM taskdb.cards WHERE idCard=?",
		IDCard,
	)
	return err
}
func (r *CardRepository) CreateCard(CardTitle string, CardDes string, BoardId int) error {
	_, err := r.store.db.Exec(
		"INSERT INTO taskdb.cards (name, description, done, idBoard) VALUES (?,?,?,?)",
		CardTitle,
		CardDes,
		false,
		BoardId,
	)

	return err
}
