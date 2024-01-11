package utils

import "math"

const earthRadiusMiles = 3959

type Point struct {
	Lat, Lng float64
}

type NearBox struct {
	MinLat, MaxLat float64
	MinLng, MaxLng float64
}

// CalculateBoundingBox calculates the bounding box around a central point
func CalculateNearBox(center Point, distance float64) *NearBox {
	// Convert distance to radians
	distanceRadians := distance / earthRadiusMiles

	// Convert latitude and longitude to radians
	latRad := center.Lat * (math.Pi / 180)
	lonRad := center.Lng * (math.Pi / 180)

	// Calculate minimum and maximum latitude and longitude
	minLat := latRad - distanceRadians
	maxLat := latRad + distanceRadians
	minLng := lonRad - distanceRadians/math.Cos(latRad)
	maxLng := lonRad + distanceRadians/math.Cos(latRad)

	// Convert back to degrees
	minLat = minLat * (180 / math.Pi)
	maxLat = maxLat * (180 / math.Pi)
	minLng = minLng * (180 / math.Pi)
	maxLng = maxLng * (180 / math.Pi)

	return &NearBox{
		MinLat: minLat,
		MaxLat: maxLat,
		MinLng: minLng,
		MaxLng: maxLng,
	}
}
