package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/f-chilmi/just-text-go/router"
)

func main() {
	r := router.Router()

	fmt.Println("Starting server on port 8080")

	log.Fatal(http.ListenAndServe(":8080", r))
}
