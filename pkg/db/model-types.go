package db

// location: {
// 	type: "Point",
// 	coordinates: [-73.856077, 40.848447]
// }

// GeoSpatial - For location coordinates
type GeoSpatial struct {
	Type        string    `bson:"type"`
	Coordinates []float64 `bson:"coordinates"`
	Location    string    `bson:"-"`
}

// NewGeoPoint - Create a geo point
func NewGeoPoint(longitude, latitude float64) *GeoSpatial {
	return &GeoSpatial{
		Type:        "Point",
		Coordinates: []float64{longitude, latitude},
	}
}

// GetLongitude - to get longitude from geo spatial coordinate
func (t GeoSpatial) GetLongitude() float64 {
	if len(t.Coordinates) == 2 {
		return t.Coordinates[0]
	}
	return 0
}

// GetLatitude - to get longitude from geo spatial coordinate
func (t GeoSpatial) GetLatitude() float64 {
	if len(t.Coordinates) == 2 {
		return t.Coordinates[1]
	}
	return 0
}
