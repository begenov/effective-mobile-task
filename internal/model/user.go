package model

import "time"

type User struct {
	ID          int64        `json:"id"`
	Name        string       `json:"name"`
	Surname     string       `json:"surname"`
	Patronymic  *string      `json:"patronymic"`
	Age         *int32       `json:"age"`
	Gender      *int         `json:"gender"`
	Nationality *Nationality `json:"nationality"`
	CreatedAt   *time.Time   `json:"created_at"`
	UpdatedAt   *time.Time   `json:"updated_at"`
}

type Nationality struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
