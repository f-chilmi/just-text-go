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

	fmt.Println("message", message)

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

	var roomRes models.RoomResponse

	if roomExisted.IdUser1 == myId {
		roomRes = models.RoomResponse{
			ID:      roomExisted.ID,
			IdUser1: roomExisted.IdUser1,
			IdUser2: roomExisted.IdUser2,
			// User1: user,
			User2: models.UserData{
				Username: user.Username,
				Phone:    user.Phone,
			},
			LastMsg:   roomExisted.LastMsg,
			CreatedAt: roomExisted.CreatedAt,
			UpdatedAt: roomExisted.UpdatedAt,
		}
	} else {
		roomRes = models.RoomResponse{
			ID:      roomExisted.ID,
			IdUser1: roomExisted.IdUser1,
			IdUser2: roomExisted.IdUser2,
			User1: models.UserData{
				Username: user.Username,
				Phone:    user.Phone,
			},
			// User2:     user,
			LastMsg:   roomExisted.LastMsg,
			CreatedAt: roomExisted.CreatedAt,
			UpdatedAt: roomExisted.UpdatedAt,
		}
	}

	// send all the users as response
	responses.JSON(w, http.StatusOK, roomRes)

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

	myId, err := auth.ExtracTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	roomChat, err := roomM.ListRoomByToken(int(myId))

	var roomR []models.RoomResponse

	for _, room := range roomChat {
		fmt.Println("room: ", room)
		var roomRes models.RoomResponse
		if room.IdUser1 == myId {
			roomRes = models.RoomResponse{
				ID:      room.ID,
				IdUser1: room.IdUser1,
				IdUser2: room.IdUser2,
				// User1: user,
				// User2: models.UserData{
				// 	Username: user.Username,
				// 	Phone:    user.Phone,
				// },
				LastMsg:   room.LastMsg,
				CreatedAt: room.CreatedAt,
				UpdatedAt: room.UpdatedAt,
			}
		} else {
			roomRes = models.RoomResponse{
				ID:      room.ID,
				IdUser1: room.IdUser1,
				IdUser2: room.IdUser2,
				// User1: models.UserData{
				// 	Username: user.Username,
				// 	Phone:    user.Phone,
				// },
				// User2:     user,
				LastMsg:   room.LastMsg,
				CreatedAt: room.CreatedAt,
				UpdatedAt: room.UpdatedAt,
			}
		}
		roomR = append(roomR, roomRes)
	}
	fmt.Println("roomR ", roomR)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// send all the users as response
	responses.JSON(w, http.StatusOK, roomR)

}
