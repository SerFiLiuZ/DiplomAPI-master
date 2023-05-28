package store

import "github.com/gopherschool/http-rest-api/internal/app/model"

// UserRepository ...
type UserRepository interface {
	Create(*model.User) error
	Find(int) (*model.User, error)
	FindByEmail(string) (*model.User, error)
	FindAllById(string) (string, error)
}

// BoardRepository ...
type BoardRepository interface {
	GetAllBoards(IDUser int) ([]*model.Board, error)
	CreateBoard(Title string, IDUser int) error
	DeleteBoard(IDBoard int) error
}

// CardRepository ...
type CardRepository interface {
	FindCardsByBoardID(IDUser int, BoardID int) ([]*model.Card, error)
	CreateCard(CardTitle string, CardDes string, BoardId int) error
	DeleteCard(IDCard int) error
}

// TaskRepository ...
type TaskRepository interface {
	FindTasksByCardID(IDUser int, CardID int) ([]*model.Task, error)
	DeleteTask(IDTask int) error
}
