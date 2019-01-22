package main

import (
	"github.com/adlrocha/goxyq/app"
	"github.com/adlrocha/goxyq/config"
)

// func main() {
// 	// Define handler function for the http server. Common function for every route
// 	proxyHandler := http.HandlerFunc(handler.ProxyRequest)
// 	// If specific routes where required this function may be used.
// 	http.HandleFunc("/test", handler.TestFunction)
// 	fmt.Println("Starting server...")
// 	err := http.ListenAndServe(config.GetConfig().Port, proxyHandler) // set listen port
// 	if err != nil {
// 		log.Fatal("ListenAndServe: ", err)
// 	}
// }

func main() {
	config := config.GetConfig()

	app := &app.App{}
	// app.Initialize(config)
	app.Run(config.Port)
}
