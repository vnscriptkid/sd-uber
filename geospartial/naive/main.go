package main

import (
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// User represents a user with latitude and longitude coordinates.
type User struct {
	ID        uint `gorm:"primaryKey"`
	Latitude  float64
	Longitude float64
}

var (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "123456"
	dbname   = "postgres"
)

func main() {
	// Database connection string.
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai", host, user, password, dbname, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the User schema.
	if err := db.AutoMigrate(&User{}); err != nil {
		panic("failed to migrate database schema")
	}

	// Create indexes on latitude and longitude for performance.
	db.Exec("CREATE INDEX IF NOT EXISTS idx_users_latitude ON users(latitude)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_users_longitude ON users(longitude)")

	r := gin.Default()

	// Endpoint to find nearby users within a radius.
	r.GET("/nearby", func(c *gin.Context) {
		// Get query parameters.
		latStr := c.Query("lat")
		lonStr := c.Query("lon")
		radiusStr := c.Query("radius")

		// Convert parameters to float64.
		lat, err := strconv.ParseFloat(latStr, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid latitude"})
			return
		}
		lon, err := strconv.ParseFloat(lonStr, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid longitude"})
			return
		}
		radiusKm, err := strconv.ParseFloat(radiusStr, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid radius"})
			return
		}

		// Find nearby users.
		nearbyUsers, err := findNearbyUsers(db, lat, lon, radiusKm)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		c.JSON(http.StatusOK, nearbyUsers)
	})

	r.Run()
}

// findNearbyUsers finds users within a given radius (in km) from the specified latitude and longitude.
func findNearbyUsers(db *gorm.DB, lat, lon, radiusKm float64) ([]User, error) {
	// Earth's radius in kilometers.
	const earthRadius = 6371.0

	// Calculate latitude and longitude boundaries for the bounding box.
	latRadius := radiusKm / earthRadius * (180 / math.Pi)
	minLat := lat - latRadius
	maxLat := lat + latRadius

	// Adjust longitude boundaries based on the latitude.
	lonRadius := radiusKm / (earthRadius * math.Cos(lat*math.Pi/180)) * (180 / math.Pi)
	minLon := lon - lonRadius
	maxLon := lon + lonRadius

	var users []User

	// Raw SQL query to find users within the bounding box and specified radius.
	query := `
		SELECT id, latitude, longitude, 
			(6371 * acos(
				cos(radians(?)) * cos(radians(latitude)) * cos(radians(longitude) - radians(?)) + 
				sin(radians(?)) * sin(radians(latitude))
			)) AS distance
		FROM users
		WHERE (latitude BETWEEN ? AND ?) AND (longitude BETWEEN ? AND ?)
			AND (6371 * acos(
				cos(radians(?)) * cos(radians(latitude)) * cos(radians(longitude) - radians(?)) + 
				sin(radians(?)) * sin(radians(latitude))
			)) < ?
		ORDER BY distance
	`

	// Execute the query with the provided parameters.
	err := db.Raw(query, lat, lon, lat, minLat, maxLat, minLon, maxLon, lat, lon, lat, radiusKm).Scan(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}
