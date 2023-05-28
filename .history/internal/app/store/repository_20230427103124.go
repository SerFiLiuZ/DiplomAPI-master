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
	Find(IDUser int, BoardId int) ([]*model.Card, error)
}

// CardRepository ...
type CardRepository interface {
}

// TaskRepository ...
type TaskRepository interface {
}
