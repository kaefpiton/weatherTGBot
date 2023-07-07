package api

type Weather struct {
	Temperature          float64
	TemperatureFeelsLike float64
	Pressure             float64
	WindSpeed            float64
	Clouds               int
	Humidity             int
}
