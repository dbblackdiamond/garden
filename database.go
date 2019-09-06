package main

import (
  "github.com/influxdata/influxdb1-client/v2"
)

const (
  database = "gardentestdb"
)

func initializeDBConnection() (client.Client, error) {
  c, err := client.NewHTTPClient(client.HTTPConfig {
    Addr: "http://192.168.1.100:8086",
  })
  if err != nil {
    return nil, err
  }
  return c, nil
}
