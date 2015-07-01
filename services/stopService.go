package services

////////////////////////////////////////////////////////////////////////////////////////////////
/// Imports
////////////////////////////////////////////////////////////////////////////////////////////////

import (
    "fmt"
    "log"
    "sync"
    "time"
    "sort"
    "encoding/json"
    "gopkg.in/redis.v2"
    "github.com/jinzhu/gorm"
    "github.com/fatih/stopwatch"
    "github.com/helyx-io/commute-api/utils"
    "github.com/helyx-io/commute-api/models"
    "github.com/helyx-io/commute-api/config"
    "github.com/helyx-io/commute-api/database"
    "github.com/helyx-io/commute-api/data"
)


////////////////////////////////////////////////////////////////////////////////////////////////
/// Structures
////////////////////////////////////////////////////////////////////////////////////////////////

type StopService struct {
    db *gorm.DB
    redis *redis.Client
    connectInfos *config.DBConnectInfos
    selectStopsByDate string
    selectStopById string
}


////////////////////////////////////////////////////////////////////////////////////////////////
/// Constructor
////////////////////////////////////////////////////////////////////////////////////////////////

func NewStopService(db *gorm.DB, connectInfos *config.DBConnectInfos, redis *redis.Client) *StopService {

    selectStopsByDate := loadAsset(connectInfos.Dialect, "select-stops-by-date")
    selectStopById := loadAsset(connectInfos.Dialect, "select-stop-by-id")
    //    log.Printf("Exec Stmt: '%s'", sc.selectStopsByDate)

    return &StopService{db, redis, connectInfos, selectStopsByDate, selectStopById}
}


////////////////////////////////////////////////////////////////////////////////////////////////
/// Helper functions
////////////////////////////////////////////////////////////////////////////////////////////////

func loadAsset(dialect, dmlKey string) string {
    filePath := fmt.Sprintf("resources/ddl/%s/%s.sql", dialect, dmlKey)
    //    log.Printf("Executing query from file path: '%s' - Params: %v", filePath, params)

    dml, err := data.Asset(filePath)
    utils.FailOnError(err, fmt.Sprintf("Could get dml resource at path '%s' for exec", filePath))

    stmt := string(dml)
    //    log.Printf("Exec Stmt: '%s'", stmt)

    return stmt
}


////////////////////////////////////////////////////////////////////////////////////////////////
/// Service functions
////////////////////////////////////////////////////////////////////////////////////////////////


func (ss *StopService) MergeFlStopNamesByTripIdWithStopRoutes(stops *models.Stops, flStopNamesByTripId map[int]models.FirstLastStopNamesByTripId) {
    for i := range *stops {
        stop := &(*stops)[i]
        for j := range stop.Routes {
            route := &stop.Routes[j]
            route.FirstStopName = flStopNamesByTripId[route.TripId].FirstStopName
            route.LastStopName = flStopNamesByTripId[route.TripId].LastStopName
        }
    }

}


func (ss *StopService) FetchFirstAndLastStopNamesByTripIds(agencyKey string, tripIds []int) map[int]models.FirstLastStopNamesByTripId {

    //    sw := stopwatch.Start(0)

    keys := make([]string, len(tripIds))

    for i, tripId := range tripIds {
        keys[i] = fmt.Sprintf("/%s/t/st/fl/%d", agencyKey, tripId)
    }

    tripPayloads, _ := ss.redis.MGet(keys...).Result()
    flStopNamesByTripIds := make(map[int]models.FirstLastStopNamesByTripId)

    for i, tripPayload := range tripPayloads {
        tripId := tripIds[i]
        //        log.Printf("tripPayload %d: %v", tripId, tripPayload)
        value := tripPayload.(string)

        tripFirstLast := make([]string, 2)

        err := json.Unmarshal([]byte(value), &tripFirstLast)
        if err != nil {
            log.Printf(" * Error: '%s' ...", err.Error())
        }

        //        log.Printf("[TRIP][FIND_STOP_TIMES_BY_TRIP_ID] Data Fetch for tripIds: '%v' Done in %v", tripIds, sw.ElapsedTime());
        //        log.Printf("[TRIP][FIND_STOP_TIMES_BY_TRIP_ID] Data Fetch for %d tripIds Done in %v", len(tripIds), sw.ElapsedTime());

        flStopNamesByTripIds[tripId] = models.FirstLastStopNamesByTripId{tripId, tripFirstLast[0], tripFirstLast[1]}
    }

    return flStopNamesByTripIds
}


func (ss *StopService) ExtractTripIds(stops models.Stops) []int {
    tripIdMap := make(map[int]bool)

    for _, stop := range stops {
        for _, route := range stop.Routes {
            if len(route.StopTimesFull) > 0 {
                tripIdMap[route.StopTimesFull[0].TripId] = true
            }
        }
    }

    tripIds := make([]int, 0, len(tripIdMap))
    for tripId := range tripIdMap {
        tripIds = append(tripIds, tripId)
    }

    return tripIds
}


func (ss *StopService) FetchStopsByDate(agencyKey, date, timeOfDay, lat, lon, distance string, limit int) models.Stops {

    sw := stopwatch.Start(0)

    schema := fmt.Sprintf("gtfs_%s", agencyKey)
    query := fmt.Sprintf(ss.selectStopsByDate, lat, lon, schema, lat, lon, distance)
    log.Printf("Exec Stmt 1: '%s'", query)

    rows, err := ss.db.Raw(query).Rows()
    defer rows.Close()

    log.Printf("[STOP_SERVICE][FECH_STOPS_BY_DATE] --- Data Fetch for [agencyKey=%s, date=%s, lat=%s, lon=%s, distance=%s] Done in %v", agencyKey, date, lat, lon, distance, sw.ElapsedTime());

    stopChan := make(chan models.Stop)

    go func() {

        utils.FailOnError(err, "Failed to execute query")

        id := 0
        name := ""
        desc := ""
        lat := ""
        lon := ""
        locationType := 0
        distance := 0.0

        sem := make(chan bool, 512)

        for rows.Next() {

            rows.Scan(&id, &name, &desc, &lat, &lon, &locationType, &distance)

            stop := models.Stop{id, name, desc, lat, lon, locationType, distance, nil}

            log.Printf("--- Stop: %v", stop)

            sem <- true

            go func(stop models.Stop) {
                defer func() { <-sem }()

                stop.Routes = ss.fetchRoutesForDateAndStop(agencyKey, date, timeOfDay, stop, limit)

                if len(stop.Routes) > 0 {
                    stopChan <- stop
                }

            }(stop)
        }

        for i := 0; i < cap(sem); i++ {
            sem <- true
        }

        close(stopChan)
    }()

    stops := make(models.Stops, 0)

    for stop := range stopChan {
        stops = append(stops, stop)
    }

    return stops
}


func (ss *StopService) FetchStopById(agencyKey, date, timeOfDay, stopId string, limit int) models.Stops {

    sw := stopwatch.Start(0)

    schema := fmt.Sprintf("gtfs_%s", agencyKey)
    query := fmt.Sprintf(ss.selectStopById, schema, stopId)
    log.Printf("Exec Stmt: '%s'", query)

    rows, err := ss.db.Raw(query).Rows()
    defer rows.Close()

    log.Printf("[STOP_SERVICE][FETCH_STOP_BY_ID] Data Fetch for [agencyKey=%s, date=%s, stopId=%s] Done in %v", agencyKey, date, stopId, sw.ElapsedTime());

    stopChan := make(chan models.Stop)

    go func() {

        utils.FailOnError(err, "Failed to execute query")

        id := 0
        name := ""
        desc := ""
        lat := ""
        lon := ""
        locationType := 0

        sem := make(chan bool, 512)

        for rows.Next() {
            rows.Scan(&id, &name, &desc, &lat, &lon, &locationType)

            stop := models.Stop{id, name, desc, lat, lon, locationType, 0, nil} // Set distance to 0

            log.Printf("Stop: %v", stop)

            sem <- true

            go func(stop models.Stop) {
                defer func() { <-sem }()

                stop.Routes = ss.fetchRoutesForDateAndStop(agencyKey, date, timeOfDay, stop, limit)

                if len(stop.Routes) > 0 {
                    stopChan <- stop
                }

            }(stop)
        }

        for i := 0; i < cap(sem); i++ {
            sem <- true
        }

        close(stopChan)
    }()

    stops := make(models.Stops, 0)

    for stop := range stopChan {
        stops = append(stops, stop)
    }

    return stops
}


func (ss *StopService) fetchRoutesForDateAndStop(agencyKey, date, timeOfDay string, stop models.Stop, limit int) models.Routes {
    //    log.Printf("Fetching routes for stop: %v", stop)

    stfs := ss.fetchStopTimesFullForDateAndStop(agencyKey, date, timeOfDay, stop, limit)

    return ss.groupStopTimesFullByRoute(stfs)
}


func (ss *StopService) groupStopTimesFullByRoute(stfs []models.StopTimeFull) models.Routes {

    stfsByRouteShortName := make(map[string][]models.StopTimeFull, 0)

    for _, stf := range stfs {
        if _, ok := stfsByRouteShortName[stf.RouteShortName]; !ok {
            stfsByRouteShortName[stf.RouteShortName] = make([]models.StopTimeFull, 0)
        }

        stfsByRouteShortName[stf.RouteShortName] = append(stfsByRouteShortName[stf.RouteShortName], stf)
    }

    routes := make(models.Routes, 0)

    for rsn, stfs := range stfsByRouteShortName {
        if len(stfs) > 0 {

            sort.Sort(models.StopTimeFullByDepartureDate(stfs))

            routes = append(routes, models.Route{rsn, stfs[0].TripId, stfs[0].RouteType, stfs[0].RouteColor, stfs[0].RouteTextColor, "", "", stfs})
        }
    }

    return routes
}


func (ss *StopService) fetchStopTimesFullForDateAndStop(agencyKey, date, timeOfDay string, stop models.Stop, limit int) []models.StopTimeFull {
    //    log.Printf("Fetching stop times full for date: %s & stop: %v", date, stop)

    day, _ := time.Parse("2006-01-02", date)
    dayOfWeek := day.Weekday().String()

    stfChan := make(chan models.StopTimeFull, 2)

    go func() {
        var wg sync.WaitGroup
        wg.Add(2)

        //        sw := stopwatch.Start(0)

        go ss.fetchStopTimesFullForCalendar(agencyKey, stop, date, dayOfWeek, stfChan, &wg)
        go ss.fetchStopTimesFullForCalendarDates(agencyKey, stop, date, stfChan, &wg)

        wg.Wait()

        //        log.Printf("[STOP_TIMES_FULL][FIND_LINES_BY_STOP_ID_AND_DATE] Data Fetch done in %v", sw.ElapsedTime());

        close(stfChan)
    }()

    stfs := make([]models.StopTimeFull, 0)

    // currentTime := time.Now().Format("15:04:05")

    for stf := range stfChan {
        //        log.Printf("stfs: %v - currentTime: %v", stf.DepartureTime, currentTime)
        if stf.DepartureTime >= timeOfDay && (limit <= 0 || len(stfs) < limit) {
            stfs = append(stfs, stf)
        }
    }

    return stfs
}


func (ss *StopService) fetchStopTimesFullForCalendar(agencyKey string, stop models.Stop, date, dayOfWeek string, stfChan chan models.StopTimeFull, wg *sync.WaitGroup) {

    calendarRows, err := database.Rows(ss.db, ss.connectInfos, "select-stop-times-by-calendars", agencyKey, agencyKey, stop.Id, date, date, dayOfWeek)

    utils.FailOnError(err, "Calendars row fetch error")

    defer func() {
        calendarRows.Close()
        wg.Done()
    }()

    var stopId, locationType, stopSequence, directionId, routeType, tripId int
    var stopName, stopDesc, stopLat, stopLon, arrivalTime, departureTime, routeShortName string

    var routeColor, routeTextColor int32

    for calendarRows.Next() {
        calendarRows.Scan(
            &stopId, &stopName, &stopDesc, &stopLat, &stopLon, &locationType, &arrivalTime, &departureTime,
            &stopSequence, &directionId, &routeShortName, &routeType, &routeColor, &routeTextColor, &tripId,
        )

        stfChan <- models.StopTimeFull{stopId, stopName, stopDesc, stopLat, stopLon, locationType, arrivalTime, departureTime, stopSequence, directionId, routeShortName, routeType, utils.Int32ToColor(routeColor), utils.Int32ToColor(routeTextColor), tripId}
    }
}


func (ss *StopService) fetchStopTimesFullForCalendarDates(agencyKey string, stop models.Stop, date string, stfChan chan models.StopTimeFull, wg *sync.WaitGroup) {

    calendarDateRows, err := database.Rows(ss.db, ss.connectInfos, "select-stop-times-by-calendar-dates", agencyKey, agencyKey, stop.Id, date)

    utils.FailOnError(err, "Calendar dates row fetch error")

    defer func() {
        calendarDateRows.Close()
        wg.Done()
    }()

    var stopId, locationType, stopSequence, directionId, routeType, tripId int
    var stopName, stopDesc, stopLat, stopLon, arrivalTime, departureTime, routeShortName string

    var routeColor, routeTextColor int32

    for calendarDateRows.Next() {
        calendarDateRows.Scan(
            &stopId, &stopName, &stopDesc, &stopLat, &stopLon, &locationType, &arrivalTime, &departureTime,
            &stopSequence, &directionId, &routeShortName, &routeType, &routeColor, &routeTextColor, &tripId,
        )

        //        log.Printf("departureTime: %s", departureTime)

        stfChan <- models.StopTimeFull{stopId, stopName, stopDesc, stopLat, stopLon, locationType, arrivalTime, departureTime, stopSequence, directionId, routeShortName, routeType, utils.Int32ToColor(routeColor), utils.Int32ToColor(routeTextColor), tripId}
    }

}
