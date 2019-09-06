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
