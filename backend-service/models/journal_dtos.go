package models

type PostLogAppend struct {
	Text      string `json:"text" validate:"required"`
	IsSystem  bool   `json:"is_system" validate:"boolean"`
	IsPrivate bool   `json:"is_private" validate:"boolean"`
}

type ResponseGetLogs struct {
	Text      string `json:"text"`
	IsSystem  bool   `json:"is_system"`
	IsPrivate bool   `json:"is_private,"`
	CreatedAt string `json:"created_at"`
}
