package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/adlrocha/goxyq/handler"
)

func main() {
	http.HandleFunc("/", handler.SayhelloName) // set router
	http.HandleFunc("/test", handler.Bypass)
	fmt.Println("Starting server...")
	err := http.ListenAndServe(":9090", nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
