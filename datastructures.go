package main

/*
{"timestamp": uint32, "box": int, "sensor": int, "value": float}
*/

type PowerMsg struct {
  Timestamp int64  `json:"timestamp"`
  Box int           `json:"box"`
  Voltage float64   `json:"voltage"`
  Percent float64  `json:"percent"`
}

type Record struct {
  Timestamp int64 `json:"timestamp"`
  Box int         `json:"box"`
  Sensor int      `json:"sensor"`
  Value float64   `json:"value"`
  VWC float64     `json:"vwc"`
}

type Records []Record

type Sensor struct {
  Number int `json:"slot"`
  Dots []Dot `json:"dots"`
  Water int `json:"water"`
}

type Sensors []Sensor

type Dot struct {
  Timestamp int64 `json:"timestamp"`
  Value float64   `json:"value"`
}

type UbidotsAPI struct {
  Count bool  `json:"count"`
  Next string `json:"next"`
  Previous string `json:"previous"`
  Results []Dot `json:"results"`
}

type Device struct {
  Key string
  Id string       `json:"device" redis:"id"`
  NumSensors int  `json:"sensors" redis:"sensors"`
  Box int         `json:"box" redis:"box"`
  Watering int    `json:"watering" redis:"watering"`
  Valve int       `json:"valve" redis:"valve"`
  Sleep int       `json:"sleep" redis:"sleep"`
}

type Garden struct {
  Devices []Device `json:"garden"`
  NumDevices int  `json:"devices"`
}
