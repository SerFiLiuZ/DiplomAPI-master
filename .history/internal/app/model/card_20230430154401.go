package model

type Card struct {
	ID          int    `json:"id"`
	IDBoard     int    `json:"id_board"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}
