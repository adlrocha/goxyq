package app

import (
	"net/http"

	"github.com/adlrocha/goxyq/handler"
	"github.com/adlrocha/goxyq/log"
	"github.com/gorilla/mux"
)

// App has router and db instances
type App struct {
	Router *mux.Router
}

// Set all required routers and handlers in MUX
// TODO: Configure this through config file.
func (a *App) setRouters() {
	// Main handler
	a.Router.PathPrefix("/api/v1/").Handler(http.HandlerFunc(handler.ProxyRequest))
	// Additional routing paths
	a.Router.HandleFunc("/alive", handler.AliveFunction) // Check if the proxy is alive.
}

// Run the app on it's router
func (a *App) Run(host string) {
	a.Router = mux.NewRouter()
	a.setRouters()
	log.Infof("[APP] Running server at port %v", host)
	log.Fatalf("%v", http.ListenAndServe(host, a.Router))
}
