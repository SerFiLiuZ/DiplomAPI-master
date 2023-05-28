package sqlstore

import (
	"net/http"

	"github.com/gopherschool/http-rest-api/internal/app/model"
)

type BoardRepository struct {
	store *Store
}

// FindByEmail ...
func (r *BoardRepository) getAllBoards() (*[]model.Board, error) {
	var boards = []*model.Board{}

	// Execute the SQL query
	rows, err := r.store.db.Query("SELECT idBoard, Title, done FROM taskdb.boards")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var board model.Board{}
		err = rows.Scan(&board.ID, &board.Title, &board.Done)
		if err != nil {
			return nil, err
		}
		boards = append(boards, board)
	}

	err = rows.Err()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
