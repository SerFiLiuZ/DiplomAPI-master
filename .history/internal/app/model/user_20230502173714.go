package model

// User ...
type User struct {
	ID            int     `json:"id"`
	Email         string  `json:"email"`
	Password      string  `json:"password,omitempty"`
	FIO           string  `json:"FIO"`
	Status        string  `json:"status"`
	IdLegalEntity int     `json:"idLegalEntity"`
	AssignedTo    *int    `json:"assignedTo"`
	phoneNumber   *string `json:"phoneNumber"`
	ChatID        *int    `json:"chatId"`
	UserNameTg    *string `json:"userNameTg"`
	PasswordCorp  string  `json:"passwordCorp"`
	NameCorp      string  `json:"nameCorp"`
}
