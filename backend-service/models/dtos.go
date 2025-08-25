package models

type PostUserRegister struct {
	Name                 string `json:"name" validate:"required"`
	Username             string `json:"username" validate:"required,min=3,max=20"`
	Email                string `json:"email" validate:"required,email"`
	Password             string `json:"password" validate:"required,min=6"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,eqfield=Password"`
}

type PostUserLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type PostLogAppend struct {
	Text      string `json:"text" validate:"required"`
	IsSystem  bool   `json:"is_system" validate:"boolean"`
	IsPrivate bool   `json:"is_private" validate:"boolean"`
}

type ResponseGetLogs struct {
}
