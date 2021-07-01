package models

import (
	"time"

	"github.com/f-chilmi/just-text-go/db"
	"github.com/f-chilmi/just-text-go/helpers"
)

type RoomDb struct {
	ID        int64     `json:"id"`
	IdUser1   int64     `json:"id_user1"`
	IdUser2   int64     `json:"id_user2"`
	LastMsg   string    `json:"last_msg"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RoomList struct {
	ID        int64     `json:"id"`
	IdUser1   int64     `json:"id_user1"`
	Username1 string    `json:"username1"`
	Phone1    string    `json:"phone1"`
	IdUser2   int64     `json:"id_user2"`
	Username2 string    `json:"username2"`
	Phone2    string    `json:"phone2"`
	LastMsg   string    `json:"last_msg"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type Room struct {
	ID             int64     `json:"id"`
	IdRecipient    int64     `json:"id_recipient"`
	UnameRecipient string    `json:"uname_recipient"`
	PhoneRecipient string    `json:"phone_recipient"`
	LastMsg        string    `json:"last_msg"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type RoomResponse struct {
	ID          int64     `json:"id"`
	IdRecipient int64     `json:"id_recipient"`
	Recipient   UserData  `json:"recipient"`
	LastMsg     string    `json:"last_msg"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (r *Room) FindRoom(myId int64, idUser2 int64) (Room, error) {
	// create the db connection
	db := db.CreateConnection()

	// close the db connection
	defer db.Close()

	// create the select query
	// sqlStatement := `SELECT * FROM rooms WHERE (id_user1=$1 AND id_user2=$2) OR (id_user1=$2 AND id_user2=$1);`
	sqlStatement := `
		SELECT 
			rooms.id, 
			id_user1, 
			a.username as username1, 
			a.phone as phone1, 
			id_user2, 
			b.username as username2, 
			b.phone as phone2, 
			last_msg, 
			rooms.created_at, 
			rooms.updated_at from rooms 
		INNER JOIN users a on rooms.id_user1 = a.id
		INNER JOIN users b on rooms.id_user2 = b.id
		WHERE (id_user1=$1 AND id_user2=$2) OR (id_user1=$2 AND id_user2=$1)`

	var room RoomList
	var newR Room
	// execute the sql statement
	row := db.QueryRow(sqlStatement, myId, idUser2)
	err := row.Scan(
		&room.ID,
		&room.IdUser1,
		&room.Username1,
		&room.Phone1,
		&room.IdUser2,
		&room.Username2,
		&room.Phone2,
		&room.LastMsg,
		&room.CreatedAt,
		&room.UpdatedAt,
	)
	if err != nil {
		return newR, err
	}

	newR = Room{
		ID:        room.ID,
		LastMsg:   room.LastMsg,
		CreatedAt: room.CreatedAt,
		UpdatedAt: room.UpdatedAt,
	}
	if room.IdUser1 == myId {
		newR.IdRecipient = room.IdUser2
		newR.UnameRecipient = room.Username2
		newR.PhoneRecipient = room.Phone2
	} else {
		newR.IdRecipient = room.IdUser1
		newR.UnameRecipient = room.Username1
		newR.PhoneRecipient = room.Phone1
	}
	// return the inserted message
	return newR, err
}

func (r *Room) FindRoomById(id int64) (RoomList, error) {
	// create the db connection
	db := db.CreateConnection()

	// close the db connection
	defer db.Close()

	// create the select query
	// sqlStatement := `SELECT * FROM rooms WHERE id=$1;`
	sqlStatement := `
		SELECT 
			rooms.id, 
			id_user1, 
			a.username as username1, 
			a.phone as phone1, 
			id_user2, 
			b.username as username2, 
			b.phone as phone2, 
			last_msg, 
			rooms.created_at, 
			rooms.updated_at from rooms 
		INNER JOIN users a on rooms.id_user1 = a.id
		INNER JOIN users b on rooms.id_user2 = b.id
		WHERE rooms.id=$1`

	// inserted id will store in this id
	var room RoomList

	// execute the sql statement
	row := db.QueryRow(sqlStatement, id)
	err := row.Scan(
		&room.ID,
		&room.IdUser1,
		&room.Username1,
		&room.Phone1,
		&room.IdUser2,
		&room.Username2,
		&room.Phone2,
		&room.LastMsg,
		&room.CreatedAt,
		&room.UpdatedAt,
	)

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
	// sqlStatement := `SELECT * FROM rooms WHERE id_user1=$1 OR id_user2=$1`
	sqlStatement := `
	SELECT 
		rooms.id, 
		id_user1, 
		a.username as username1, 
		a.phone as phone1, 
		id_user2, 
		b.username as username2, 
		b.phone as phone2, 
		last_msg, 
		rooms.created_at, 
		rooms.updated_at from rooms 
	INNER JOIN users a on rooms.id_user1 = a.id
	INNER JOIN users b on rooms.id_user2 = b.id
	WHERE id_user1=$1 OR id_user2=$1`

	// execute the sql statement
	rows, err := db.Query(sqlStatement, idR)

	helpers.CheckError("Unable to execute the query.", err)

	// close the statement
	defer rows.Close()

	// iterate over the rows
	for rows.Next() {
		var room RoomList
		// unmarshal the row object to user
		err = rows.Scan(
			&room.ID,
			&room.IdUser1,
			&room.Username1,
			&room.Phone1,
			&room.IdUser2,
			&room.Username2,
			&room.Phone2,
			&room.LastMsg,
			&room.CreatedAt,
			&room.UpdatedAt,
		)
		if err != nil {
			return rooms, err
		}

		newR := Room{
			ID:        room.ID,
			LastMsg:   room.LastMsg,
			CreatedAt: room.CreatedAt,
			UpdatedAt: room.UpdatedAt,
		}
		if room.IdUser1 == int64(idR) {
			newR.IdRecipient = room.IdUser2
			newR.UnameRecipient = room.Username2
			newR.PhoneRecipient = room.Phone2
		} else {
			newR.IdRecipient = room.IdUser1
			newR.UnameRecipient = room.Username1
			newR.PhoneRecipient = room.Phone1
		}

		// append the user in the users slice
		rooms = append(rooms, newR)

	}

	// return empty user on error
	return rooms, err
}

func (ru *Room) NewRoom(r RoomDb) (int64, error) {
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

func (r *Room) UpdateLastMsg(idRoom int64, msg string) error {

	// create the postgres db connection
	db := db.CreateConnection()

	// close the db connection
	defer db.Close()

	// create the update sql query
	sqlStatement := `UPDATE rooms SET last_msg=$2, updated_at=CURRENT_TIMESTAMP WHERE id=$1`

	// execute the sql statement
	res, err := db.Exec(sqlStatement, idRoom, msg)
	if err != nil {
		return err
	}

	// check how many rows affected
	_, err = res.RowsAffected()

	return err
}
