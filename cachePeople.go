package main

import (
  "github.com/kellydunn/golang-geo"
  "time"
  "fmt"
)

const peopleExpireDuration = time.Second * 10

type People struct {
  Lat, Lon float64
  Nick string
  AliveTime time.Time
}

type PeopleManager struct {
  People []People
}

func (pm *PeopleManager) PutIfNeeded(p *People) {
  for _, peopleInList := range pm.People {
    if(peopleInList.Nick == p.Nick) {
      return
    }
  }
  pm.People = append(pm.People, *p)
}

func (pm *PeopleManager) SumPeopleInArea(lat, lon float64)int {
  central := geo.NewPoint(lat, lon)

  up := central.PointAtDistanceAndBearing(0.5, 0)
  down := central.PointAtDistanceAndBearing(0.5, 180)
  left := central.PointAtDistanceAndBearing(0.5, 270)
  right := central.PointAtDistanceAndBearing(0.5, 90)

  sum := 0
  for _, peopleInList := range pm.People {
    if(peopleInList.Lat <= up.Lat() &&
      peopleInList.Lat >= down.Lat() &&
      peopleInList.Lon >= left.Lng() &&
      peopleInList.Lon <= right.Lng()) {
          sum++
    }
  }
  return sum
}

func (pm *PeopleManager) Clean(delay time.Duration) chan bool {
  stop := make(chan bool)

  go func(delay time.Duration) {
    for {
      select {
      case <- time.After(delay):
        fmt.Println("Cleaning people")
        now :=  time.Now()

        for i := len(pm.People)-1; i >= 0; i-- {
          people := pm.People[i]
          if now.Sub(people.AliveTime) > peopleExpireDuration {
          pm.People = append(pm.People[:i], pm.People[i + 1:]...)
          }
        }
      case <- stop:
            return
      }
    }
  }(delay)

  return stop
}
