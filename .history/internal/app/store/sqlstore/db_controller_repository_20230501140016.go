package sqlstore

type DBControllerRepository struct {
	store *Store
}

func (r *DBControllerRepository) UpdateDoneByTask(IDTask int) error {

	return nil
}

// UPDATE taskdb.cards
// SET done = 1
// WHERE idCard = (
//     SELECT idCard
//     FROM taskdb.tasks
//     WHERE idTask = 10
// )
// AND (
//     SELECT COUNT(*)
//     FROM taskdb.tasks
//     WHERE idCard = (
//         SELECT idCard
//         FROM taskdb.tasks
//         WHERE idTask = 10
//     )
//     AND done = 0
// ) = 0;
