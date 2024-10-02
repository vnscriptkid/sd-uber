package main

import (
	"fmt"
	"math"
	"strings"
)

// User represents a user with an ID and location.
type User struct {
	ID  int
	Lat float64
	Lon float64
}

// base32 map for geohash encoding
const base32 = "0123456789bcdefghjkmnpqrstuvwxyz"

// geohashEncode encodes latitude and longitude into a geohash string with given precision.
func geohashEncode(lat, lon float64, precision int) string {
	// Define the bit values for each position in the 5-bit encoding
	var bits = []int{16, 8, 4, 2, 1}
	// 16 = 2^4 (10000 in binary)
	// 8 = 2^3 (01000 in binary)
	// 4 = 2^2 (00100 in binary)
	// 2 = 2^1 (00010 in binary)
	// 1 = 2^0 (00001 in binary)

	// Initialize the latitude and longitude ranges
	var latRange = []float64{-90.0, 90.0}   // Valid latitude range
	var lonRange = []float64{-180.0, 180.0} // Valid longitude range

	// Initialize a string builder to construct the geohash
	var hash strings.Builder

	// Initialize variables for the encoding process
	var ch int       // Holds the current 5-bit chunk (a number between 0 and 31)
	var even = true  // Flag to alternate between longitude and latitude
	var bitIndex = 0 // Index to track position within the current 5-bit chunk

	// Continue encoding until we reach the desired precision
	// With each iteration, we narrow down the latitude or longitude range by half, effectively increasing precision.
	for hash.Len() < precision {
		if even {
			// Process longitude (even bits)
			mid := (lonRange[0] + lonRange[1]) / 2
			if lon >= mid {
				// If longitude is in the upper half, set the bit and adjust the range
				// (bitwise OR assignment, which sets the bit without affecting other bits)
				ch |= bits[bitIndex]
				lonRange[0] = mid
			} else {
				// If longitude is in the lower half, just adjust the range
				lonRange[1] = mid
			}
		} else {
			// Process latitude (odd bits)
			mid := (latRange[0] + latRange[1]) / 2
			if lat >= mid {
				// If latitude is in the upper half, set the bit and adjust the range
				ch |= bits[bitIndex]
				latRange[0] = mid
			} else {
				// If latitude is in the lower half, just adjust the range
				latRange[1] = mid
			}
		}

		// Switch between longitude and latitude for the next iteration
		even = !even

		// Move to the next bit in the 5-bit chunk
		if bitIndex < 4 {
			bitIndex++
		} else {
			// We've completed a 5-bit chunk, encode it as a base32 character
			// Append the character to the hash string
			hash.WriteByte(base32[ch])
			// Reset for the next 5-bit chunk
			ch = 0
			bitIndex = 0
		}
	}

	// Return the completed geohash string
	return hash.String()
}

// haversine calculates the great-circle distance between two points in meters.
func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371e3 // Earth radius in meters

	phi1 := lat1 * math.Pi / 180
	phi2 := lat2 * math.Pi / 180
	deltaPhi := (lat2 - lat1) * math.Pi / 180
	deltaLambda := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(deltaPhi/2)*math.Sin(deltaPhi/2) +
		math.Cos(phi1)*math.Cos(phi2)*math.Sin(deltaLambda/2)*math.Sin(deltaLambda/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := R * c // in meters

	return distance
}

// getPrecisionForRadius determines geohash precision level for given radius in meters.
func getPrecisionForRadius(radius float64) int {
	// Approximate mapping from radius to geohash precision
	switch {
	case radius >= 2500e3:
		return 1
	case radius >= 630e3:
		return 2
	case radius >= 78e3:
		return 3
	case radius >= 20e3:
		return 4
	case radius >= 2.4e3:
		return 5
	case radius >= 0.61e3:
		return 6
	case radius >= 0.076e3:
		return 7
	default:
		return 8
	}
}

func main() {
	// Sample users with IDs and locations
	users := []User{
		{ID: 1, Lat: 37.7749, Lon: -122.4194}, // San Francisco
		{ID: 2, Lat: 34.0522, Lon: -118.2437}, // Los Angeles
		{ID: 3, Lat: 40.7128, Lon: -74.0060},  // New York
		{ID: 4, Lat: 37.8044, Lon: -122.2711}, // Oakland
		{ID: 5, Lat: 37.7749, Lon: -122.4194}, // San Francisco
	}

	// Build geohash index
	geohashUsers := make(map[string][]User)
	for _, user := range users {
		// Use maximum precision for indexing
		geohash := geohashEncode(user.Lat, user.Lon, 8)
		geohashUsers[geohash] = append(geohashUsers[geohash], user)
	}

	// // User's current location
	myLat := 37.7749 // San Francisco
	myLon := -122.4194

	// // Desired radius in meters
	radius := 5000.0 // 5 km

	// // Determine geohash precision
	precision := getPrecisionForRadius(radius)

	// Get my geohash at the determined precision
	myGeohash := geohashEncode(myLat, myLon, precision)

	fmt.Printf(">> My Geohash: %s\n", myGeohash)

	nearbyUsers := []User{}

	for hash, users := range geohashUsers {
		fmt.Printf("Geohash: %s, Users: %v\n", hash, users)

		if strings.HasPrefix(hash, myGeohash) {
			fmt.Printf(">> found a prefix match\n")
			nearbyUsers = append(nearbyUsers, users...)
		}
	}

	// Filter users within the radius
	result := []User{}
	for _, user := range nearbyUsers {
		d := haversine(myLat, myLon, user.Lat, user.Lon)
		if d <= radius {
			result = append(result, user)
		}
	}

	fmt.Printf("Result: %v\n", result)
}
