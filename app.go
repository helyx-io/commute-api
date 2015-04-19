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
    "github.com/jinzhu/gorm"
    "gopkg.in/redis.v2"
    "github.com/helyx-io/gtfs-api/database"
//    "github.com/davecheney/profile"
//    "os/signal"
)



////////////////////////////////////////////////////////////////////////////////////////////////
/// Variables
////////////////////////////////////////////////////////////////////////////////////////////////

var (
    DB              *gorm.DB
    RedisClient     *redis.Client
)

////////////////////////////////////////////////////////////////////////////////////////////////
/// Main Function
////////////////////////////////////////////////////////////////////////////////////////////////

func main() {
    defer Close()

	// Init Runtime
	runtime.GOMAXPROCS(runtime.NumCPU())

//	// Init Profiling
//	// defer profile.Start(profile.MemProfile).Stop()
//	p := profile.Start(profile.CPUProfile)
//
//    c := make(chan os.Signal, 1)
//    signal.Notify(c, os.Interrupt)
//    go func(){
//        p.Stop()
//    }()


    // Init Config
    config := config.Init();

    // Init Logger
    logWriter, err := os.Create(config.LoggerInfos.Path)
    utils.FailOnError(err, fmt.Sprintf("Could not access log"))
    defer logWriter.Close()


    DB, err = database.InitDB(config.ConnectInfos)
    utils.FailOnError(err, fmt.Sprintf("Could not init Database"))

    RedisClient = redis.NewTCPClient(&redis.Options{
        Addr: fmt.Sprintf("%s:%d", config.RedisInfos.Host, config.RedisInfos.Port),
        Password: "", // no password set
        DB:       0,  // use default DB
        PoolSize: 16,
    })

	// Init Router
	router := initRouter(DB, config.ConnectInfos, RedisClient)
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

func Close() {
    if DB != nil {
        defer DB.Close()
    }
}


////////////////////////////////////////////////////////////////////////////////////////////////
/// Router Configuration
////////////////////////////////////////////////////////////////////////////////////////////////

func initRouter(db *gorm.DB, connectConfig *config.DBConnectInfos, redis *redis.Client) *mux.Router {
	r := mux.NewRouter()

	new(controller.IndexController).Init(r.PathPrefix("/").Subrouter())
	new(controller.AgencyController).Init(db, connectConfig, redis, r.PathPrefix("/api/agencies").Subrouter())

	return r
}
