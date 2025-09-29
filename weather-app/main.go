package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// Struct to map JSON response
type WeatherResponse struct {
	Name string `json:"name"`
	Main struct {
		Temp     float64 `json:"temp"`
		Humidity int     `json:"humidity"`
	} `json:"main"`
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <city>")
		return
	}

	city := os.Args[1]
	apiKey := "YOUR_API_KEY" // 🔑 replace with your OpenWeatherMap key
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric", city, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println("Failed to get weather data. Status:", resp.Status)
		return
	}

	var data WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	fmt.Printf("\n🌍 Weather in %s:\n", data.Name)
	fmt.Printf("🌡️  Temperature: %.1f°C\n", data.Main.Temp)
	fmt.Printf("💧 Humidity: %d%%\n", data.Main.Humidity)
	fmt.Printf("☁️  Conditions: %s\n", data.Weather[0].Description)
}
