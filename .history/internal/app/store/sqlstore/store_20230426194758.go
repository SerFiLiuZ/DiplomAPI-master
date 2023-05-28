package sqlstore

import (
	"database/sql"

	"github.com/gopherschool/http-rest-api/internal/app/store"
)

// Store ...
type Store struct {
	db              *sql.DB
	userRepository  *UserRepository
	cardRepository  *CardRepository
	boardRepository *BoardRepository
	taskRepository  *TaskRepository
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

// User ...
func (s *Store) Board() store.BoardRepository {
	if s.boardRepository != nil {
		return s.boardRepository
	}

	s.boardRepository = &BoardRepository{
		store: s,
	}

	return s.boardRepository
}
