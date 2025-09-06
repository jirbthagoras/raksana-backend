package models

type ResponsePacketDetail struct {
	PacketID           int64                       `json:"packet_id"`
	Username           string                      `json:"username"`
	PacketName         string                      `json:"packet_name"`
	Target             string                      `json:"target"`
	Description        string                      `json:"description"`
	CompletedTask      int32                       `json:"completed_task"`
	ExpectedTask       int32                       `json:"expected_task"`
	AssignedTask       int32                       `json:"assigned_task"`
	TaskCompletionRate string                      `json:"task_completion_rate"`
	TaskPerDay         int32                       `json:"task_per_day"`
	Completed          bool                        `json:"completed"`
	CreatedAt          string                      `json:"created_at"`
	Habits             []ResponsePacketDetailHabit `json:"habits"`
}

type ResponsePacketDetailHabit struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Difficulty  string `json:"difficulty"`
	Locked      bool   `json:"locked"`
	ExpGain     int32  `json:"exp_gain"`
}

type ResponseGetPacket struct {
	Id             int32  `json:"id"`
	Name           string `json:"name"`
	Target         string `json:"target"`
	Description    string `json:"description"`
	CompletedTask  int32  `json:"completed_task"`
	ExpectedTask   int32  `json:"expected_task"`
	AssignedTask   int32  `json:"assigned_task"`
	CompletionRate string `json:"completion_rate"`
	TaskPerDay     int32  `json:"task_per_day"`
	Completed      bool   `json:"completed"`
	CreatedAt      string `json:"created_at"`
}

type EcoachCreatePacketResponse struct {
	Name         string `json:"name"`
	ExpectedTask int    `json:"expected_task"`
	TaskPerDay   int    `json:"task_per_day"`
	Habits       []EcoachHabitResponse
}

type EcoachHabitResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Difficulty  string `json:"difficulty"`
}

type PostPacketCreate struct {
	Target      string `json:"target" validate:"required"`
	Description string `json:"description" validate:"required"`
}
