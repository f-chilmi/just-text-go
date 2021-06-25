package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/f-chilmi/just-text-go/helpers"
	"github.com/f-chilmi/just-text-go/models"
	"github.com/f-chilmi/just-text-go/responses"
	"github.com/gorilla/mux"
)

func FindAll(w http.ResponseWriter, r *http.Request) {
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
	// get the userid from the request params, key is "id"
	params := mux.Vars(r)

	// convert the id type from string to int
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	userM := models.User{}

	user, err := userM.GetUser(int64(id))
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// send all the users as response
	responses.JSON(w, http.StatusOK, user)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
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
