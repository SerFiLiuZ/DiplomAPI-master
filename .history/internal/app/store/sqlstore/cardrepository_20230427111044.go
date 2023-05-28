package sqlstore

import "github.com/gopherschool/http-rest-api/internal/app/model"

type CardRepository struct {
	store *Store
}

func (r *CardRepository) Find(IDUser int, IDBoard int) ([]*model.Card, error) {
	// Execute the SQL query
	rows, err := r.store.db.Query("SELECT idCards, name, description, done FROM taskdb.cards WHERE idBoard = ? AND idBoard IN (SELECT idBoard FROM taskdb.boards WHERE idUser = ?)", IDBoard, IDUser)
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
