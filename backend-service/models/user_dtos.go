package models

type ResponseGetUserProfileStatistic struct {
	Id                     int     `json:"id"`
	Name                   string  `json:"name"`
	Username               string  `json:"username"`
	Email                  string  `json:"email"`
	CurrentExp             int64   `json:"current_exp"`
	ExpNeeded              int64   `json:"exp_needed"`
	Level                  int32   `json:"level"`
	Points                 int64   `json:"points"`
	ProfileUrl             string  `json:"profile_url"`
	Challenges             int32   `json:"challenges"`
	Events                 int32   `json:"events"`
	Quests                 int32   `json:"quests"`
	Treasures              int32   `json:"treasures"`
	LongestStreak          int32   `json:"longest_streak"`
	CurrentStreak          int     `json:"current_Streak"`
	TreeGrown              int32   `json:"tree_grown"`
	NeededExpPreviousLevel int     `json:"needed_exp_previous_level"`
	CompletedTask          int32   `json:"completed_task"`
	AssignedTask           int32   `json:"assigend_task"`
	TaskCompletionRate     string  `json:"task_completion_rate"`
	Rank                   int     `json:"rank"`
	Badges                 []Badge `json:"badges"`
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
	Filename    string `json:"filename" validate:"required"`
	ContentType string `json:"content_type"`
}

type ResponseUser struct {
	ID         int    `json:"id"`
	Level      int    `json:"level"`
	ProfileUrl string `json:"profile_url"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Streak     int    `json:"streak"`
}

type RequestStatistics struct {
	Challenges    int `json:"challenges"`
	Events        int `json:"events"`
	Quests        int `json:"quests"`
	Treasures     int `json:"treasures"`
	LongestStreak int `json:"longest_streak"`
}
