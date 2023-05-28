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
	SELECT taskdb.applications.passwordCorp, taskdb.applications.FIO, taskdb.applications.phoneNumber, taskdb.applications.userName, taskdb.applications.dataApplication  
	FROM taskdb.applications WHERE taskdb.applications.passwordCorp = (SELECT users.passwordCorp from taskdb.users where users.idUsers = 1 );
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appications []*model.Appication

}
