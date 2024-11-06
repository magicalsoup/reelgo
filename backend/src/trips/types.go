package trips

type TripAttractions struct {
	Tid int64 `sql:"primary_key"`
	TripName string
	Attractions []struct {
		Aid int64 `sql:"primary_key"`
		AttractionName string
		AttractionLocation string
	}
}