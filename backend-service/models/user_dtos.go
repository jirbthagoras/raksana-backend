package models

type ResponseGetUserProfileStatistic struct {
	Id                 int     `json:"id"`
	Name               string  `json:"name"`
	Username           string  `json:"username"`
	Email              string  `json:"email"`
	CurrentExp         int64   `json:"current_exp"`
	ExpNeeded          int64   `json:"exp_needed"`
	Level              int32   `json:"level"`
	Points             int64   `json:"points"`
	ProfileUrl         string  `json:"profile_url"`
	Challenges         int32   `json:"challenges"`
	Events             int32   `json:"events"`
	Quests             int32   `json:"quests"`
	Treasures          int32   `json:"treasures"`
	LongestStreak      int32   `json:"longest_streak"`
	CurrentStreak      int     `json:"current_Streak"`
	TreeGrown          int32   `json:"tree_grown"`
	CompletedTask      int32   `json:"completed_task"`
	AssignedTask       int32   `json:"assigend_task"`
	TaskCompletionRate string  `json:"task_completion_rate"`
	Badges             []Badge `json:"badges"`
}

type Badge struct {
	Category  string `json:"category"`
	Name      string `json:"name"`
	Frequency int    `json:"frequency"`
}

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

type PutUserEditProfile struct {
	ProfileKey string `json:"profile_key" validate:"required"`
}
