package service

//Response type used to represent final response
type Response struct {
	Location   string     `json:"location_name"`
	Temp       string     `json:"temperature"`
	Feel       string     `json:"real_feel_temperature"`
	Min        string     `json:"minimum_temperature"`
	Max        string     `json:"maximum_temperature"`
	Wind       string     `json:"wind"`
	Cloudiness string     `json:"cloudiness"`
	Pressure   string     `json:"pressure"`
	Humidity   string     `json:"humidity"`
	Sunrise    string     `json:"sunrise"`
	Sunset     string     `json:"sunset"`
	Coord      string     `json:"geo_coordinates"`
	ReqTime    string     `json:"requested_time"`
	Forecast   []forecast `json:"forecast"`
}

type forecast struct {
	ForecastedDate string `json:"forecasted_datetime"`
	Temp           string `json:"temperature"`
	Feel           string `json:"real_feel_temperature"`
	Min            string `json:"minimum_temperature"`
	Max            string `json:"maximum_temperature"`
	Cloudiness     string `json:"cloudiness"`
	Humidity       string `json:"humidity"`
}

type mainWeatherInfo struct {
	FeelsLike float64 `json:"feels_like"`
	Humidity  int     `json:"humidity"`
	Pressure  int     `json:"pressure"`
	Temp      float64 `json:"temp"`
	TempMax   float64 `json:"temp_max"`
	TempMin   float64 `json:"temp_min"`
}

type cloudInfo struct {
	Description string `json:"description"`
}

type weatherResponse struct {
	Coord struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	} `json:"coord"`
	Main mainWeatherInfo `json:"main"`
	Name string          `json:"name"`
	Sys  struct {
		Country string `json:"country"`
		Sunrise int    `json:"sunrise"`
		Sunset  int    `json:"sunset"`
	} `json:"sys"`
	Weather []cloudInfo `json:"weather"`
	Wind    struct {
		Deg   int     `json:"deg"`
		Speed float64 `json:"speed"`
	} `json:"wind"`
}

type forecastResponse struct {
	Forecast []struct {
		Dt      int             `json:"dt"`
		Main    mainWeatherInfo `json:"main"`
		Weather []cloudInfo     `json:"weather"`
	} `json:"list"`
}
