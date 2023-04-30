package models

import "time"

type Key struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`

	UserId string `json:"-"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
