package main

////////////////////////////////////////////////////////////////////////////////////////////////
/// Imports
////////////////////////////////////////////////////////////////////////////////////////////////

import (
	"fmt"
	appHandlers "github.com/helyx-io/gtfs-api/handlers"
	"github.com/helyx-io/gtfs-api/config"
	"github.com/helyx-io/gtfs-api/controller"
	"github.com/helyx-io/gtfs-api/utils"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"

	"log"
	"net/http"
	"os"
	"runtime"
)


////////////////////////////////////////////////////////////////////////////////////////////////
/// Main Function
////////////////////////////////////////////////////////////////////////////////////////////////

func main() {

	// Init Runtime
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Init Profiling
	//	defer profile.Start(profile.MemProfile).Stop()
	//	defer profile.Start(profile.CPUProfile).Stop()

	// Init Logger
	logWriter, err := os.Create("/var/log/gtfs-api/access.log")
	utils.FailOnError(err, fmt.Sprintf("Could not access log"))
	defer logWriter.Close()

	// Init Config
	err = config.Init();
	utils.FailOnError(err, fmt.Sprintf("Could not init Configuration"))
	defer config.Close()

	// Init Router
	router := initRouter()
	http.Handle("/", router)

	handlerChain := alice.New(
		appHandlers.LoggingHandler(logWriter),
//		appHandlers.ThrottleHandler,
//		appHandlers.TimeoutHandler,
	).Then(router)

	// Init HTTP Server

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", config.Http.Port),
		Handler: handlerChain,
	}

	log.Println(fmt.Sprintf("Listening on port '%d' ...", config.Http.Port))

	err = server.ListenAndServe()
    utils.FailOnError(err, fmt.Sprintf("Could not listen and server"))
}



////////////////////////////////////////////////////////////////////////////////////////////////
/// Router Configuration
////////////////////////////////////////////////////////////////////////////////////////////////

func initRouter() *mux.Router {
	r := mux.NewRouter()

	new(controller.IndexController).Init(r.PathPrefix("/").Subrouter())
	new(controller.AgencyController).Init(r.PathPrefix("/api/agencies").Subrouter())

	return r
}
