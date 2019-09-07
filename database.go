package main

import (
  "github.com/influxdata/influxdb1-client/v2"
  "database/sql"
  _ "github.com/lib/pq"
  "fmt"
)

const (
  host = "openstack01"
  port = 5432
  user = "postgres"
  password = "password01"
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

func initializePQConnection() (driver.Conn, error) {
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    return nil, err
  }
  return db, nil
}
