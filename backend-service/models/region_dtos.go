package models

type ResponseRegion struct {
	Id         int     `json:"id"`
	Name       string  `json:"name"`
	Location   string  `json:"location"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	TreeAmount int     `json:"tree_amount"`
}
