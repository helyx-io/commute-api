package controller

////////////////////////////////////////////////////////////////////////////////////////////////
/// Imports
////////////////////////////////////////////////////////////////////////////////////////////////

import (
    "fmt"
	"log"
    "net/http"
    "encoding/json"
    "github.com/gorilla/mux"
    "github.com/fatih/stopwatch"
    "github.com/helyx-io/gtfs-api/config"
    "github.com/helyx-io/gtfs-api/utils"
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

type AgencyController struct { }

func (ac *AgencyController) Init(r *mux.Router) {
    // Init Router
    r.HandleFunc("/nearest", ac.NearestAgencies)

    new(StopController).Init(r.PathPrefix("/{agencyKey}/stops").Subrouter())
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

    agencies := fetchNearestAgencies(lat, lon)

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


func fetchNearestAgencies(lat, lon string) Agencies {

    query := fmt.Sprintf("select agency_key, agency_id, agency_name, agency_url, agency_timezone, agency_lang, agency_min_lat, agency_max_lat, agency_min_lon, agency_max_lon from agencies where agency_min_lat <= %s and agency_max_lat >= %s and agency_min_lon <= %s and agency_max_lon >= %s", lat, lat, lon, lon)
    sw := stopwatch.Start(0)

    log.Printf("Query: %s", query)
    rows, err := config.DB.Raw(query).Rows()
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
