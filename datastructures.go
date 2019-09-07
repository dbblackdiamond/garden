package main

/*
{"timestamp": uint32, "box": int, "sensor": int, "value": float}
*/

type PowerMsg struct {
  Box int           `json:"box"`
  Capacity int  `json:"capacity"`
  Charge int    `json:"charge"`
  Current int   `json:"current"`
  Health int    `json:"health"`
  Voltage float64   `json:"voltage"`
  Signal int       `json:"signal"`
}

type Record struct {
  Box int         `json:"box"`
  Sensor int      `json:"sensor"`
  Value float64   `json:"value"`
  VWC float64     `json:"vwc"`
}

type Records []Record

type Sample struct {
  Timestamp int64 `json:"t"`
  Value float64 `json:"v"`
}

type Samples []Sample

type Sensor struct {
  Number int `json:"slot"`
  Samples []Sample `json:"samples"`
  Water int `json:"water"`
}

type Sensors []Sensor

type Device struct {
  Key string
  Id string       `json:"device" redis:"id"`
  NumSensors int  `json:"sensors" redis:"sensors"`
  Box int         `json:"box" redis:"box"`
  Valve int       `json:"valve" redis:"valve"`
  Sleep int       `json:"sleep" redis:"sleep"`
  Start int       `json:"start" redis:"start"`
  End into        `json:"end" redis:"end"`
}

type Garden struct {
  Devices []Device `json:"garden"`
  NumDevices int  `json:"devices"`
}
