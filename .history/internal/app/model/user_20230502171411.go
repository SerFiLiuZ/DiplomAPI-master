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
	PhoneNumber   *string `json:"phoneNumber"`
	ChatID        *int    `json:"chatId"`
	UserNameTg    *string `json:"userNameTg"`
	PasswordCorp  string  `json:"passwordCorp"`
	NameCorp      string  `json:"nameCorp"`
}

// // BeforeCreate ...
// func (u *User) BeforeCreate() error {
// 	if len(u.Password) > 0 {
// 		enc, err := encryptString(u.Password)
// 		if err != nil {
// 			return err
// 		}

// 		u.Password = enc
// 	}

// 	return nil
// }

// // Sanitize ...
// func (u *User) Sanitize() {
// 	u.Password = ""
// }

// // ComparePassword ...
// func (u *User) ComparePassword(password string) bool {
// 	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil
// }

// func encryptString(s string) (string, error) {
// 	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
// 	if err != nil {
// 		return "", err
// 	}

// 	return string(b), nil
// }
