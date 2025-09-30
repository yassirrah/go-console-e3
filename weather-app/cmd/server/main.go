package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yassirrah/go-console-e3/weather-app/internal/store"
	"github.com/yassirrah/go-console-e3/weather-app/internal/weather"
)

func main() {
	apiKey := strings.TrimSpace(os.Getenv("OPENWEATHER_API_KEY"))
	if apiKey == "" {
		log.Fatal("OPENWEATHER_API_KEY is required")
	}

	dsn := strings.TrimSpace(os.Getenv("DB_DSN"))
	if dsn == "" {
		log.Fatal("DB_DSN is required (e.g. postgres://user:pass@db:5432/weather?sslmode=disable)")
	}

	st, err := store.Open(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer st.Close()

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

		dto := weather.ToDTO(res)
		// Persist
		_ = st.InsertWeather(c.Request.Context(), store.WeatherRow{
			City: city, Units: units, Temperature: dto.Temperature,
			Humidity: dto.Humidity, Conditions: dto.Conditions,
		}) // keep simple; handle/log err in real app

		c.JSON(http.StatusOK, dto)
	})

	addr := ":8080"
	log.Printf("listening on %s ...", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
