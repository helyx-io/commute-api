package controller

////////////////////////////////////////////////////////////////////////////////////////////////
/// Imports
////////////////////////////////////////////////////////////////////////////////////////////////

import (
	"log"
    "net/http"
    "encoding/json"
    "github.com/gorilla/mux"
    "github.com/fatih/stopwatch"
    "github.com/helyx-io/gtfs-api/config"
    "github.com/helyx-io/gtfs-api/utils"
    "github.com/jinzhu/gorm"
    "gopkg.in/redis.v2"
    "github.com/helyx-io/gtfs-api/database"
)


////////////////////////////////////////////////////////////////////////////////////////////////
/// Structures
////////////////////////////////////////////////////////////////////////////////////////////////

type Agencies []Agency

type Agency struct {
    Key string
    Id int
    Name string
    Url string
    Timezone string
    Lang string
    MinLat float64
    MaxLat float64
    MinLon float64
    MaxLon float64
}


type JsonAgency struct {
    Key string `json:"key"`
    Id int `json:"id"`
    Name string `json:"name"`
    Url string `json:"url"`
    Timezone string `json:"timezone"`
    Lang string `json:"lang"`
    MinLat float64 `json:"min_lat"`
    MaxLat float64 `json:"max_lat"`
    MinLon float64 `json:"min_lon"`
    MaxLon float64 `json:"max_lon"`
}


////////////////////////////////////////////////////////////////////////////////////////////////
/// Agency Controller
////////////////////////////////////////////////////////////////////////////////////////////////

type AgencyController struct {
    db *gorm.DB
    redis *redis.Client
    connectInfos *config.DBConnectInfos
}

func (ac *AgencyController) Init(db *gorm.DB, connectInfos *config.DBConnectInfos, redis *redis.Client, r *mux.Router) {
    // Init Router
    r.HandleFunc("/nearest", ac.NearestAgencies)

    ac.db = db
    ac.connectInfos = connectInfos
    ac.redis = redis

    new(StopController).Init(db, connectInfos, redis, r.PathPrefix("/{agencyKey}/stops").Subrouter())
}

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

    agencies := ac.fetchNearestAgencies(lat, lon)

    jsm, err := json.Marshal(agencies.toJsonAgencies())

    utils.FailOnError(err, "Failed on marshalling json agencies")
    w.Write(jsm)
}

func (as Agencies) toJsonAgencies() []JsonAgency {
    jas := make([]JsonAgency, len(as))

    for i, a := range as {
        jas[i] = a.toJsonAgency()
    }

    return jas
}

func (a Agency) toJsonAgency() JsonAgency {
    return JsonAgency{a.Key, a.Id, a.Name, a.Url, a.Timezone, a.Lang, a.MinLat, a.MaxLat, a.MinLon, a.MaxLon}
}


func (ac *AgencyController) fetchNearestAgencies(lat, lon string) Agencies {

    sw := stopwatch.Start(0)

    rows, err := database.Rows(ac.db, ac.connectInfos, "select-nearest-stations", lat, lat, lon, lon)
    defer rows.Close()

    log.Printf("[STOP_SERVICE][FIND_NEAREST_AGENCIES] Data Fetch for [lat=%s, lon=%s] Done in %v", lat, lon, sw.ElapsedTime());

    utils.FailOnError(err, "Failed to execute query")

    key := ""
    id := 0
    name := ""
    url := ""
    timezone := ""
    lang := ""
    minLat := 0.0
    maxLat := 0.0
    minLon := 0.0
    maxLon := 0.0

    agencies := make([]Agency, 0)

    for rows.Next() {
        rows.Scan(&key, &id, &name, &url, &timezone, &lang, &minLat, &maxLat, &minLon, &maxLon)

        agency := Agency{key, id, name, url, timezone, lang, minLat, maxLat, minLon, maxLon}

        agencies = append(agencies, agency)
    }

    return agencies
}
