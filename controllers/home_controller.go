package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/f-chilmi/just-text-go/models"
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
