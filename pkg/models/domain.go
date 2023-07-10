package models

import "time"

type Resource struct {
	Id      string `json:"id" gorm:"primaryKey"`
	Parent  string `json:"parent"`
	Name    string `json:"name"`
	Version string `json:"version"`

	// json
	Data string `json:"data"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Domain struct {
	Id     string `json:"id" gorm:"primaryKey"`
	Name   string `json:"name"`
	Parent string `json:"parent"`

	Resources []Resource `json:"resources" gorm:"foreignKey:Parent"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
