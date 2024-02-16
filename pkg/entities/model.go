package entities

import "time"

type ModelID struct {
	Id uint `gorm:"primaryKey" json:"id"`
}

type ModelTimeStamps struct {
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
