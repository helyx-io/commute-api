package controller

////////////////////////////////////////////////////////////////////////////////////////////////
/// Imports
////////////////////////////////////////////////////////////////////////////////////////////////

import (
    "fmt"
	"log"
    "sync"
    "time"
    "sort"
    "net/http"
    "encoding/json"
    "gopkg.in/redis.v2"
    "github.com/gorilla/mux"
    "github.com/fatih/stopwatch"
    "github.com/helyx-io/gtfs-api/config"
    "github.com/helyx-io/gtfs-api/utils"
)


////////////////////////////////////////////////////////////////////////////////////////////////
/// Structures
////////////////////////////////////////////////////////////////////////////////////////////////

type Stops []Stop

type Stop struct {
    Id int
    Name string
    Desc string
    Lat string
    Lon string
    LocationType int
    Distance float64
    Routes Routes
}

type StopGroups []StopGroup

type StopGroup struct {
    Ids []int
    Name string
    Desc string
    Lat string
    Lon string
    LocationType int
    Distance float64
    Routes Routes
}

type StopGroupByDistance []StopGroup
func (sgbd StopGroupByDistance) Len() int { return len(sgbd) }
func (sgbd StopGroupByDistance) Swap(i, j int) { sgbd[i], sgbd[j] = sgbd[j], sgbd[i] }
func (sgbd StopGroupByDistance) Less(i, j int) bool { return sgbd[i].Distance < sgbd[j].Distance }


type Routes []Route

type Route struct {
    Name string
    TripId int
    RouteType int
    RouteColor string
    RouteTextColor string
    FirstStopName string
    LastStopName string
    StopTimesFull []StopTimeFull
}

type StopTimeFull struct {
    StopId int
    StopName string
    StopDesc string
    StopLat string
    StopLon string
    LocationType int
    ArrivalTime string
    DepartureTime string
    StopSequence int
    DirectionId int
    RouteShortName string
    RouteType int
    RouteColor string
    RouteTextColor string
    TripId int
}

type StopTimeFullByDepartureDate []StopTimeFull
func (stf StopTimeFullByDepartureDate) Len() int { return len(stf) }
func (stf StopTimeFullByDepartureDate) Swap(i, j int) { stf[i], stf[j] = stf[j], stf[i] }
func (stf StopTimeFullByDepartureDate) Less(i, j int) bool { return stf[i].DepartureTime < stf[j].DepartureTime }

type RouteByName []Route
func (rn RouteByName) Len() int { return len(rn) }
func (rn RouteByName) Swap(i, j int) { rn[i], rn[j] = rn[j], rn[i] }
func (rn RouteByName) Less(i, j int) bool { return rn[i].Name < rn[j].Name }


type FirstLastStopNamesByTripId struct {
    TripId int
    FirstStopName string
    LastStopName string
}

type JsonStopGroup struct {
    Ids []int `json:"ids"`
    Name string `json:"name"`
    Desc string `json:"desc"`
    Distance float64 `json:"distance"`
    LocationType int `json:"location_type"`
    GeoLocation JsonGeoLocation `json:"geo_location"`
    Routes []JsonRoute `json:"routes"`
}

type JsonGeoLocation struct {
    Lat string `json:"lat"`
    Lon string `json:"lon"`
}

type JsonRoute struct {
    Name string `json:"name"`
    TripId int `json:"trip_id"`
    RouteType int `json:"route_type"`
    RouteColor string `json:"route_color"`
    RouteTextColor string `json:"route_text_color"`
    FirstStopName string `json:"first_stop_name"`
    LastStopName string `json:"last_stop_name"`
    StopTimes []string `json:"stop_times"`
}


////////////////////////////////////////////////////////////////////////////////////////////////
/// Variables
////////////////////////////////////////////////////////////////////////////////////////////////

var (
    redisClient *redis.Client
)


////////////////////////////////////////////////////////////////////////////////////////////////
/// Agency Controller
////////////////////////////////////////////////////////////////////////////////////////////////

type StopController struct { }

func (sc *StopController) Init(r *mux.Router) {
    redisClient = redis.NewTCPClient(&redis.Options{
        Addr: fmt.Sprintf("%s:%d", config.RedisInfos.Host, config.RedisInfos.Port),
        Password: "", // no password set
        DB:       0,  // use default DB
        PoolSize: 16,
    })

    // Init Router
    r.HandleFunc("/{date}/nearest", sc.NearestStops)
}


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

    log.Printf("Agency Key: %s", agencyKey)
    log.Printf("Lat: %s", lat)
    log.Printf("Lon: %s", lon)
    log.Printf("Distance: %s", distance)
    log.Printf("Date: %s", date)


    log.Printf("Fetching stops by date ...")
    stops := fetchStopsByDate(agencyKey, date, lat, lon, distance)

    log.Printf("Extracting Trip Ids ...")
    tripIds := extractTripIds(stops)

    log.Printf("Fetching First And Last StopNames By Trip Ids ...")
    flStopNamesByTripId := fetchFirstAndLastStopNamesByTripIds(agencyKey, tripIds)

    log.Printf("Merge First and Last StopNames By TripId With Stop Routes ...")
    mergeFlStopNamesByTripIdWithStopRoutes(&stops, flStopNamesByTripId)

    log.Printf("-----------------------------------------------------------------------------------")
    log.Printf("--- Nearest stops. ElapsedTime: %v", sw.ElapsedTime())
    log.Printf("-----------------------------------------------------------------------------------")


    w.Header().Set("X-Response-Time", sw.ElapsedTime().String())

    stopGroups := stops.toStopGroups()
    sort.Sort(StopGroupByDistance(stopGroups))

    jsm, err := json.Marshal(stopGroups.toJsonStopGroups())

    utils.FailOnError(err, "Failed on marshalling json stops")
    w.Write(jsm)
}


func (ss Stops) toStopGroups() StopGroups {
    sgsByKey := make(map[string]StopGroup)

    for _, s := range ss {
        key := fmt.Sprintf("%s%s%d", s.Name, s.Desc, s.LocationType)

        if sg, ok := sgsByKey[key]; !ok {
            sgsByKey[key] = s.toStopGroup()
        } else {
            sg.Ids = append(sg.Ids, s.Id)
            for _, r := range s.Routes {
                sg.Routes = append(sg.Routes, r)
            }

            sort.Sort(RouteByName(sg.Routes))

            sgsByKey[key] = sg
        }
    }
    sgs := make(StopGroups, 0)

    for _, value := range sgsByKey {
//        log.Printf("Key: %s - Value: %v", key, value)
        sgs = append(sgs, value)
    }

    return sgs
}

func (s Stop) toStopGroup() StopGroup {
    return StopGroup{[]int{s.Id}, s.Name, s.Desc, s.Lat, s.Lon, s.LocationType, s.Distance, s.Routes}
}


func (sgs StopGroups) toJsonStopGroups() []JsonStopGroup {
    jsgs := make([]JsonStopGroup, len(sgs))

    for i, sg := range sgs {
        jsgs[i] = sg.toJsonStopGroup()
//        log.Printf("Index: %d - json: %v", i, jsgs[i])
    }

    return jsgs
}

func (rs Routes) toJsonRoutes() []JsonRoute {
    jrs := make([]JsonRoute, len(rs))

    for i, r := range rs {
        jrs[i] = r.toJsonRoute()
    }

    return jrs
}

func (sg StopGroup) toJsonStopGroup() JsonStopGroup {
    return JsonStopGroup{sg.Ids, sg.Name, sg.Desc, sg.Distance, sg.LocationType, JsonGeoLocation{sg.Lat, sg.Lon}, sg.Routes.toJsonRoutes()}
}

func (r *Route) toJsonRoute() JsonRoute {

    jstfs := make([]string, len(r.StopTimesFull))

    for i, stf := range r.StopTimesFull {
        jstfs[i] = stf.DepartureTime
    }

    return JsonRoute{r.Name, r.TripId, r.RouteType, r.RouteColor, r.RouteTextColor, r.FirstStopName, r.LastStopName, jstfs}
}

func mergeFlStopNamesByTripIdWithStopRoutes(stops *Stops, flStopNamesByTripId map[int]FirstLastStopNamesByTripId) {
    for i := range *stops {
        stop := &(*stops)[i]
        for j := range stop.Routes {
            route := &stop.Routes[j]
            route.FirstStopName = flStopNamesByTripId[route.TripId].FirstStopName
            route.LastStopName = flStopNamesByTripId[route.TripId].LastStopName
        }
    }

}

func fetchFirstAndLastStopNamesByTripIds(agencyKey string, tripIds []int) map[int]FirstLastStopNamesByTripId {

    flStopNamesByTripIdChan := make(chan FirstLastStopNamesByTripId)

    sem := make(chan bool, 64)

    go func() {
        for _, tripId := range tripIds {

            sem <- true

            go func(tripId int) {

                defer func() { <-sem }()

//                sw := stopwatch.Start(0)

                key := fmt.Sprintf("/%s/t/st/fl/%d", agencyKey, tripId)
                tripPayload := redisClient.Get(key)
                value := tripPayload.Val()

                tripFirstLast := make([]string, 2)

                err := json.Unmarshal([]byte(value), &tripFirstLast)
                if err != nil {
                    log.Printf(" * Error: '%s' ...", err.Error())
                }

//                log.Printf("[TRIP][FIND_STOP_TIMES_BY_TRIP_ID] Data Fetch for key: '%s' Done in %v", key, sw.ElapsedTime());

                flStopNamesByTripIdChan <- FirstLastStopNamesByTripId{tripId, tripFirstLast[0], tripFirstLast[1]}
            }(tripId)
        }

        for i := 0; i < cap(sem); i++ {
            sem <- true
        }

        close(flStopNamesByTripIdChan)
    }()

    flStopNamesByTripIds := make(map[int]FirstLastStopNamesByTripId)

    for flStopNamesByTripId := range flStopNamesByTripIdChan {
        flStopNamesByTripIds[flStopNamesByTripId.TripId] = flStopNamesByTripId
    }

    return flStopNamesByTripIds
}

func extractTripIds(stops Stops) []int {
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

func fetchStopsByDate(agencyKey, date, lat, lon, distance string) Stops {

    query := fmt.Sprintf("select s.stop_id, s.stop_name, s.stop_desc, s.stop_lat, s.stop_lon, s.location_type, 111195 * st_distance(point(%s, %s), s.stop_geo) as stop_distance from gtfs_%s.stops s where 111195 * st_distance(point(%s, %s), s.stop_geo) < %s order by stop_distance asc", lat, lon, agencyKey, lat, lon, distance)
    sw := stopwatch.Start(0)

    log.Printf("Query: %s", query)
    rows, err := config.DB.Raw(query).Rows()
    defer rows.Close()

    log.Printf("[STOP_SERVICE][FIND_NEAREST_STOPS] Data Fetch for [agencyKey=%s, date=%s, lat=%s, lon=%s, distance=%s] Done in %v", agencyKey, date, lat, lon, distance, sw.ElapsedTime());

    stopChan := make(chan Stop)

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

            stop := Stop{id, name, desc, lat, lon, locationType, distance, nil}

//            log.Printf("Stop: %v", stop)

            sem <- true

            go func(stop Stop) {
                defer func() { <-sem }()

                stop.Routes = fetchRoutesForDateAndStop(agencyKey, date, stop)

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

    stops := make(Stops, 0)

    for stop := range stopChan {
        stops = append(stops, stop)
    }

    return stops
}

func fetchRoutesForDateAndStop(agencyKey, date string, stop Stop) Routes {
//    log.Printf("Fetching routes for stop: %v", stop)

    stfs := fetchStopTimesFullForDateAndStop(agencyKey, date, stop)

    return groupStopTimesFullByRoute(stfs)
}

func groupStopTimesFullByRoute(stfs []StopTimeFull) Routes {

    stfsByRouteShortName := make(map[string][]StopTimeFull, 0)

    for _, stf := range stfs {
        if _, ok := stfsByRouteShortName[stf.RouteShortName]; !ok {
            stfsByRouteShortName[stf.RouteShortName] = make([]StopTimeFull, 0)
        }

        stfsByRouteShortName[stf.RouteShortName] = append(stfsByRouteShortName[stf.RouteShortName], stf)
    }

    routes := make(Routes, 0)

    for rsn, stfs := range stfsByRouteShortName {
        if len(stfs) > 0 {

            sort.Sort(StopTimeFullByDepartureDate(stfs))

            routes = append(routes, Route{rsn, stfs[0].TripId, stfs[0].RouteType, stfs[0].RouteColor, stfs[0].RouteTextColor, "", "", stfs})
        }
    }

    return routes
}

func fetchStopTimesFullForDateAndStop(agencyKey, date string, stop Stop) []StopTimeFull {
//    log.Printf("Fetching stop times full for date: %s & stop: %v", date, stop)

    day, _ := time.Parse("2006-01-02", date)
    dayOfWeek := day.Weekday().String()

    stfChan := make(chan StopTimeFull, 2)

    go func() {
        var wg sync.WaitGroup
        wg.Add(2)

//        sw := stopwatch.Start(0)

        go fetchStopTimesFullForCalendar(agencyKey, stop, date, dayOfWeek, stfChan, &wg)
        go fetchStopTimesFullForCalendarDates(agencyKey, stop, date, stfChan, &wg)

        wg.Wait()

//        log.Printf("[STOP_TIMES_FULL][FIND_LINES_BY_STOP_ID_AND_DATE] Data Fetch done in %v", sw.ElapsedTime());

        close(stfChan)
    }()

    stfs := make([]StopTimeFull, 0)

    for stf := range stfChan {
        if (len(stfs) < 5) {
            stfs = append(stfs, stf)
        }
    }

    return stfs
}

func fetchStopTimesFullForCalendar(agencyKey string, stop Stop, date, dayOfWeek string, stfChan chan StopTimeFull, wg *sync.WaitGroup) {
    queryCalendar := fmt.Sprintf(
        "select" +
        "    stf.stop_id," +
        "    stf.stop_name," +
        "    stf.stop_desc," +
        "    stf.stop_lat," +
        "    stf.stop_lon," +
        "    stf.location_type," +
        "    stf.arrival_time," +
        "    stf.departure_time," +
        "    stf.stop_sequence," +
        "    stf.direction_id," +
        "    stf.route_short_name," +
        "    stf.route_type," +
        "    stf.route_color," +
        "    stf.route_text_color," +
        "    stf.trip_id" +
        " from " +
        "    gtfs_%s.stop_times_full stf inner join" +
        "    gtfs_%s.calendars c on stf.service_id=c.service_id" +
        " where " +
        "    stf.stop_id=%d and" +
        "    c.start_date <= '%s' and" +
        "    c.end_date >= '%s' and %s=1", agencyKey, agencyKey, stop.Id, date, date, dayOfWeek)
//        "    c.end_date >= '%s' and %s=1", agencyKey, agencyKey, stop.Id, date, date, dayOfWeek/*, time.Now().Format("15:04:05"), 5*/)

//    log.Printf("Query calendar: %s", queryCalendar)

    calendarRows, err := config.DB.Raw(queryCalendar).Rows()

    utils.FailOnError(err, "Calendar row fetch error")

    defer func() {
        calendarRows.Close()
        wg.Done()
    }()

    var stopId, locationType, stopSequence, directionId, routeType, tripId int
    var stopName, stopDesc, stopLat, stopLon, arrivalTime, departureTime, routeShortName, routeColor, routeTextColor string

    for calendarRows.Next() {
        calendarRows.Scan(
        &stopId, &stopName, &stopDesc, &stopLat, &stopLon, &locationType, &arrivalTime, &departureTime,
        &stopSequence, &directionId, &routeShortName, &routeType, &routeColor, &routeTextColor, &tripId,
        )

        log.Printf("StopId: %s", stopId)
        log.Printf("StopName: %s", stopName)

        stfChan <- StopTimeFull{stopId, stopName, stopDesc, stopLat, stopLon, locationType, arrivalTime, departureTime, stopSequence, directionId, routeShortName, routeType, routeColor, routeTextColor, tripId}
    }
}

func fetchStopTimesFullForCalendarDates(agencyKey string, stop Stop, date string, stfChan chan StopTimeFull, wg *sync.WaitGroup) {
    queryCalendarDate := fmt.Sprintf(
        "select" +
        "    stf.stop_id," +
        "    stf.stop_name," +
        "    stf.stop_desc," +
        "    stf.stop_lat," +
        "    stf.stop_lon," +
        "    stf.location_type," +
        "    stf.arrival_time," +
        "    stf.departure_time," +
        "    stf.stop_sequence," +
        "    stf.direction_id," +
        "    stf.route_short_name," +
        "    stf.route_type," +
        "    stf.route_color," +
        "    stf.route_text_color," +
        "    stf.trip_id" +
        " from" +
        "    gtfs_%s.stop_times_full stf inner join" +
        "    gtfs_%s.calendar_dates cd on stf.service_id=cd.service_id" +
        " where" +
        "    stf.stop_id=%d and cd.date = '%s'", agencyKey, agencyKey, stop.Id, date)
//        "    stf.stop_id=%d and cd.date = '%s'", agencyKey, agencyKey, stop.Id, date/*, time.Now().Format("15:04:05"), 5*/)

//    log.Printf("Query calendar dates : %s", queryCalendarDate)

    calendarDateRows, _ := config.DB.Raw(queryCalendarDate).Rows()

    defer func() {
        calendarDateRows.Close()
        wg.Done()
    }()

    var stopId, locationType, stopSequence, directionId, routeType, tripId int
    var stopName, stopDesc, stopLat, stopLon, arrivalTime, departureTime, routeShortName, routeColor, routeTextColor string

    for calendarDateRows.Next() {
        calendarDateRows.Scan(
        &stopId, &stopName, &stopDesc, &stopLat, &stopLon, &locationType, &arrivalTime, &departureTime,
        &stopSequence, &directionId, &routeShortName, &routeType, &routeColor, &routeTextColor, &tripId,
        )

        stfChan <- StopTimeFull{stopId, stopName, stopDesc, stopLat, stopLon, locationType, arrivalTime, departureTime, stopSequence, directionId, routeShortName, routeType, routeColor, routeTextColor, tripId}
    }

}
