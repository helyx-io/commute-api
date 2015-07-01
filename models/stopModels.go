package models


////////////////////////////////////////////////////////////////////////////////////////////////
/// Imports
////////////////////////////////////////////////////////////////////////////////////////////////

import (
    "fmt"
    "sort"
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
/// Convert functions
////////////////////////////////////////////////////////////////////////////////////////////////

func (ss Stops) ToStopGroups() StopGroups {
    sgsByKey := make(map[string]StopGroup)

    for _, s := range ss {
        key := fmt.Sprintf("%s%s%d", s.Name, s.Desc, s.LocationType)

        if sg, ok := sgsByKey[key]; !ok {
            sgsByKey[key] = s.ToStopGroup()
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

func (s Stop) ToStopGroup() StopGroup {
    return StopGroup{[]int{s.Id}, s.Name, s.Desc, s.Lat, s.Lon, s.LocationType, s.Distance, s.Routes}
}


func (sgs StopGroups) ToJsonStopGroups() []JsonStopGroup {
    jsgs := make([]JsonStopGroup, len(sgs))

    for i, sg := range sgs {
        jsgs[i] = sg.ToJsonStopGroup()
        //        log.Printf("Index: %d - json: %v", i, jsgs[i])
    }

    return jsgs
}

func (rs Routes) ToJsonRoutes() []JsonRoute {
    jrs := make([]JsonRoute, len(rs))

    for i, r := range rs {
        jrs[i] = r.ToJsonRoute()
    }

    return jrs
}

func (sg StopGroup) ToJsonStopGroup() JsonStopGroup {
    return JsonStopGroup{sg.Ids, sg.Name, sg.Desc, sg.Distance, sg.LocationType, JsonGeoLocation{sg.Lat, sg.Lon}, sg.Routes.ToJsonRoutes()}
}

func (r *Route) ToJsonRoute() JsonRoute {

    jstfs := make([]string, len(r.StopTimesFull))

    for i, stf := range r.StopTimesFull {
        jstfs[i] = stf.DepartureTime
    }

    return JsonRoute{r.Name, r.TripId, r.RouteType, r.RouteColor, r.RouteTextColor, r.FirstStopName, r.LastStopName, jstfs}
}
