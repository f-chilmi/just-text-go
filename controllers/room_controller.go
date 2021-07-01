package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/f-chilmi/just-text-go/auth"
	"github.com/f-chilmi/just-text-go/models"
	"github.com/f-chilmi/just-text-go/responses"
	"github.com/gorilla/mux"
)

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

	err = roomM.UpdateLastMsg(int64(idRoom), message.Content)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

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

		newR := models.RoomDb{
			IdUser1: myId,
			IdUser2: user.ID,
			LastMsg: "",
		}

		idRoom, err := roomM.NewRoom(newR)
		if err != nil {
			responses.ERROR(w, http.StatusBadRequest, err)
			return
		}

		var roomF models.RoomList
		roomF, err = roomM.FindRoomById(idRoom)
		if err != nil {
			responses.ERROR(w, http.StatusBadRequest, err)
			return
		}

		roomExisted.ID = roomF.ID
		roomExisted.LastMsg = roomF.LastMsg
		roomExisted.CreatedAt = roomF.CreatedAt
		roomExisted.UpdatedAt = roomF.UpdatedAt
		if roomF.IdUser1 == myId {
			roomExisted.IdRecipient = roomF.IdUser2
			roomExisted.UnameRecipient = roomF.Username2
			roomExisted.PhoneRecipient = roomF.Phone2
		} else {
			roomExisted.IdRecipient = roomF.IdUser1
			roomExisted.UnameRecipient = roomF.Username1
			roomExisted.PhoneRecipient = roomF.Phone1
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

	myId, err := auth.ExtracTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	roomChat, err := roomM.ListRoomByToken(int(myId))
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// send all the users as response
	responses.JSON(w, http.StatusOK, roomChat)

}
