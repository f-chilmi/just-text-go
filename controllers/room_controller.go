package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/f-chilmi/just-text-go/auth"
	"github.com/f-chilmi/just-text-go/models"
	"github.com/f-chilmi/just-text-go/responses"
	"github.com/gorilla/mux"
)

func NewMsg(w http.ResponseWriter, r *http.Request) {
	// set the header to content type x-www-form-urlencoded
	// Allow all origin to handle cors issue
	// w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	// w.Header().Set("Access-Control-Allow-Methods", "POST")
	// w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

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

func SendMsg(w http.ResponseWriter, r *http.Request) {
	messageM := models.Message{}
	roomM := models.Room{}

	params := mux.Vars(r)
	idRoom, err := (strconv.Atoi(params["id"]))
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	myId, err := auth.ExtracTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// check if rooms existed
	_, err = roomM.FindRoomById(int64(idRoom))
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	var message models.Message
	err = json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	message.IdSender = myId
	message.IdRoom = int64(idRoom)

	newM, err := messageM.NewMsg(message)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	message.ID = newM.ID
	message.CreatedAt = newM.CreatedAt
	message.UpdatedAt = newM.UpdatedAt
	res := responseNew{
		Message: message,
	}

	// send all the users as response
	responses.JSON(w, http.StatusOK, res)
}

func FindRoomByPhone(w http.ResponseWriter, r *http.Request) {
	roomM := models.Room{}
	userM := models.User{}

	params := mux.Vars(r)
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

	myId, err := auth.ExtracTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// if user found, then check is there room for id token (me) and user.ID
	roomExisted, err := roomM.FindRoom(int64(myId), int64(user.ID))

	switch err {
	case sql.ErrNoRows:

		newR := models.Room{
			IdUser1: myId,
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

	responses.JSON(w, http.StatusOK, roomChat)
}

func ListRoom(w http.ResponseWriter, r *http.Request) {

	var err error

	roomM := models.Room{}

	// for token id
	// params := mux.Vars(r)

	idR, err := auth.ExtracTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	roomChat, err := roomM.ListRoomByToken(int(idR))

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// send all the users as response
	responses.JSON(w, http.StatusOK, roomChat)

}
