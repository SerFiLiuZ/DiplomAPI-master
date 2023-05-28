package sqlstore

import "github.com/gopherschool/http-rest-api/internal/app/model"

type DBControllerRepository struct {
	store *Store
}

func (r *DBControllerRepository) UpdateDoneByTask(IDTask int) error {
	_, err := r.store.db.Exec(`
		UPDATE taskdb.cards 
		SET done = 1 
		WHERE idCard = (
			SELECT idCard 
			FROM taskdb.tasks 
			WHERE idTask = ?
		) AND (
			SELECT COUNT(*) 
			FROM taskdb.tasks 
			WHERE idCard = (
				SELECT idCard 
				FROM taskdb.tasks 
				WHERE idTask = ?
			) AND done = 0
		) = 0;
	`, IDTask, IDTask)
	return err
}

func (r *DBControllerRepository) GetAppications(IDManager int) ([]*model.Appication, error) {
	rows, err := r.store.db.Query(`
	SELECT applications.passwordCorp, applications.FIO, applications.phoneNumber, applications.userName, applications.dataApplication, applications.chatID    
	FROM taskdb.applications WHERE taskdb.applications.passwordCorp = (SELECT users.passwordCorp from taskdb.users where users.idUsers = 1 );
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applications := []*model.Appication{}

	for rows.Next() {
		application := &model.Appication{}
		err = rows.Scan(&application.PasswordCorp, &application.FIO, &application.PhoneNumber, &application.UserName, &application.DataApplication, &application.ChatID)
		if err != nil {
			return nil, err
		}
		applications = append(applications, application)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return applications, nil

}

func (r *DBControllerRepository) AcceptApplication(chatIDUser int, IDManager int) error {

	err := r.RejectApplication(chatIDUser)
	return err
}

func (r *DBControllerRepository) RejectApplication(chatIDUser int) error {
	_, err := r.store.db.Exec(
		"DELETE FROM taskdb.applications WHERE chatID=?",
		chatIDUser,
	)
	return err
}
