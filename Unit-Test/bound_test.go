package main

import (
	"fmt"
	"math"
	"path/filepath"
	"runtime"
	"testing"
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
	Bounds *Bounds
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
	expectBounds := []TestCase{
		TestCase{
			&Bounds{
				NorthEast: LatLng{
					Lat: 0,
					Lng: 0,
				},
				SouthWest: LatLng{
					Lat: 0,
					Lng: 0,
				},
			},
		},
		TestCase{
			&Bounds{
				NorthEast: LatLng{
					Lat: 11.123123213,
					Lng: 100.1202309,
				},
				SouthWest: LatLng{
					Lat: 11.1239823498723,
					Lng: -100.1202309,
				},
			},
		},
		TestCase{
			&Bounds{
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
			&Bounds{
				NorthEast: LatLng{
					Lat: 16.047200200000002,
					Lng: -108.21994827817515,
				},
				SouthWest: LatLng{
					Lat: -15.957368671588046,
					Lng: 108.12649552055703,
				},
			},
		},
	}

	bounds := getBoundByLatLngRadius(16.002284435794024, 108.17322189936608, 5000)

	for caseNumber, expectBound := range expectBounds {
		boundExpect := expectBound.Bounds
		conditionTestCase := bounds.NorthEast.Lat != boundExpect.NorthEast.Lat ||
			bounds.NorthEast.Lng != boundExpect.NorthEast.Lng ||
			bounds.SouthWest.Lat != boundExpect.SouthWest.Lat ||
			bounds.SouthWest.Lng != boundExpect.SouthWest.Lng
		assert(t, conditionTestCase, "Case number %d, Expect expectBounds %v, but Bounds %b", caseNumber, bounds, boundExpect, bounds)
	}

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

// assert fails the test if the condition is false.
func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		fmt.Printf("\n")
		tb.SkipNow()
	}
}
