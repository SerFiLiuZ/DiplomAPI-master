package model

// User ...
type User struct {
	ID           int          `json:"id"`
	FIO          string       `json:"FIO"`
	Status       string       `json:"status"`
	Autorization Autorization `json:"idAutorization"`
	LegalEntity  LegalEntity  `json:"idLegalEntity"`
	TgData       TgData       `json:"idTgData"`
	AssignedTo   *int         `json:"assignedTo"`
	PhoneNumber  *string      `json:"phoneNumber"`
}
