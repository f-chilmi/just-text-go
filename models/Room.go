package models

import (
	"time"

	"github.com/f-chilmi/just-text-go/db"
	"github.com/f-chilmi/just-text-go/helpers"
)

type Room struct {
	ID        int64     `json:"id"`
	IdUser1   int64     `json:"id_user1"`
	IdUser2   int64     `json:"id_user2"`
	LastMsg   string    `json:"last_msg"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (r *Room) FindRoom(idUser1 int64, idUser2 int64) (Room, error) {
	// create the db connection
	db := db.CreateConnection()

	// close the db connection
	defer db.Close()

	// create the select query
	sqlStatement := `SELECT * FROM rooms WHERE (id_user1=$1 AND id_user2=$2) OR (id_user1=$2 AND id_user2=$1);`

	// inserted id will store in this id
	var room Room

	// execute the sql statement
	row := db.QueryRow(sqlStatement, idUser1, idUser2)

	err := row.Scan(&room.ID, &room.IdUser1, &room.IdUser2, &room.LastMsg, &room.CreatedAt, &room.UpdatedAt)

	// return the inserted message
	return room, err
}

func (r *Room) FindRoomById(id int64) (Room, error) {
	// create the db connection
	db := db.CreateConnection()

	// close the db connection
	defer db.Close()

	// create the select query
	sqlStatement := `SELECT * FROM rooms WHERE id=$1;`

	// inserted id will store in this id
	var room Room

	// execute the sql statement
	row := db.QueryRow(sqlStatement, id)

	err := row.Scan(&room.ID, &room.IdUser1, &room.IdUser2, &room.LastMsg, &room.CreatedAt, &room.UpdatedAt)

	// return the inserted message
	return room, err
}

func (r *Room) OpenRoomChat(idR int) ([]Message, error) {
	// create the postgres db connection
	db := db.CreateConnection()

	var err error

	// close the db connection
	defer db.Close()

	var chats []Message

	// create the select sql query
	sqlStatement := `SELECT * FROM messages WHERE id_room=$1`

	// execute the sql statement
	rows, err := db.Query(sqlStatement, idR)

	helpers.CheckError("Unable to execute the query.", err)

	// close the statement
	defer rows.Close()

	// iterate over the rows
	for rows.Next() {
		var chat Message

		// unmarshal the row object to user
		err = rows.Scan(&chat.ID, &chat.IdSender, &chat.IdRecipient, &chat.Content, &chat.CreatedAt, &chat.UpdatedAt, &chat.IdRoom)

		helpers.CheckError("Unable to scan the row.", err)

		// append the user in the users slice
		chats = append(chats, chat)

	}

	// return empty user on error
	return chats, err
}

func (r *Room) ListRoomByToken(idR int) ([]Room, error) {
	// create the postgres db connection
	db := db.CreateConnection()

	var err error

	// close the db connection
	defer db.Close()

	var rooms []Room

	// create the select sql query
	sqlStatement := `SELECT * FROM rooms WHERE id_user1=$1 OR id_user2=$1`

	// execute the sql statement
	rows, err := db.Query(sqlStatement, idR)

	helpers.CheckError("Unable to execute the query.", err)

	// close the statement
	defer rows.Close()

	// iterate over the rows
	for rows.Next() {
		var room Room

		// unmarshal the row object to user
		err = rows.Scan(&room.ID, &room.IdUser1, &room.IdUser2, &room.LastMsg, &room.CreatedAt, &room.UpdatedAt)

		helpers.CheckError("Unable to scan the row.", err)

		// append the user in the users slice
		rooms = append(rooms, room)

	}

	// return empty user on error
	return rooms, err
}

func (ru *Room) NewRoom(r Room) (int64, error) {
	// create the db connection
	db := db.CreateConnection()

	// close the db connection
	defer db.Close()

	// create the insert query
	// returning userid will return the id of the inserted user
	sqlStatement := `INSERT INTO rooms (id_user1, id_user2, last_msg) VALUES ($1, $2, $3) RETURNING id;`

	// inserted id will store in this id
	var idRoom int64

	// execute the sql statement
	// scan function will save the inserted id in the id
	err := db.QueryRow(sqlStatement, r.IdUser1, r.IdUser2, r.LastMsg).Scan(&idRoom)

	// return the inserted message
	return idRoom, err
}
