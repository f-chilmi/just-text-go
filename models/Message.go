package models

import (
	"time"

	"github.com/f-chilmi/just-text-go/db"
)

type Message struct {
	ID          int64     `json:"id"`
	IdSender    int64     `json:"id_sender"`
	IdRecipient int64     `json:"id_recipient"`
	IdRoom      int64     `json:"id_room"`
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (m *Message) NewMsg(message Message) (Message, error) {
	// create the db connection
	db := db.CreateConnection()

	// close the db connection
	defer db.Close()

	// create the insert query
	// returning userid will return the id of the inserted user
	sqlStatement := `INSERT INTO messages (id_sender, id_recipient, id_room, content) VALUES ($1, $2, $3, $4) RETURNING *;`

	// inserted id will store in this id
	var messages Message

	// execute the sql statement
	// scan function will save the inserted id in the id
	row := db.QueryRow(sqlStatement, message.IdSender, message.IdRecipient, message.IdRoom, message.Content)

	err := row.Scan(&messages.ID, &messages.IdSender, &messages.IdRecipient, &messages.Content, &messages.CreatedAt, &messages.UpdatedAt, &message.IdRoom)

	// return the inserted message
	return messages, err
}
