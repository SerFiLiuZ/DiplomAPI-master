package sqlstore

import "github.com/gopherschool/http-rest-api/internal/app/model"

type TaskRepository struct {
	store *Store
}

func (r *TaskRepository) Find(IDUser int, TaskId int) (*model.Task, error)
