package dto

import "github.com/google/uuid"

type WeekCourse struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type WeekTee struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
