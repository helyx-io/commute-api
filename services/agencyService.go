package services

////////////////////////////////////////////////////////////////////////////////////////////////
/// Imports
////////////////////////////////////////////////////////////////////////////////////////////////

import (
	"log"
    "github.com/fatih/stopwatch"
    "github.com/helyx-io/commute-api/config"
    "github.com/helyx-io/commute-api/models"
    "github.com/helyx-io/commute-api/utils"
    "github.com/jinzhu/gorm"
    "gopkg.in/redis.v2"
    "github.com/helyx-io/commute-api/database"
)



////////////////////////////////////////////////////////////////////////////////////////////////
/// Agency Controller
////////////////////////////////////////////////////////////////////////////////////////////////

type AgencyService struct {
    db *gorm.DB
    redis *redis.Client
    connectInfos *config.DBConnectInfos
}


////////////////////////////////////////////////////////////////////////////////////////////////
/// Constructor
////////////////////////////////////////////////////////////////////////////////////////////////

func NewAgencyService(db *gorm.DB, connectInfos *config.DBConnectInfos, redis *redis.Client) *AgencyService {
    return &AgencyService{db, redis, connectInfos}
}


////////////////////////////////////////////////////////////////////////////////////////////////
/// Service functions
////////////////////////////////////////////////////////////////////////////////////////////////

func (as *AgencyService) FetchNearestAgencies(lat, lon string) models.Agencies {

    sw := stopwatch.Start(0)

    rows, err := database.Rows(as.db, as.connectInfos, "select-nearest-stations", lat, lat, lon, lon)
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

    agencies := make([]models.Agency, 0)

    for rows.Next() {
        rows.Scan(&key, &id, &name, &url, &timezone, &lang, &minLat, &maxLat, &minLon, &maxLon)

        agency := models.Agency{key, id, name, url, timezone, lang, minLat, maxLat, minLon, maxLon}

        agencies = append(agencies, agency)
    }

    return agencies
}
