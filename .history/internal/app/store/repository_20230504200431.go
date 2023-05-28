package store

import "github.com/gopherschool/http-rest-api/internal/app/model"

// UserRepository ...
type UserRepository interface {
	Find(int) (*model.User, error)
	FindByEmail(string) (*model.User, error)
	FindAllById(string) (string, error)
	FindAllByIdManager(IDManager int) ([]*model.User, error)
	DeliteUser(ID int) error
}

// BoardRepository ...
type BoardRepository interface {
	GetBoard(IDBoard int) (*model.Board, error)
	GetAllBoards(IDUser int, UserStatus string) ([]*model.Board, error)
	CreateBoard(Title string, IDUser int) error
	DeleteBoard(IDBoard int) error
	InsertBoardToDeleteBoard(IDBoard int) error
}

// CardRepository ...
type CardRepository interface {
	FindCardsByBoardID(IDUser int, UserStatus string, BoardID int) ([]*model.Card, error)
	CreateCard(CardTitle string, CardDes string, BoardId int) error
	DeleteCard(IDCard int) error
	FindCardsByBoardIDSimple(BoardID int) ([]*model.Card, error)
	InsertCardToDeleteCard(IDCard int) error
	GetCard(IDCard int) (*model.Card, error)
}

// TaskRepository ...
type TaskRepository interface {
	FindTasksByCardID(IDUser int, UserStatus string, CardID int) ([]*model.Task, error)
	DeleteTask(IDTask int) error
	CreateTask(TaskTitle string, TaskDes string, TaskDueDate string, TaskSelectedWorkers []*model.User, CardId int) error
	CompliteTask(IDTask int) error
	InsertTaskToDeleteTask(IDTask int) error
	GetTask(IDTask int) (*model.Task, error)
	FindTasksByCardIDSimple(CardID int) ([]*model.Task, error)
}

type DBControllerRepository interface {
	UpdateDoneByTask(IDTask int) error
	GetAppications(IDManager int) ([]*model.Appication, error)
	AcceptApplication(chatIDUser int, IDManager int) error
	RejectApplication(chatIDUser int) error
	SendMesage(TaskTitle string, TaskDes string, TaskDueDate string, TaskSelectedWorkers []*model.User) error
	GetBoardDone(BoardId int) (int, error)
	GetCardDone(CardId int) (int, error)
}
