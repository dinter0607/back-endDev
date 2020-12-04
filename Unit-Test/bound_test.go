package main

import (
	"math"
	"reflect"
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

type TestCase struct {
	Lat          float64
	Lng          float64
	Radius       float64
	ExpectBounds *Bounds
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
	testSuite := []TestCase{
		TestCase{
			Lat:    16.002284435794024,
			Lng:    108.17322189936608,
			Radius: 5000.0,
			ExpectBounds: &Bounds{
				NorthEast: LatLng{
					Lat: 16.047200200000002,
					Lng: 108.21994827817515,
				},
				SouthWest: LatLng{
					Lat: 15.957368671588046,
					Lng: 108.12649552055703,
				},
			},
		},
		TestCase{
			Lat:    0,
			Lng:    0,
			Radius: 5000.0,
			ExpectBounds: &Bounds{
				NorthEast: LatLng{
					Lat: 0.04491576420597608,
					Lng: 0.04491576420597608,
				},
				SouthWest: LatLng{
					Lat: -0.04491576420597608,
					Lng: -0.04491576420597608,
				},
			},
		},
		TestCase{
			Lat:    11.3966303,
			Lng:    106.8267534,
			Radius: 5000.0,
			ExpectBounds: &Bounds{
				NorthEast: LatLng{
					Lat: 11.441546064205975,
					Lng: 106.87257259064775,
				},
				SouthWest: LatLng{
					Lat: 11.351714535794024,
					Lng: 106.78093420935225,
				},
			},
		},
	}
	for caseNumber, testCase := range testSuite {
		bounds := getBoundByLatLngRadius(testCase.Lat, testCase.Lng, testCase.Radius)

		assert.Equal(t, testCase.ExpectBounds.NorthEast.Lat, bounds.NorthEast.Lat, "Case number %d, should have Equal NorthEast latitude", caseNumber)
		assert.Equal(t, testCase.ExpectBounds.NorthEast.Lng, bounds.NorthEast.Lng, "Case number %d, should have Equal NorthEast longitude", caseNumber)
		assert.Equal(t, testCase.ExpectBounds.SouthWest.Lat, bounds.SouthWest.Lat, "Case number %d, should have Equal NorthEast latitude", caseNumber)
		assert.Equal(t, testCase.ExpectBounds.SouthWest.Lng, bounds.SouthWest.Lng, "Case number %d, should have Equal NorthEast longitude", caseNumber)
	}

}

func getBoundByLatLngRadius(lat, lng float64, radius float64) *Bounds {
	if radius <= 0 {
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

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		return
	}
	// debug.PrintStack()
	t.Errorf("Received %v (type %v), expected %v (type %v)", a, reflect.TypeOf(a), b, reflect.TypeOf(b))
}
