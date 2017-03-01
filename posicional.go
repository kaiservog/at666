package main

type IsCoordinateInsideFunc func(posicional Posicional, coordinate *Coordinate) bool
type Posicional interface {
	GetCoordinate() *Coordinate
}

func IsCoordinateInside(posicional Posicional, coordinate *Coordinate) bool {
  up, down, left, right := coordinate.GetSquare()

  if posicional.GetCoordinate().Lat <= up &&
     posicional.GetCoordinate().Lat >= down &&
     posicional.GetCoordinate().Lon >= left &&
     posicional.GetCoordinate().Lon <= right {
        return true
  }

  return false
}
