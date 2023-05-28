package sqlstore

import (
	"strconv"

	"github.com/gopherschool/http-rest-api/internal/app/model"
)

type TaskRepository struct {
	store *Store
}

func (r *TaskRepository) FindTasksByCardID(IDUser int, UserStatus string, TaskId int) ([]*model.Task, error) {
	var response string = ""

	if UserStatus == "manager" {
		response = `
		"SELECT taskdb.tasks.idTask, taskdb.tasks.name, taskdb.tasks.description, taskdb.tasks.due_date, taskdb.tasks.done, taskdb.tasks.assigned 
		FROM taskdb.tasks JOIN taskdb.cards 
			ON taskdb.tasks.idCard = taskdb.cards.idCard 
		JOIN taskdb.boards 
			ON taskdb.cards.idBoard = taskdb.boards.idBoard	
				WHERE taskdb.boards.idUser = ? AND taskdb.cards.idCard = ?;"
		`
	}
	if UserStatus == "worker" {
		response = `
		SELECT idBoard, title, done FROM taskdb.boards WHERE idBoard IN (
			SELECT DISTINCT idBoard FROM taskdb.cards WHERE idCard IN (
				SELECT DISTINCT idCard FROM taskdb.tasks WHERE FIND_IN_SET(?, assigned) > 0
			)
		);
		`
	}

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

func (r *TaskRepository) DeleteTask(IDTask int) error {
	_, err := r.store.db.Exec(
		"DELETE FROM taskdb.tasks WHERE idTask=?",
		IDTask,
	)
	return err
}

func (r *TaskRepository) CreateTask(TaskTitle string, TaskDes string, TaskDueDate string, TaskSelectedWorkers []*model.User, CardId int) error {
	assignedIds := ""
	for i, worker := range TaskSelectedWorkers {
		if i > 0 {
			assignedIds += ","
		}
		assignedIds += strconv.Itoa(worker.ID)
	}
	_, err := r.store.db.Exec(
		"INSERT INTO taskdb.tasks (name, description, due_date, idCard, assigned, done) VALUES (?,?,?,?,?,?)",
		TaskTitle,
		TaskDes,
		TaskDueDate,
		CardId,
		assignedIds,
		false,
	)
	return err
}

func (r *TaskRepository) CompliteTask(IDTask int) error {
	_, err := r.store.db.Exec(
		"UPDATE taskdb.tasks SET done=1 WHERE idTask=?",
		IDTask,
	)
	return err
}
