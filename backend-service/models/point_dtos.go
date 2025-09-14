package models

type RequestPostConvert struct {
	Amount   int `json:"amount" validate:"required"`
	RegionId int `json:"region_id" validate:"required"`
}
