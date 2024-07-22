package location

type Location struct {
	latitude  float64
	longitude float64
}

func (l Location) Latitude() float64  { return l.latitude }
func (l Location) Longitude() float64 { return l.longitude }

func (l Location) ToJSON() JSONLocation {
	return JSONLocation{
		Latitude:  l.latitude,
		Longitude: l.longitude,
	}
}
func NewLocation(latitude, longitude float64) Location {
	return Location{latitude: latitude, longitude: longitude}
}

type JSONLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (j JSONLocation) ToLocation() Location {
	return Location{
		latitude:  j.Latitude,
		longitude: j.Longitude,
	}
}
