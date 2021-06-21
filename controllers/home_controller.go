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

type responseNew struct {
	Message models.Message `json:"message,omitempty"`
}

// func newRespMsg() *models.Message {
// 	return &models.Message{
// 		CreatedAt: time.Now(),
// 		UpdatedAt: time.Now(),
// 	}
// }

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

	insertId := InsertUser(user)

	res := response{
		ID:      insertId,
		Message: "user created successfully",
	}

	// send the response
	json.NewEncoder(w).Encode(res)
}

func InsertUser(user models.User) int64 {
	// create the db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create the insert query
	// returning userid will return the id of the inserted user
	sqlStatement := `INSERT INTO users (username, phone, password) VALUES ($1, $2, $3) RETURNING id;`

	// inserted id will store in this id
	var id int64

	// execute the sql statement
	// scan function will save the inserted id in the id
	err := db.QueryRow(sqlStatement, user.Username, user.Phone, user.Password).Scan(&id)

	helpers.CheckError("Unable to execute the query.", err)

	fmt.Printf("Inserted a single record %v", id)

	// return the inserted id
	return id
}

func FindAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	users, err := getUsers()

	responses.ERROR(w, http.StatusBadRequest, err)

	// send all the users as response
	responses.JSON(w, http.StatusOK, users)

}

// get one user from the DB by its userid
func getUsers() ([]models.User, error) {
	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	var users []models.User

	// create the select sql query
	sqlStatement := `SELECT * FROM users`

	// execute the sql statement
	rows, err := db.Query(sqlStatement)

	helpers.CheckError("Unable to execute the query.", err)

	// close the statement
	defer rows.Close()

	// iterate over the rows
	for rows.Next() {
		var user models.User

		// unmarshal the row object to user
		err := rows.Scan(&user.ID, &user.Username, &user.Phone, &user.Password, &user.CreatedAt, &user.UpdatedAt)

		helpers.CheckError("Unable to scan the row.", err)

		// append the user in the users slice
		users = append(users, user)

	}

	// return empty user on error
	return users, err
}

func FindById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// get the userid from the request params, key is "id"
	params := mux.Vars(r)

	// convert the id type from string to int
	id, err := strconv.Atoi(params["id"])

	helpers.CheckError("Unable to convert the string into int", err)

	// call the getUser function with user id to retrieve a single user
	user, err := getUser(int64(id))

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// send all the users as response
	responses.JSON(w, http.StatusOK, user)

}

// get one user from the DB by its userid
func getUser(id int64) (models.User, error) {
	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create a user of models.User type
	var user models.User

	// create the select sql query
	sqlStatement := `SELECT * FROM users WHERE id=$1`

	// execute the sql statement
	row := db.QueryRow(sqlStatement, id)

	// unmarshal the row object to user
	err := row.Scan(&user.ID, &user.Username, &user.Phone, &user.Password, &user.CreatedAt, &user.UpdatedAt)

	// switch err {
	// case sql.ErrNoRows:
	// 	fmt.Println("No rows were returned!", err)
	// 	return user, nil
	// case nil:
	// 	return user, nil
	// default:
	// 	log.Fatalf("Unable to scan the row. %v", err)
	// }

	// return empty user on error
	return user, err
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
	updatedRows := updateUser(int64(id), user)

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

// update user in the DB
func updateUser(id int64, user models.User) int64 {

	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create the update sql query
	sqlStatement := `UPDATE users SET username=$2, phone=$3, updated_at=CURRENT_TIMESTAMP WHERE id=$1`

	// execute the sql statement
	res, err := db.Exec(sqlStatement, id, user.Username, user.Phone)

	helpers.CheckError("Unable to execute the query", err)

	// check how many rows affected
	rowsAffected, err := res.RowsAffected()

	helpers.CheckError("Error while cheking the affected rows", err)

	fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected
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
