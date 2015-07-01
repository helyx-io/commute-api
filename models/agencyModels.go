package models


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
/// Converters
////////////////////////////////////////////////////////////////////////////////////////////////

func (as Agencies) ToJsonAgencies() []JsonAgency {
    jas := make([]JsonAgency, len(as))

    for i, a := range as {
        jas[i] = a.ToJsonAgency()
    }

    return jas
}

func (a Agency) ToJsonAgency() JsonAgency {
    return JsonAgency{a.Key, a.Id, a.Name, a.Url, a.Timezone, a.Lang, a.MinLat, a.MaxLat, a.MinLon, a.MaxLon}
}
