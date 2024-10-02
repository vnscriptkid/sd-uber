package main

import (
	"testing"
)

func TestGeohashEncode(t *testing.T) {
	testCases := []struct {
		name      string
		lat       float64
		lon       float64
		precision int
		expected  string
	}{
		{"Eiffel Tower", 48.8584, 2.2945, 6, "u09tun"},
		{"Louvre Museum", 48.8606, 2.3376, 6, "u09tvn"},
		{"Notre-Dame Cathedral", 48.852968, 2.349902, 6, "u09tvm"},
		{"Tokyo Tower", 35.6586, 139.7454, 6, "xn76gg"},
		{"Senso-ji Temple", 35.7148, 139.7967, 6, "xn77jj"},
		{"Shibuya Crossing", 35.6595, 139.7006, 6, "xn76fg"},
		{"Brandenburg Gate", 52.5163, 13.3777, 6, "u33db2"},
		{"Neuschwanstein Castle", 47.5576, 10.7498, 6, "u0rws9"},
		{"The Great Wall", 40.4319, 116.5704, 6, "wx4yh8"},
		{"The Forbidden City", 39.9163, 116.3972, 6, "wx4g0d"},
		{"Terracotta Army", 34.3833, 109.2772, 6, "wqjewe"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := geohashEncode(tc.lat, tc.lon, tc.precision)
			if result != tc.expected {
				t.Errorf("geohashEncode(%f, %f, %d) = %s; want %s",
					tc.lat, tc.lon, tc.precision, result, tc.expected)
			}
		})
	}
}
