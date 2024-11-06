package trips

import (
	"database/sql"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/magicalsoup/reelgo/.gen/reelgo/public/model"
	. "github.com/magicalsoup/reelgo/.gen/reelgo/public/table"
	"github.com/magicalsoup/reelgo/src/gcs"
)


func AddAttraction(db *sql.DB, attraction gcs.Attraction, user *model.Users) (error) {
	trip_name, err := gcs.GenerateTripName(attraction)
	if err != nil {
		return err
	}

	get_trip_stmt := Trips.SELECT(Trips.AllColumns).WHERE(Trips.UID.EQ(Int32(user.UID)).AND(Trips.TripName.EQ(String(trip_name)))).LIMIT(1)
	
	trip := &model.Trips{}

	err = get_trip_stmt.Query(db, trip)


	// no existing trip, we add it to database
	if err == qrm.ErrNoRows {
		insert_trip_stmt := Trips.INSERT(Trips.UID, Trips.TripName).VALUES(user.UID, trip_name).RETURNING(Trips.AllColumns)
		err = insert_trip_stmt.Query(db, trip)

		if err != nil {
			return err
		}

	} else if err != nil { // some database error
		return err
	} 

	// we can be for sure trip exists
	insert_attraction_stmt := Attractions.INSERT(Attractions.UID, Attractions.Tid, Attractions.AttractionName, Attractions.AttractionLocation).VALUES(user.UID, trip.Tid, attraction.Name, attraction.Location)
	_, err = insert_attraction_stmt.Exec(db)
	return nil
}

func GetTrips(db *sql.DB, user *model.Users) ([]TripAttractions, error) {
	stmt := SELECT(	Trips.Tid, 
					Trips.TripName, 
					Attractions.Aid, 
					Attractions.AttractionName, 
					Attractions.AttractionLocation,
			).FROM(
				Trips.INNER_JOIN(Attractions, Attractions.UID.EQ(Trips.UID).AND(Attractions.Tid.EQ(Trips.Tid))),
			).WHERE(Trips.UID.EQ(Int32(user.UID))).
			ORDER_BY(Trips.Tid, Attractions.Aid)
	
	trips := []TripAttractions{}

	err := stmt.Query(db, &trips)

	if err != nil {
		return nil, err
	}

	return trips, nil
}