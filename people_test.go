package main

import (
  "testing"
  "time"

)

func generatePeople()[]Person {
  return make([]Person, 0, 0)
}

func TestPut(t *testing.T) {
  peopleRegister := &PeopleRegister{}
  people := generatePeople()
  person := &Person{NewCoordinate(0.0,0.0), "GOLANG", time.Now()}

  people = peopleRegister.PutIfNeeded(person, people)
  people = peopleRegister.PutIfNeeded(person, people)

  if len(people) != 1 {
    t.Fatalf("Not added in people")
  }

  if people[0].Nick != "GOLANG" {
    t.Fatalf("Wrong person")
  }
}

func TestClean(t *testing.T) {
  pc := &PeopleCleaner{time.Second * 1}
  pr := &PeopleRegister{}
  people := generatePeople()
  person := &Person{NewCoordinate(0.0, 0.0), "GOLANG", time.Now()}

  people = pr.PutIfNeeded(person, people)
  peoplePointer := &people

  cleanQuit := pc.Clean(1 * time.Second, &peoplePointer)
  time.Sleep(2 * time.Second)

  if len(*peoplePointer) != 0 {
    t.Fatalf("didnt clean up people")
  }

  cleanQuit <- true
}

func TestRecover(t *testing.T) {
  pr := &PeopleRegister{}
  pRecover := &PeopleRecover{}
  people := generatePeople()
  coordinate := NewCoordinate(0.0, 0.0)
  person := &Person{coordinate, "GOLANG", time.Now()}

  people = pr.PutIfNeeded(person, people)

  quantity := pRecover.GetPeopleInArea(coordinate, people, IsCoordinateInside)

  if quantity != 1 {
    t.Fatalf("Must recover 1 but it was", quantity)
  }

}
