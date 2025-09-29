package weather

// Upstream (trimmed)
type Response struct {
	Name string `json:"name"`
	Main struct {
		Temp     float64 `json:"temp"`
		Humidity int     `json:"humidity"`
	} `json:"main"`
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
}

// Our API contract
type DTO struct {
	City        string  `json:"city"`
	Temperature float64 `json:"temperature_c"`
	Humidity    int     `json:"humidity_pct"`
	Conditions  string  `json:"conditions"`
	Source      string  `json:"source"`
}

func ToDTO(r Response) DTO {
	desc := ""
	if len(r.Weather) > 0 {
		desc = r.Weather[0].Description
	}
	return DTO{
		City:        r.Name,
		Temperature: r.Main.Temp,
		Humidity:    r.Main.Humidity,
		Conditions:  desc,
		Source:      "openweathermap",
	}
}
