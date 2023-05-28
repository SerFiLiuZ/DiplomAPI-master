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
			log.Println(err)
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

func (r *TaskRepository) GetTask(IDTask int) (*model.Task, error) {
	rows := r.store.db.QueryRow(`
        SELECT idTask, name, description, due_date, idCard, assigned, done FROM tasks WHERE idTask = ?`,
		IDTask)

	t := &model.Task{}

	if err := rows.Scan(&t.ID, &t.Name, &t.Description, &t.Due_date, &t.CardID, &t.Assigned, &t.Done); err != nil {
		log.Println(err)
		return nil, err
	}

	return t, nil
}

func (r *TaskRepository) InsertTaskToDeleteTask(IDTask int) error {
	t, err := r.GetTask(IDTask)
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = r.store.db.Exec(
		"INSERT INTO deleted_tasks (last_idTask, last_name, last_description, last_due_date, last_idCard, last_assigned, last_done) VALUES (?,?,?,?,?,?,?);",
		IDTask,
		t.Name,
		t.Description,
		t.Due_date,
		t.CardID,
		t.Assigned,
		t.Done)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (r *TaskRepository) DeleteTask(IDTask int) error {
	log.Print("Переношу таск")
	err := r.InsertTaskToDeleteTask(IDTask)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Print("Перенёс таск")
	log.Print("Удаляю таск")
	_, err = r.store.db.Exec(
		"DELETE FROM tasks WHERE idTask=?",
		IDTask,
	)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Print("Удалил таск")
	return nil
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
	log.Println(err)
	return err
}

func (r *TaskRepository) CompliteTask(IDTask int) error {
	_, err := r.store.db.Exec(
		"UPDATE tasks SET done=1 WHERE idTask=?",
		IDTask,
	)
	log.Println(err)
	return err
}

func (r *TaskRepository) FindTasksByCardIDSimple(CardID int) ([]*model.Task, error) {
	// Execute the SQL query
	rows, err := r.store.db.Query(`
		SELECT idTask, name, description, due_date, idCard, assigned, done FROM tasks WHERE idCard=?`,
		CardID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	// Loop through the rows of data and store it in a slice of Board structs
	tasks := []*model.Task{}

	for rows.Next() {
		task := &model.Task{}
		err = rows.Scan(&task.ID, &task.Name, &task.Description, &task.Due_date, &task.CardID, &task.Assigned, &task.Done)
		if err != nil {
			log.Println(err)
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
