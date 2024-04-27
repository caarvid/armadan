package dto

import (
	"github.com/google/uuid"
)

type hole struct {
	Id    uuid.UUID `json:"id"`
	Par   int32     `json:"par"`
	Nr    int32     `json:"nr"`
	Index int32     `json:"index"`
}

type HoleList []hole
