package models

import "time"

type Message struct {
	ID          int64     `json:"id"`
	IdSender    int64     `json:"id_sender"`
	IdRecipient int64     `json:"id_recipient"`
	IdRoom      int64     `json:"id_room"`
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
