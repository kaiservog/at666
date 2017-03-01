package main

import (
  "github.com/kellydunn/golang-geo"
)

type Coordinate struct {
  Lat, Lon float64
}

const Distance = 0.5

func NewCoordinate(lat, lon float64) *Coordinate {
  return &Coordinate{lat, lon}
}

func (c *Coordinate) GetSquare() (up, down, left, right float64) {
  central := geo.NewPoint(c.Lat, c.Lon)

  up = central.PointAtDistanceAndBearing(Distance, 0).Lat()
  down = central.PointAtDistanceAndBearing(Distance, 180).Lat()
  left = central.PointAtDistanceAndBearing(Distance, 270).Lng()
  right = central.PointAtDistanceAndBearing(Distance, 90).Lng()

  return  up, down, left, right
}
