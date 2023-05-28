package model

type Board struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	IDUser int    `json:"id_user"`
	Done   bool   `json:"done"`
}
