package app

import (
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"gitlab.com/api/app/handler"
	"gitlab.com/api/config"
)

// App has router
type App struct {
	Config *config.Config
	Router *mux.Router
}

// App initialize with predefined configuration
func (a *App) Initialize(config *config.Config) {
	a.Config = config
	a.Config.Logger.Formatter = &log.TextFormatter{
		FullTimestamp: true,
	}
	a.Config.Logger.Info("Initializing...")
	a.Router = mux.NewRouter()
	a.setRouters()
}

// Set all required routers
func (a *App) setRouters() {
	a.Config.Logger.Info("Setting Routers...")
	a.Get("/properties", a.GetAllProperties)
}

// Wrap the router for GET method
func (a *App) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("GET")
}

// Handlers to manage Employee Data
func (a *App) GetAllProperties(w http.ResponseWriter, r *http.Request) {
	a.Config.Logger.WithFields(log.Fields{
		"URL":    r.URL,
		"header": r.Header,
	}).Info("Requesting all Properties")
	handler.GetAllProperties(a.Config, w, r)
}

// Run the app on it's router
func (a *App) Run(host string) {
	a.Config.Logger.Info("Listening to the port", host)
	log.Fatal(http.ListenAndServe(host, a.Router))
}
