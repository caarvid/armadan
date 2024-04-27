package dto

import "github.com/google/uuid"

type Tee struct {
	Id    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Slope int32     `json:"slope"`
	Cr    float32   `json:"cr"`
}

type TeeList []Tee
