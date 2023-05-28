package sqlstore

import (
	"log"

	"github.com/gopherschool/http-rest-api/internal/app/model"
)

type CardRepository struct {
	store *Store
}

func (r *CardRepository) FindCardsByBoardID(IDUser int, UserStatus string, BoardID int) ([]*model.Card, error) {

	var response string = ""

	if UserStatus == "manager" {
		response = `
			SELECT idCard, name, description, done, idBoard 
				FROM cards 
				WHERE idBoard = ? AND idBoard 
				IN (SELECT idBoard FROM boards WHERE idUser = ?)`
	}
	if UserStatus == "worker" {
		response = `
			SELECT idCard, name, description, done, idBoard FROM cards WHERE idBoard = ? AND idCard IN (
				SELECT DISTINCT idCard FROM tasks WHERE FIND_IN_SET(?, assigned) > 0
		);
		`
	}
	// Execute the SQL query
	rows, err := r.store.db.Query(response, BoardID, IDUser)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	// Loop through the rows of data and store it in a slice of Board structs
	boardCards := []*model.Card{}

	for rows.Next() {
		card := &model.Card{}
		err = rows.Scan(&card.ID, &card.Name, &card.Description, &card.Done, &card.IDBoard)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		boardCards = append(boardCards, card)
	}

	err = rows.Err()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return boardCards, nil
}

func (r *CardRepository) FindCardsByBoardIDSimple(BoardID int) ([]*model.Card, error) {
	// Execute the SQL query
	rows, err := r.store.db.Query(`
		SELECT idCard, idBoard, name, description, done FROM cards WHERE idBoard=?`,
		BoardID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	// Loop through the rows of data and store it in a slice of Board structs
	boardCards := []*model.Card{}

	for rows.Next() {
		card := &model.Card{}
		err = rows.Scan(&card.ID, &card.IDBoard, &card.Name, &card.Description, &card.Done, &card.IDBoard)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		boardCards = append(boardCards, card)
	}

	err = rows.Err()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return boardCards, nil
}

func (r *CardRepository) GetCard(IDCard int) (*model.Card, error) {
	rows := r.store.db.QueryRow(`
        SELECT idCard, idBoard, name, description, done FROM cards WHERE idCard = ?`,
		IDCard)

	c := &model.Card{}

	if err := rows.Scan(&c.ID, &c.IDBoard, &c.Name, &c.Description, &c.Done); err != nil {
		log.Println(err)
		return nil, err
	}

	return c, nil
}

func (r *CardRepository) InsertCardToDeleteCard(IDCard int) error {
	c, err := r.GetCard(IDCard)
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = r.store.db.Exec(
		"INSERT INTO deleted_cards (last_idCard, last_idBoard, last_name, last_description, last_done) VALUES (?,?,?,?,?);",
		c.ID, c.IDBoard, c.Name, c.Description, c.Done)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (r *CardRepository) DeleteCard(IDCard int) error {
	_, err := r.store.db.Exec(
		"DELETE FROM cards WHERE idCard=?",
		IDCard,
	)
	log.Println(err)
	return err
}
func (r *CardRepository) CreateCard(CardTitle string, CardDes string, BoardId int) error {
	_, err := r.store.db.Exec(
		"INSERT INTO cards (name, description, done, idBoard) VALUES (?,?,?,?)",
		CardTitle,
		CardDes,
		false,
		BoardId,
	)

	log.Print(err)
	return err
}
