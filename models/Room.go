package models

import "time"

type Room struct {
	ID        int64     `json:"id"`
	IdUser1   int64     `json:"id_user1"`
	IdUser2   int64     `json:"id_user2"`
	LastMsg   string    `json:"last_msg"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
