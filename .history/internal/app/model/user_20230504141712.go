package model

// User ...
type User struct {
	ID          int         `json:"id"`
	FIO         string      `json:"FIO"`
	Status      string      `json:"status"`
	LegalEntity LegalEntity `json:"idLegalEntity"`
	AssignedTo  *int        `json:"assignedTo"`
	PhoneNumber *string     `json:"phoneNumber"`
}
