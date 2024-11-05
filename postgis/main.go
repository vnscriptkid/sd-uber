package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

type Driver struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Location string `gorm:"type:geometry(Point,4326)" json:"-"`
}

func initDB() {
	var err error
	dsn := "host=localhost user=postgres dbname=postgres sslmode=disable password=123456"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.Exec("CREATE EXTENSION IF NOT EXISTS postgis")

	db.AutoMigrate(&Driver{})

	// Create a spatial index on the location column
	db.Exec("CREATE INDEX idx_drivers_location ON drivers USING GIST(location)")

}

func CreateDriver(c *gin.Context) {
	var input struct {
		Latitude  float64 `json:"latitude" binding:"required"`
		Longitude float64 `json:"longitude" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wkt := fmt.Sprintf("POINT(%f %f)", input.Longitude, input.Latitude)

	driver := Driver{
		Location: wkt,
	}

	if err := db.Create(&driver).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, driver)
}

func UpdateDriverLocation(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid driver ID"})
		return
	}

	var input struct {
		Latitude  float64 `json:"latitude" binding:"required"`
		Longitude float64 `json:"longitude" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wkt := fmt.Sprintf("POINT(%f %f)", input.Longitude, input.Latitude)

	if err := db.Model(&Driver{}).Where("id = ?", id).Update("location", gorm.Expr("ST_GeomFromText(?, 4326)", wkt)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id, "location": wkt})
}

func FindNearbyDrivers(c *gin.Context) {
	latitudeParam := c.Query("latitude")
	longitudeParam := c.Query("longitude")
	radiusParam := c.Query("radius")

	latitude, err1 := strconv.ParseFloat(latitudeParam, 64)
	longitude, err2 := strconv.ParseFloat(longitudeParam, 64)
	radius, err3 := strconv.ParseFloat(radiusParam, 64)

	if err1 != nil || err2 != nil || err3 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parameters"})
		return
	}

	point := fmt.Sprintf("SRID=4326;POINT(%f %f)", longitude, latitude)

	var drivers []struct {
		ID        uint    `json:"id"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Distance  float64 `json:"distance"`
	}

	err := db.Raw(`
        SELECT id,
               ST_Y(location::geometry) AS latitude,
               ST_X(location::geometry) AS longitude,
               ST_DistanceSphere(location, ST_GeomFromText(?, 4326)) AS distance
        FROM drivers
        WHERE ST_DWithin(location::geography, ST_GeomFromText(?, 4326)::geography, ?)
        ORDER BY distance
    `, point, point, radius*1000).Scan(&drivers).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, drivers)
}

func main() {
	initDB()
	r := gin.Default()

	r.POST("/drivers", CreateDriver)
	r.PUT("/drivers/:id/location", UpdateDriverLocation)
	r.GET("/drivers/nearby", FindNearbyDrivers)

	r.Run(":8080")
}
