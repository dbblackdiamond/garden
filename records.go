package main

import (
  "github.com/influxdata/influxdb1-client/v2"
  "encoding/json"
  "strconv"
  "fmt"
  "github.com/garyburd/redigo/redis"
  "time"
)

func analyzeRecords(c client.Client, results []client.Result) []Record {
  records := make([]Record, len(results[0].Series[0].Values))
  for i, value := range results[0].Series[0].Values {
//    fmt.Println("Record #",i ,":\tBox = ", value[1], "\tSensor = ", value[2], "\tValue = ", value[3], "\tTimestamp = ", value[0])
    timestamp, _ := time.Parse(time.RFC3339, value[0].(string))
    records[i].Timestamp = int64(timestamp.Unix())
    box, _ := value[1].(string)
    records[i].Box, _ = strconv.Atoi(box)
    sensor, _ := value[2].(string)
    records[i].Sensor, _ = strconv.Atoi(sensor)
    records[i].Value, _ = value[3].(json.Number).Float64()
  }
  return records
}

func getAllRecords(c client.Client) ([]client.Result, error) {
  q := client.Query{
    Command: "select * from moisture",
    Database: database,
  }

  resp, err := c.Query(q)
  if err != nil {
    return nil, err
  }
  if resp.Error() != nil {
    return nil, resp.Error()
  }
  return resp.Results, nil
}

func getBoxRecords(c client.Client, box int) (res []client.Result, err error) {
  q := client.Query{
    Command: fmt.Sprintf("select * from moisture where box = '%d' AND time > now() - 24h tz('America/Edmonton')", box),
    Database: database,
  }

  resp, err := c.Query(q)
  if err != nil {
    return nil, err
  }
  if resp.Error() != nil {
    return nil, resp.Error()
  }
  return resp.Results, nil
}

/*func getListOfBoxes(c client.Client) ([]client.Result, error) {
  q := client.Query{
    Command: "select distinct(box) from (select box, value from moisture)",
    Database: database,
  }

  resp, err := c.Query(q)
  if err != nil {
    return nil, err
  }
  if resp.Error() != nil {
    return nil, resp.Error()
  }
  return resp.Results, nil
}*/

func loadGarden() {
  c := Pool.Get()
  defer c.Close()
  res, err := redis.Strings(c.Do("KEYS", "photon:*"))
  fmt.Println("loadGarden:\tGetting those keys from redis:", res)
  check("Keys", err)

  garden.Devices = make([]Device, len(res))
  garden.NumDevices = len(res)
  for index, photon := range res {
    r, err := redis.Values(c.Do("HGETALL", photon))
    err = redis.ScanStruct(r, &(garden.Devices[index]))
    garden.Devices[index].Key = photon
    check("ScanStruct", err)
  }
  fmt.Println("loadGarden:\tgarden = ", garden)
}
