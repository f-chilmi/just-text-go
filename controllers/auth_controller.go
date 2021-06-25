package controllers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/f-chilmi/just-text-go/auth"
	"github.com/f-chilmi/just-text-go/models"
	"github.com/f-chilmi/just-text-go/responses"
)

func Login(w http.ResponseWriter, r *http.Request) {
	userM := models.User{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	err = json.Unmarshal(body, &userM)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	userM.Prepare()
	err = userM.Validate("login")
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	// token, err := userM.Login(userM.Phone, userM.Password)
	userExisted, err := userM.GetUserByPhone(userM.Phone)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	err = auth.CheckPasswordHash(userM.Password, userExisted.Password)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
	}

	response, validToken, err := auth.GenerateJWT(userM.ID, userM.Username, userM.Phone)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
	}

	res := models.ResLoginWithToken{
		ID:       userExisted.ID,
		Phone:    response.Phone,
		Username: userExisted.Username,
		Exp:      response.Exp,
		Token:    validToken,
	}

	responses.JSON(w, http.StatusOK, res)
}

func Register(w http.ResponseWriter, r *http.Request) {
	userM := models.User{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	err = json.Unmarshal(body, &userM)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	userM.Prepare()
	err = userM.Validate("register")
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	_, err = userM.GetUserByPhone(userM.Phone)

	switch err {
	// if no user found, so user can create new (register)
	case sql.ErrNoRows:
		newP, err := auth.GeneratehashPassword(userM.Password)

		if err != nil {
			responses.ERROR(w, http.StatusBadRequest, err)
		}

		newU := models.User{
			Username: userM.Username,
			Phone:    userM.Phone,
			Password: newP,
		}
		_, err = userM.InsertUser(newU)
		if err != nil {
			responses.ERROR(w, http.StatusBadRequest, err)
		}

		res := basicRes{Message: "user created successfully"}
		responses.JSON(w, http.StatusOK, res)

	// if user found
	case nil:
		responses.ERROR(w, http.StatusBadRequest, errors.New("user already exist"))
		return

	default:
		break
	}

}
