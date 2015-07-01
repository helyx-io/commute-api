package controllers

////////////////////////////////////////////////////////////////////////////////////////////////
/// Imports
////////////////////////////////////////////////////////////////////////////////////////////////

import (
	"log"
    "time"
    "sort"
    "net/http"
    "encoding/json"
    "github.com/gorilla/mux"
    "github.com/fatih/stopwatch"
    "github.com/helyx-io/commute-api/utils"
    "github.com/helyx-io/commute-api/models"
    "github.com/helyx-io/commute-api/services"
)


////////////////////////////////////////////////////////////////////////////////////////////////
/// Structures
////////////////////////////////////////////////////////////////////////////////////////////////

type StopController struct {
    ss *services.StopService
}


////////////////////////////////////////////////////////////////////////////////////////////////
/// Constructor
////////////////////////////////////////////////////////////////////////////////////////////////

func NewStopController(ss *services.StopService, r *mux.Router) *StopController {

    sc := StopController{ss}

    // Init Router
    r.HandleFunc("/{date}/nearest", sc.NearestStops)
    r.HandleFunc("/{date}/{stopId}", sc.StopById)

    return &sc
}


////////////////////////////////////////////////////////////////////////////////////////////////
/// Controller functions
////////////////////////////////////////////////////////////////////////////////////////////////

func (sc *StopController) NearestStops(w http.ResponseWriter, r *http.Request) {

    defer utils.RecoverFromError(w)

    sw := stopwatch.Start(0)

    params := mux.Vars(r)

    agencyKey := params["agencyKey"]
    date := params["date"]

    lat := r.URL.Query().Get("lat")
    lon := r.URL.Query().Get("lon")
    distance := r.URL.Query().Get("distance")

    if len(distance) <= 0 {
        distance = "1000"
    }

    timeOfDay := time.Now().Format("15:04:05")

//    log.Printf("Agency Key: %s", agencyKey)
//    log.Printf("Lat: %s", lat)
//    log.Printf("Lon: %s", lon)
//    log.Printf("Distance: %s", distance)
//    log.Printf("Date: %s", date)

    log.Printf("Fetching stops by date ...")
    stops := sc.ss.FetchStopsByDate(agencyKey, date, timeOfDay, lat, lon, distance, 3)
    log.Printf("Stops[%d] %v", len(stops), stops)

    log.Printf("Extracting Trip Ids ...")
    tripIds := sc.ss.ExtractTripIds(stops)
    log.Printf("TripIds[%d] %v", len(tripIds), tripIds)

    log.Printf("Fetching First And Last StopNames By Trip Ids ...")
    flStopNamesByTripId := sc.ss.FetchFirstAndLastStopNamesByTripIds(agencyKey, tripIds)

    log.Printf("Merge First and Last StopNames By TripId With Stop Routes ...")
    sc.ss.MergeFlStopNamesByTripIdWithStopRoutes(&stops, flStopNamesByTripId)

    stopGroups := stops.ToStopGroups()
    sort.Sort(models.StopGroupByDistance(stopGroups))

    jsm, err := json.Marshal(stopGroups.ToJsonStopGroups())

    utils.FailOnError(err, "Failed on marshalling json stops")

    w.Header().Set("Content-Type", "application/json");
    w.Header().Set("X-Response-Time", sw.ElapsedTime().String())
    w.Write(jsm)

    log.Printf("-----------------------------------------------------------------------------------")
    log.Printf("--- Nearest stops. ElapsedTime: %v", sw.ElapsedTime())
    log.Printf("-----------------------------------------------------------------------------------")
}


func (sc *StopController) StopById(w http.ResponseWriter, r *http.Request) {

    defer utils.RecoverFromError(w)

    sw := stopwatch.Start(0)

    params := mux.Vars(r)

    agencyKey := params["agencyKey"]
    date := params["date"]

    stopId := params["stopId"]

    limit := r.URL.Query().Get("limit")

    if len(limit) <= 0 {
        limit = "3"
    }

    timeOfDay := time.Now().Format("15:04:05")

    //    log.Printf("Agency Key: %s", agencyKey)
    //    log.Printf("Lat: %s", lat)
    //    log.Printf("Lon: %s", lon)
    //    log.Printf("Distance: %s", distance)
    //    log.Printf("Date: %s", date)

    //    log.Printf("Fetching stops by date ...")
    stops := sc.ss.FetchStopById(agencyKey, date, timeOfDay, stopId, 3)
    //    log.Printf("Stops[%d] %v", len(stops), stops)

    //    log.Printf("Extracting Trip Ids ...")
    tripIds := sc.ss.ExtractTripIds(stops)
    //    log.Printf("TripIds[%d] %v", len(tripIds), tripIds)

    //    log.Printf("Fetching First And Last StopNames By Trip Ids ...")
    flStopNamesByTripId := sc.ss.FetchFirstAndLastStopNamesByTripIds(agencyKey, tripIds)

    //    log.Printf("Merge First and Last StopNames By TripId With Stop Routes ...")
    sc.ss.MergeFlStopNamesByTripIdWithStopRoutes(&stops, flStopNamesByTripId)


    stopGroups := stops.ToStopGroups()
    sort.Sort(models.StopGroupByDistance(stopGroups))

    jsm, err := json.Marshal(stopGroups.ToJsonStopGroups())

    utils.FailOnError(err, "Failed on marshalling json stops")

    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("X-Response-Time", sw.ElapsedTime().String())
    w.Write(jsm)

    log.Printf("-----------------------------------------------------------------------------------")
    log.Printf("--- Stop by Id: '%s'. ElapsedTime: %v", stopId, sw.ElapsedTime())
    log.Printf("-----------------------------------------------------------------------------------")
}
