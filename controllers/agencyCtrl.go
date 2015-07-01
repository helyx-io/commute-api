package controllers

////////////////////////////////////////////////////////////////////////////////////////////////
/// Imports
////////////////////////////////////////////////////////////////////////////////////////////////

import (
	"log"
    "net/http"
    "encoding/json"
    "github.com/gorilla/mux"
    "github.com/fatih/stopwatch"
    "github.com/helyx-io/commute-api/services"
    "github.com/helyx-io/commute-api/utils"
)



////////////////////////////////////////////////////////////////////////////////////////////////
/// Structures
////////////////////////////////////////////////////////////////////////////////////////////////

type AgencyController struct {
    as *services.AgencyService
    ss *services.StopService
}


////////////////////////////////////////////////////////////////////////////////////////////////
/// Constructor
////////////////////////////////////////////////////////////////////////////////////////////////

func NewAgencyController(as *services.AgencyService, ss *services.StopService, r *mux.Router) *AgencyController {

    ac := AgencyController{as, ss}

    // Init Router
    r.HandleFunc("/nearest", ac.NearestAgencies)

    NewStopController(ss, r.PathPrefix("/{agencyKey}/stops").Subrouter())

    return &ac
}


////////////////////////////////////////////////////////////////////////////////////////////////
/// Controller functions
////////////////////////////////////////////////////////////////////////////////////////////////

func (ac *AgencyController) NearestAgencies(w http.ResponseWriter, r *http.Request) {

    defer utils.RecoverFromError(w)

    sw := stopwatch.Start(0)

    lat := r.URL.Query().Get("lat")
    lon := r.URL.Query().Get("lon")

    log.Printf("Lat: %s", lat)
    log.Printf("Lon: %s", lon)

    log.Printf("-----------------------------------------------------------------------------------")
    log.Printf("--- Nearest agencies. ElapsedTime: %v", sw.ElapsedTime())
    log.Printf("-----------------------------------------------------------------------------------")

    w.Header().Set("X-Response-Time", sw.ElapsedTime().String())

    agencies := ac.as.FetchNearestAgencies(lat, lon)

    jsm, err := json.Marshal(agencies.ToJsonAgencies())

    utils.FailOnError(err, "Failed on marshalling json agencies")
    w.Write(jsm)
}
