package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/adlrocha/goxyq/handler"
)

func main() {
	http.HandleFunc("/", handler.sayhelloName) // set router
	http.HandleFunc("/test", handler.bypass)
	fmt.Println("Starting server...")
	err := http.ListenAndServe(":9090", nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
