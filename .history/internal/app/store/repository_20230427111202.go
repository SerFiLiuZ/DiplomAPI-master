package store

import "github.com/gopherschool/http-rest-api/internal/app/model"

// UserRepository ...
type UserRepository interface {
	Create(*model.User) error
	Find(int) (*model.User, error)
	FindByEmail(string) (*model.User, error)
}

// BoardRepository ...
type BoardRepository interface {
	GetAllBoards(IDUser int) ([]*model.Board, error)
}

// CardRepository ...
type CardRepository interface {
	FindCardByBoardID(IDUser int, BoardID int) ([]*model.Card, error)
}

// TaskRepository ...
type TaskRepository interface {
	FindTaskByCardID(IDUser int, CardID int) ([]*model.Task, error)
}
