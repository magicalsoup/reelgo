package trips

type TripAttractions struct {
	Tid      int64 `sql:"primary_key" json:"tid"`
	TripName string `json:"trip_name"`
	Attractions [] AttractionPayload `json:"attractions"`
}

type AttractionPayload struct {
	Aid                int64 `sql:"primary_key" json:"aid"`
	AttractionName     string `json:"attraction_name"`
	AttractionLocation string `json:"attraction_location"`
}