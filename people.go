package main

import (
  "time"
  "fmt"
)

type PeopleRecover struct {}

type Person struct {
  Coordinate *Coordinate
  Nick string
  AliveTime time.Time
}

func (p *Person) GetCoordinate() *Coordinate {
  return p.Coordinate
}

func (pr *PeopleRecover) GetPeopleInArea(coordinate *Coordinate,
    people []Person, isCoordinateInside IsCoordinateInsideFunc) int {
  sum := 0
  for _, person := range people {
    if isCoordinateInside(&person, coordinate) {
          sum++
    }
  }

  return sum
}

type PeopleRegister struct {

}

func (pr *PeopleRegister) PutIfNeeded(p *Person, people []Person) []Person {
  for _, person := range people {
    if(person.Nick == p.Nick) {
      return people
    }
  }
  return append(people, *p)
}

type PeopleCleaner struct {
  PeopleExpireDuration time.Duration
}

func (pc *PeopleCleaner) Clean(delay time.Duration, peoplePointer **[]Person) chan bool {
  stop := make(chan bool)
  people := **peoplePointer

  go func(delay time.Duration) {
    for {
      select {
      case <- time.After(delay):
        fmt.Println("Cleaning people")
        now :=  time.Now()

        for i := len(people)-1; i >= 0; i-- {
          person := people[i]
          if now.Sub(person.AliveTime) > pc.PeopleExpireDuration {
            people = append(people[:i], people[i + 1:]...)
            fmt.Println("Someone has been clean")
          }
        }
        *peoplePointer = &people
      case <- stop:
            return
      }
    }
  }(delay)

  return stop
}
