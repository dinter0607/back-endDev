package main

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Bounds struct {
	NorthEast LatLng
	SouthWest LatLng
}
type LatLng struct {
	Lat float64
	Lng float64
}

var (
	earthRadius = 6378137.0 // meters

	//MinLatitude is the minimum possible latitude
	minLatitude = deg2rad(-90)

	//MaxLatitude is the maxiumum possible latitude
	maxLatitude = deg2rad(90)

	//MinLongitude is the minimum possible longitude
	minLongitude = deg2rad(-180)

	//MaxLongitude is the maxiumum possible longitude
	maxLongitude = deg2rad(180)
)

func deg2rad(d float64) float64 {
	return d * math.Pi / 180.0
}

func rad2deg(r float64) float64 {
	return 180.0 * r / math.Pi
}

func TestGetBounds(t *testing.T) {
	assert := assert.New(t)
	expectBounds := &Bounds{
		NorthEast: LatLng{
			Lat: 16.047200200000002,
			Lng: 108.21994827817515,
		},
		SouthWest: LatLng{
			Lat: 15.957368671588046,
			Lng: 108.12649552055703,
		},
	}

	bounds := getBoundByLatLngRadius(16.002284435794024, 108.17322189936608, 5000)

	// assert equality
	assert.Equal(bounds, expectBounds, "they should be equal")

}

func getBoundByLatLngRadius(lat, lng float64, radius float64) *Bounds {
	if lat < 0 || lng < 0 || radius < 0 {
		return nil
	}

	radDist := radius / earthRadius
	radLat := deg2rad(lat)
	radLon := deg2rad(lng)
	minLat := radLat - radDist
	maxLat := radLat + radDist

	var minLng, maxLng float64
	if minLat > minLatitude && maxLat < maxLatitude {
		deltaLon := math.Asin(math.Sin(radDist) / math.Cos(radLat))
		minLng = radLon - deltaLon
		if minLng < minLongitude {
			minLng += 2 * math.Pi
		}
		maxLng = radLon + deltaLon
		if maxLng > maxLongitude {
			maxLng -= 2 * math.Pi
		}
	} else {
		minLat = math.Max(minLat, minLatitude)
		maxLat = math.Min(maxLat, maxLatitude)
		minLng = minLongitude
		maxLng = maxLongitude
	}

	bounds := &Bounds{
		NorthEast: LatLng{
			Lat: rad2deg(maxLat),
			Lng: rad2deg(maxLng),
		},
		SouthWest: LatLng{
			Lat: rad2deg(minLat),
			Lng: rad2deg(minLng),
		},
	}

	return bounds
}
