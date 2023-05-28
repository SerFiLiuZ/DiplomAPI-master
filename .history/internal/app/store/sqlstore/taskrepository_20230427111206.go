package sqlstore

import "github.com/gopherschool/http-rest-api/internal/app/model"

type TaskRepository struct {
	store *Store
}

func (r *TaskRepository) FindTaskByCardID(IDUser int, TaskId int) ([]*model.Task, error) {
	// Execute the SQL query
	rows, err := r.store.db.Query("SELECT taskdb.tasks.idTask, taskdb.tasks.name, taskdb.tasks.description, taskdb.tasks.due_date, taskdb.tasks.done, taskdb.tasks.assigned FROM taskdb.tasks JOIN taskdb.cards ON taskdb.tasks.idCard = taskdb.cards.idCard JOIN taskdb.boards ON taskdb.cards.idBoard = taskdb.boards.idBoard	WHERE taskdb.boards.idUser = ? AND taskdb.cards.idCard = ?;", IDUser, TaskId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Loop through the rows of data and store it in a slice of Board structs
	tasks := []*model.Task{}

	for rows.Next() {
		task := &model.Task{}
		err = rows.Scan(&task.ID, &task.Name, &task.Description, &task.Due_date, &task.Done, &task.Assigned)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return tasks, nil
}
