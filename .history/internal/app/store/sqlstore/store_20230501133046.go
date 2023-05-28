package sqlstore

import (
	"database/sql"

	"github.com/gopherschool/http-rest-api/internal/app/store"
)

// Store ...
type Store struct {
	db                     *sql.DB
	userRepository         *UserRepository
	cardRepository         *CardRepository
	boardRepository        *BoardRepository
	taskRepository         *TaskRepository
	dbcontrollerRepository *DBControllerTaskRepository
}

// New ...
func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

// User ...
func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}

	return s.userRepository
}

func (s *Store) Board() store.BoardRepository {
	if s.boardRepository != nil {
		return s.boardRepository
	}

	s.boardRepository = &BoardRepository{
		store: s,
	}

	return s.boardRepository
}

// Card ...
func (s *Store) Card() store.CardRepository {
	if s.cardRepository != nil {
		return s.cardRepository
	}

	s.cardRepository = &CardRepository{
		store: s,
	}

	return s.cardRepository
}

// Task ...
func (s *Store) Task() store.TaskRepository {
	if s.taskRepository != nil {
		return s.taskRepository
	}

	s.taskRepository = &TaskRepository{
		store: s,
	}

	return s.taskRepository
}

func (s *Store) DBController() store.DBControllerRepository {
	if s.dbcontrollerRepository != nil {
		return s.dbcontrollerRepository
	}

	s.dbcontrollerRepository = &DBControllerRepository{
		store: s,
	}

	return s.dbcontrollerRepository
}
