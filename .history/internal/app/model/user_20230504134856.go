package model

// User ...
type User struct {
	ID            int     `json:"id"`
	Email         string  `json:"email"`
	Password      string  `json:"password"`
	FIO           string  `json:"FIO"`
	Status        string  `json:"status"`
	IdLegalEntity int     `json:"idLegalEntity"`
	AssignedTo    *int    `json:"assignedTo"`
	PhoneNumber   *string `json:"phoneNumber"`
	ChatID        *int    `json:"chatId"`
	UserNameTg    *string `json:"userNameTg"`
	PasswordCorp  string  `json:"passwordCorp"`
	NameCorp      string  `json:"nameCorp"`
}
