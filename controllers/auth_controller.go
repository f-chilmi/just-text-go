package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/f-chilmi/just-text-go/models"
	"github.com/f-chilmi/just-text-go/responses"
)

func Login(w http.ResponseWriter, r *http.Request) {
	userM := models.User{}

	// set the header to content type x-www-form-urlencoded
	// Allow all origin to handle cors issue
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

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

	fmt.Println("aman")

	userM.Prepare()
	err = userM.Validate("login")
	fmt.Println("validate")
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	token, err := userM.Login(userM.Phone, userM.Password)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	fmt.Println("login")
	responses.JSON(w, http.StatusOK, token)
}

func Register(w http.ResponseWriter, r *http.Request) {
	userM := models.User{}

	// set the header to content type x-www-form-urlencoded
	// Allow all origin to handle cors issue
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

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

	fmt.Println(userM)

	_, err = userM.Register(userM.Username, userM.Phone, userM.Password)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	res := basicRes{Message: "user created successfully"}
	responses.JSON(w, http.StatusOK, res)
}
