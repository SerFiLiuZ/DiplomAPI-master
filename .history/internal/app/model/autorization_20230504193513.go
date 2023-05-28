package model

type Autorization struct {
	ID       int    `json:"idAutorization"`
	Email    string `json:"login"`
	Password string `json:"password"`
}
