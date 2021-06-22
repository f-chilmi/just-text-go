package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/f-chilmi/just-text-go/helpers"
	"github.com/f-chilmi/just-text-go/models"
	"github.com/f-chilmi/just-text-go/responses"
)

type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

type basicRes struct {
	Message string `json:"message,omitempty"`
}

type responseNew struct {
	Message models.Message `json:"message,omitempty"`
}

func createConnection() *sql.DB {
	fmt.Println("create connection")
	// load .env file
	err := godotenv.Load(".env")

	helpers.CheckError("Error loading env files (create connection)", err)

	// initialize db credential
	DbHost := os.Getenv("HOST")
	DbPort := os.Getenv("PORT")
	DbUser := os.Getenv("USER")
	DbPassword := os.Getenv("PASSWORD")
	DbName := os.Getenv("DBNAME")

	DbUrl := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)

	// open the connection
	db, err := sql.Open("postgres", DbUrl)
	if err != nil {
		panic(err)
	}

	// check the connection
	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected to database")

	// return the connection
	return db
}

func HomeController(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Home controller called")
	err := json.NewEncoder(w).Encode("Welcome to the awesome chat app")
	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
	}
}

func CreateUser(w http.ResponseWriter, r *http.Request) {

	fmt.Println("create user called")

	userM := models.User{}

	// set the header to content type x-www-form-urlencoded
	// Allow all origin to handle cors issue
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// create an empty user of type models.User
	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)

	helpers.CheckError("Unable to decode the request body", err)

	insertId, err := userM.InsertUser(user)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	res := response{
		ID:      insertId,
		Message: "user created successfully",
	}

	// send the response
	json.NewEncoder(w).Encode(res)
}

func FindAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	userM := models.User{}

	users, err := userM.GetUsers()

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// send all the users as response
	responses.JSON(w, http.StatusOK, users)
}

func FindById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// get the userid from the request params, key is "id"
	params := mux.Vars(r)

	// convert the id type from string to int
	id, err := strconv.Atoi(params["id"])

	helpers.CheckError("Unable to convert the string into int", err)

	userM := models.User{}

	user, err := userM.GetUser(int64(id))

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// send all the users as response
	responses.JSON(w, http.StatusOK, user)
}

// UpdateUser update user's detail in the postgres db
func UpdateUser(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// get the userid from the request params, key is "id"
	params := mux.Vars(r)

	// convert the id type from string to int
	id, err := strconv.Atoi(params["id"])

	helpers.CheckError("Unable to convert the string into int", err)

	// create an empty user of type models.User
	var user models.User

	// decode the json request to user
	err = json.NewDecoder(r.Body).Decode(&user)

	helpers.CheckError("Unable to decode the request body", err)

	// call update user to update the user
	userM := models.User{}

	updatedRows := userM.UpdateUser(int64(id), user)

	// format the message string
	msg := fmt.Sprintf("Total rows/record affected %v", updatedRows)

	// format the response message
	res := response{
		ID:      int64(id),
		Message: msg,
	}

	if updatedRows < 1 {
		responses.JSON(w, http.StatusBadRequest, res)
		return
	}

	// send the response
	responses.JSON(w, http.StatusOK, res)
}

func NewMsg(w http.ResponseWriter, r *http.Request) {
	// set the header to content type x-www-form-urlencoded
	// Allow all origin to handle cors issue
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// create an empty user of type models.User
	var message models.Message

	err := json.NewDecoder(r.Body).Decode(&message)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// check if rooms existed
	roomsExisted, err := findRoom(message.IdSender, message.IdRecipient)

	fmt.Println("roomsExisted", roomsExisted, err)

	var idRoom int64

	// create new room if no room found
	if err == sql.ErrNoRows {
		fmt.Println("no rooms")
		var dataRoom models.Room
		dataRoom.IdUser1 = message.IdRecipient
		dataRoom.IdUser2 = message.IdSender
		dataRoom.LastMsg = message.Content

		idRoom, err = newRoom(dataRoom)

		if err != nil {
			responses.ERROR(w, http.StatusBadRequest, err)
			return
		}
	}

	idRoom = roomsExisted.ID

	// create new message
	message.IdRoom = idRoom
	newM, err := newMsg(message)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	newM.IdRoom = idRoom
	res := responseNew{
		Message: newM,
	}

	// send all the users as response
	responses.JSON(w, http.StatusOK, res)
	// json.NewEncoder(w).Encode(res)
}

func newMsg(message models.Message) (models.Message, error) {
	// create the db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create the insert query
	// returning userid will return the id of the inserted user
	sqlStatement := `INSERT INTO messages (id_sender, id_recipient, id_room, content) VALUES ($1, $2, $3, $4) RETURNING *;`

	// inserted id will store in this id
	var messages models.Message

	// execute the sql statement
	// scan function will save the inserted id in the id
	row := db.QueryRow(sqlStatement, message.IdSender, message.IdRecipient, message.IdRoom, message.Content)

	err := row.Scan(&messages.ID, &messages.IdSender, &messages.IdRecipient, &messages.Content, &messages.CreatedAt, &messages.UpdatedAt, &message.IdRoom)

	// return the inserted message
	return messages, err
}

func findRoom(idUser1 int64, idUser2 int64) (models.Room, error) {
	// create the db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create the select query
	sqlStatement := `SELECT * FROM rooms WHERE (id_user1=$1 AND id_user2=$2) OR (id_user1=$2 AND id_user2=$1);`

	// inserted id will store in this id
	var room models.Room

	// execute the sql statement
	row := db.QueryRow(sqlStatement, idUser1, idUser2)

	err := row.Scan(&room.ID, &room.IdUser1, &room.IdUser2, &room.LastMsg, &room.CreatedAt, &room.UpdatedAt)

	// return the inserted message
	return room, err
}

func findRoomById(id int64) (models.Room, error) {
	// create the db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create the select query
	sqlStatement := `SELECT * FROM rooms WHERE id=$1;`

	// inserted id will store in this id
	var room models.Room

	// execute the sql statement
	row := db.QueryRow(sqlStatement, id)

	err := row.Scan(&room.ID, &room.IdUser1, &room.IdUser2, &room.LastMsg, &room.CreatedAt, &room.UpdatedAt)

	// return the inserted message
	return room, err
}

func newRoom(r models.Room) (int64, error) {
	// create the db connection
	db := createConnection()

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

func FindByPhone(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// get the userid from the request params, key is "id"
	params := mux.Vars(r)

	// convert the id type from string to int
	phone := params["phone"]

	user, err := getUserByPhone(phone)

	switch err {
	case sql.ErrNoRows:
		res := basicRes{Message: "no user found"}
		responses.JSON(w, http.StatusBadRequest, res)
		return
	case nil:
		break
	default:
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// if user found, now check is there room for id token (me) and user.ID
	roomExisted, err := findRoom(int64(1), int64(user.ID))

	switch err {
	case sql.ErrNoRows:

		newR := models.Room{
			IdUser1: 1,
			IdUser2: user.ID,
			LastMsg: "",
		}

		idRoom, err := newRoom(newR)
		if err != nil {
			responses.ERROR(w, http.StatusBadRequest, err)
			return
		}

		roomExisted, err = findRoomById(idRoom)
		if err != nil {
			responses.ERROR(w, http.StatusBadRequest, err)
			return
		}

	case nil:
		break
	default:
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// send all the users as response
	responses.JSON(w, http.StatusOK, roomExisted)

}

func getUserByPhone(phone string) (models.User, error) {
	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	var user models.User

	// create the select sql query
	sqlStatement := `SELECT * FROM users WHERE phone=$1`

	// execute the sql statement
	row := db.QueryRow(sqlStatement, phone)

	err := row.Scan(&user.ID, &user.Username, &user.Phone, &user.Password, &user.CreatedAt, &user.UpdatedAt)

	// return empty user on error
	return user, err
}

func OpenRoom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var err error

	params := mux.Vars(r)

	idR, err := strconv.Atoi(params["id"])
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	roomChat, err := openRoomChat(idR)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// send all the users as response
	responses.JSON(w, http.StatusOK, roomChat)

}

func ListRoom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var err error

	// for token id
	// params := mux.Vars(r)

	// idR, err := strconv.Atoi(params["id"])
	idR := 1
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	roomChat, err := listRoomByToken(idR)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// send all the users as response
	responses.JSON(w, http.StatusOK, roomChat)

}

func openRoomChat(idR int) ([]models.Message, error) {
	// create the postgres db connection
	db := createConnection()

	var err error

	// close the db connection
	defer db.Close()

	var chats []models.Message

	// create the select sql query
	sqlStatement := `SELECT * FROM messages WHERE id_room=$1`

	// execute the sql statement
	rows, err := db.Query(sqlStatement, idR)

	helpers.CheckError("Unable to execute the query.", err)

	// close the statement
	defer rows.Close()

	// iterate over the rows
	for rows.Next() {
		var chat models.Message

		// unmarshal the row object to user
		err = rows.Scan(&chat.ID, &chat.IdSender, &chat.IdRecipient, &chat.Content, &chat.CreatedAt, &chat.UpdatedAt, &chat.IdRoom)

		helpers.CheckError("Unable to scan the row.", err)

		// append the user in the users slice
		chats = append(chats, chat)

	}

	// return empty user on error
	return chats, err
}

func listRoomByToken(idR int) ([]models.Room, error) {
	// create the postgres db connection
	db := createConnection()

	var err error

	// close the db connection
	defer db.Close()

	var rooms []models.Room

	// create the select sql query
	sqlStatement := `SELECT * FROM rooms WHERE id_user1=$1 OR id_user2=$1`

	// execute the sql statement
	rows, err := db.Query(sqlStatement, idR)

	helpers.CheckError("Unable to execute the query.", err)

	// close the statement
	defer rows.Close()

	// iterate over the rows
	for rows.Next() {
		var room models.Room

		// unmarshal the row object to user
		err = rows.Scan(&room.ID, &room.IdUser1, &room.IdUser2, &room.LastMsg, &room.CreatedAt, &room.UpdatedAt)

		helpers.CheckError("Unable to scan the row.", err)

		// append the user in the users slice
		rooms = append(rooms, room)

	}

	// return empty user on error
	return rooms, err
}
