package main

import (
	"fmt"
	"go-crud/router"
	"log"
	"net/http"
)



func main() {
	r := router.Router()
	fmt.Println("starting server on port 8080...")

	log.Fatal(http.ListenAndServe(":8080", r))
}