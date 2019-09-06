package main

import (
        "log"
        "io/ioutil"
        "net/http"
        "net/url"
        "encoding/json"
        "fmt"
)

type Wind struct {
    Speed float64 `json:"speed"`
    Deg   float64 `json:"deg"`
}

type Rain struct {
    ThreeH float64 `json:"3h"`
}

type Snow struct {
    ThreeH float64 `json:"3h"`
}

type Main struct {
    Temp      float64 `json:"temp"`
    TempMin   float64 `json:"temp_min"`
    TempMax   float64 `json:"temp_max"`
    Pressure  float64 `json:"pressure"`
    SeaLevel  float64 `json:"sea_level"`
    GrndLevel float64 `json:"grnd_level"`
    Humidity  int     `json:"humidity"`
}

type Clouds struct {
    All int `json:"all"`
}

type Weather struct {
    ID          int    `json:"id"`
    Main        string `json:"main"`
    Description string `json:"description"`
    Icon        string `json:"icon"`
}

type ForecastWeatherList struct {
    Dt      int       `json:"dt"`
    Main    Main      `json:"main"`
    Weather []Weather `json:"weather"`
    Clouds  Clouds    `json:"clouds"`
    Wind    Wind      `json:"wind"`
    Rain    Rain      `json:"rain"`
    Snow    Snow      `json:"snow"`
    DtTxt   string     `json:"dt_txt"`
}

type Coordinates struct {
    Longitude float64 `json:"lon"`
    Latitude  float64 `json:"lat"`
}

type ForecastSys struct {
    Population int `json:"population"`
}

type City struct {
    ID         int         `json:"id"`
    Name       string      `json:"name"`
    Coord      Coordinates `json:"coord"`
    Country    string      `json:"country"`
    Population int         `json:"population"`
    Sys        ForecastSys `json:"sys"`
}

type ForecastWeatherData struct {
    // COD     string                `json:"cod"`
    // Message float64               `json:"message"`
    City City                   `json:"city"`
    Cnt  int                    `json:"cnt"`
    List []ForecastWeatherList `json:"list"`
}

func getWeatherData() ForecastWeatherData {
  var (
    apiKey = "f014a6a37b52a589b4e7254a70372e54"
    id = "5913490"
  )
  urlString := fmt.Sprintf("http://api.openweathermap.org/data/2.5/forecast?id=%s&APPID=%s", id, apiKey)
  u, err := url.Parse(urlString)
  res, err := http.Get(u.String())
  if err != nil {
    log.Fatal(err)
  }

  jsonBlob, err := ioutil.ReadAll(res.Body)
  res.Body.Close()
  if err != nil {
    log.Fatal(err)
  }

  var data ForecastWeatherData
  err = json.Unmarshal(jsonBlob, &data)
  if err != nil {
    fmt.Println("error:", err)
  }
  return data
}

func extractWeatherData(data []ForecastWeatherList) []int {
  weather := make([]int, len(data))
  for i, value := range data {
    weather[i] = value.Weather[0].ID
  }
  //fmt.Println("extractWeatherData: the extracted weather codes are: ", weather)
  return weather
}

//returns true if rain/drizzle/snow in the forecast
//returns false if cloud/clear sky/thunderstorm in the forecast
func weatherDataAnalysis(data ForecastWeatherData) bool {
  decision := 0
  weather := extractWeatherData(data.List)
  for i := 0; i < 8; i++ {
    fmt.Println("weatherDataAnalysis: i = ", i, " and weather[i] = ", weather[i], " and decision = ", decision)
    switch {
    // Thunderstorm weather codes
    case weather[i] >= 200 && weather[i] <= 299:
      decision += 1
    // Drizzle weather codes
    case weather[i] >= 300 && weather[i] <= 399:
      decision -= 1
    // Rain weather codes
    case weather[i] >= 500 && weather[i] <= 599:
      decision -= 1
    // Snow weather code
    case weather[i] >= 600 && weather[i] <= 699:
      decision -= 1
    default:
      decision += 1
    }
  }
  fmt.Println("weatherDataAnalysis: based on this weather data: ", weather, " the decision is ", decision)
  //if decision = 6 or lower that means we have at least 1 3hr period of rain/drizzle/snow in next 24hrs, therefore we return true
  if decision < 6 {
    return true
  } else {
    return false
  }
  return false
}
