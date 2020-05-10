package models

//easyjson:json
type User struct {
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
	Nickname string `json:"nickname,omitempty"`
	About    string `json:"about,omitempty"`
}

//easyjson:json
type Users []User

//easyjson:json
type UserUpdate struct {
	About    string `json:"about,omitempty"`
	Email    string `json:"email,omitempty"`
	Fullname string `json:"fullname,omitempty"`
}

//easyjson:json
type UserUpdates []UserUpdate

func (user User) Validate() error {
	if err := validateEmail(user.Email); err != nil {
		return err
	}
	if err := validateNickname(user.Nickname); err != nil {
		return err
	}
	return nil
}

func (userUpdate UserUpdate) Validate() error {
	return validateEmail(userUpdate.Email)
}
