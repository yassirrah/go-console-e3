package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"https://github.com/yassirrah/weather-app/internal/weather"

	"github.com/gin-gonic/gin"
)

func main() {
	apiKey := strings.TrimSpace(os.Getenv("OPENWEATHER_API_KEY"))
	if apiKey == "" {
		log.Fatal("OPENWEATHER_API_KEY is required")
	}
	wc := weather.NewClient(apiKey)

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "ts": time.Now().UTC()})
	})

	// GET /api/v1/weather?city=Rabat[&units=metric|imperial]
	r.GET("/api/v1/weather", func(c *gin.Context) {
		city := strings.TrimSpace(c.Query("city"))
		if city == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing city query param"})
			return
		}
		units := strings.TrimSpace(c.DefaultQuery("units", "metric"))

		res, err := wc.GetByCity(city, units)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, weather.ToDTO(res))
	})

	addr := ":8080"
	log.Printf("listening on %s ...", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
