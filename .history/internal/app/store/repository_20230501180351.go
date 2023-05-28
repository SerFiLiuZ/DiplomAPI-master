package store

import "github.com/gopherschool/http-rest-api/internal/app/model"

// UserRepository ...
type UserRepository interface {
	Find(int) (*model.User, error)
	FindByEmail(string) (*model.User, error)
	FindAllById(string) (string, error)
	FindAllByIdManager(IDManager int) ([]*model.User, error)
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
	CreateTask(TaskTitle string, TaskDes string, TaskDueDate string, TaskSelectedWorkers []*model.User, CardId int) error
	CompliteTask(IDTask int) error
}

type DBControllerRepository interface {
	UpdateDoneByTask(IDTask int) error
	GetAppications(IDManager int) ([]*model.Appication, error)
	AcceptApplication(chatIDUser int, IDManager int) error
	RejectApplication(chatIDUser int, IDManager int) error
}
