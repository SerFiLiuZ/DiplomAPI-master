package store

// Store ...
type Store interface {
	User() UserRepository
	Board() BoardRepository
	Card() CardRepository
	Task() TaskRepository
	DBController() DBControllerRepository
}
