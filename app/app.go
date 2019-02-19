package app

import (
	"net/http"

	"github.com/adlrocha/goxyq/config"
	"github.com/adlrocha/goxyq/handler"
	"github.com/adlrocha/goxyq/log"
	"github.com/gorilla/mux"
)

// App has router and db instances
type App struct {
	Router *mux.Router
}

// Set all required routers and handlers in MUX
func (a *App) setRouters() {
	// Main Proxy handler
	a.Router.PathPrefix(config.GetConfig().ProxyPathPrefix).Handler(http.HandlerFunc(handler.ProxyRequest))
	// Additional routing paths
	a.Router.HandleFunc("/alive", handler.AliveFunction)              // Check if the proxy is alive.
	a.Router.HandleFunc("/queue/{queueID}", handler.GetQueue)         // Get the status of a queue
	a.Router.HandleFunc("/queue/{queueID}/empty", handler.EmptyQueue) // Empty jobs of a queue

}

// Run the app on it's router
func (a *App) Run(host string) {
	a.Router = mux.NewRouter()
	a.setRouters()
	log.Infof("[APP] Running server at port %v", host)
	log.Fatalf("%v", http.ListenAndServe(host, a.Router))
}
