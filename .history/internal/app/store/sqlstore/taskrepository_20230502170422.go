package sqlstore

import (
	"log"
	"strconv"

	"github.com/gopherschool/http-rest-api/internal/app/model"
)

type TaskRepository struct {
	store *Store
}

func (r *TaskRepository) FindTasksByCardID(IDUser int, UserStatus string, CardId int) ([]*model.Task, error) {
	var response string = ""

	if UserStatus == "manager" {
		response = `
		SELECT tasks.idTask, tasks.name, tasks.description, tasks.due_date, tasks.done, tasks.assigned 
		FROM tasks JOIN cards ON tasks.idCard = cards.idCard 
		JOIN boards ON cards.idBoard = boards.idBoard	
		WHERE boards.idUser = ? AND cards.idCard = ?;
		`
	}
	if UserStatus == "worker" {
		response = `
		SELECT tasks.idTask, tasks.name, tasks.description, tasks.due_date, tasks.done, tasks.assigned
		FROM tasks WHERE FIND_IN_SET(?, assigned) > 0 AND idCard = ?;
		`
	}

	// Execute the SQL query
	rows, err := r.store.db.Query(response, IDUser, CardId)
	if err != nil {
		log.Println(err)
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
		log.Println(err)
		return nil, err
	}

	return tasks, nil
}

func (r *TaskRepository) DeleteTask(IDTask int) error {
	_, err := r.store.db.Exec(
		"DELETE FROM tasks WHERE idTask=?",
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
		"INSERT INTO tasks (name, description, due_date, idCard, assigned, done) VALUES (?,?,?,?,?,?)",
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
