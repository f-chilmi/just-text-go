package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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

	messageM := models.Message{}
	roomM := models.Room{}

	// create an empty user of type models.User
	var message models.Message

	err := json.NewDecoder(r.Body).Decode(&message)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// check if rooms existed
	roomsExisted, err := roomM.FindRoom(message.IdSender, message.IdRecipient)

	fmt.Println("roomsExisted", roomsExisted, err)

	var idRoom int64

	// create new room if no room found
	if err == sql.ErrNoRows {
		fmt.Println("no rooms")
		var dataRoom models.Room
		dataRoom.IdUser1 = message.IdRecipient
		dataRoom.IdUser2 = message.IdSender
		dataRoom.LastMsg = message.Content

		idRoom, err = roomM.NewRoom(dataRoom)

		if err != nil {
			responses.ERROR(w, http.StatusBadRequest, err)
			return
		}
	}

	idRoom = roomsExisted.ID

	// create new message
	message.IdRoom = idRoom
	newM, err := messageM.NewMsg(message)

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

func FindByPhone(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	roomM := models.Room{}
	userM := models.User{}

	// get the userid from the request params, key is "id"
	params := mux.Vars(r)

	// convert the id type from string to int
	phone := params["phone"]

	user, err := userM.GetUserByPhone(phone)

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
	roomExisted, err := roomM.FindRoom(int64(1), int64(user.ID))

	switch err {
	case sql.ErrNoRows:

		newR := models.Room{
			IdUser1: 1,
			IdUser2: user.ID,
			LastMsg: "",
		}

		idRoom, err := roomM.NewRoom(newR)
		if err != nil {
			responses.ERROR(w, http.StatusBadRequest, err)
			return
		}

		roomExisted, err = roomM.FindRoomById(idRoom)
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

func OpenRoom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var err error

	roomM := models.Room{}

	params := mux.Vars(r)

	idR, err := strconv.Atoi(params["id"])
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	roomChat, err := roomM.OpenRoomChat(idR)

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

	roomM := models.Room{}

	// for token id
	// params := mux.Vars(r)

	// idR, err := strconv.Atoi(params["id"])
	idR := 1
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	roomChat, err := roomM.ListRoomByToken(idR)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// send all the users as response
	responses.JSON(w, http.StatusOK, roomChat)

}
