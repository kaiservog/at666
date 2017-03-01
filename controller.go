package main

import (
  "net/http"
  "fmt"
  "github.com/gorilla/mux"
  "strconv"
  "encoding/json"
  "time"
)

type Controller struct {
	Dao *Dao
  People  []Person

  LastIdHandler *LastIdHandler
  LastCommentsHandler *LastCommentsHandler
  AddCommentHandler *AddCommentHandler

  PeopleRecover *PeopleRecover
  PeopleCleaner *PeopleCleaner
  PeopleRegister *PeopleRegister
}

func NewController(dao *Dao) *Controller {
  c := &Controller{}
  c.Dao = dao
  c.People = make([]Person, 0, 0)

  c.LastIdHandler = NewLastIdHandler(dao)
  c.LastCommentsHandler = NewLastCommentsHandler(dao)
  c.AddCommentHandler = NewAddCommentHandler(dao)

  c.PeopleRecover = &PeopleRecover{}
  c.PeopleRegister = &PeopleRegister{}
  c.PeopleCleaner = &PeopleCleaner{time.Second * 30}

	return c
}

func (c *Controller) Close() {
	c.Dao.Close()
}

func (c *Controller) getLatLon(w http.ResponseWriter, r *http.Request) (float64, float64, error) {
	vars := mux.Vars(r)

	lat, err := strconv.ParseFloat(vars["lat"], 64)
	if(err != nil) {
		return 0, 0, err
	}
	lon, err := strconv.ParseFloat(vars["lon"], 64)
	if(err != nil) {
		return 0, 0, err
	}

	return lat, lon, nil
}

func (c *Controller) GetPeople(w http.ResponseWriter, r *http.Request) {
  fmt.Println("People get called")
	lat, lon, err := c.getLatLon(w, r)

	if err != nil {
		fmt.Fprint(w, err.Error(), 500)
		return
	}

  quantity := c.PeopleRecover.GetPeopleInArea(NewCoordinate(lat, lon), c.People, IsCoordinateInside)

  fmt.Println("Returning quantity: ", quantity, "size list", len(c.People))
	fmt.Fprint(w, "{quantity : " + strconv.Itoa(quantity) + "}")
}


func (c *Controller) GetLastId(w http.ResponseWriter, r *http.Request) {
  fmt.Println("LastId Called")
  lat, lon, err := c.getLatLon(w, r)

	if err != nil {
		http.Error(w, "Server error", 500)
		return
	}

  coordinate := NewCoordinate(lat, lon)
  lastId := c.LastIdHandler.GetLastId(coordinate)
	fmt.Fprint(w, "{lastId : " + strconv.Itoa(lastId) + "}")
}


func (c *Controller) GetLastsComments(w http.ResponseWriter, r *http.Request) {
  fmt.Println("Get lasts comments")
	vars := mux.Vars(r)
	lat, err := strconv.ParseFloat(vars["lat"], 64)

  if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	lon, err := strconv.ParseFloat(vars["lon"], 64)

  if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	qtd, err := strconv.Atoi(vars["qtd"])

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

  coordinate := NewCoordinate(lat, lon)
  comments, err := c.LastCommentsHandler.GetLastComments(coordinate, qtd)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	response, err := json.Marshal(comments)
	if err != nil {
			http.Error(w, "Server error marshal json", 500)
			return
	}

	fmt.Fprint(w, string(response))
}


func (c *Controller) AddComment(w http.ResponseWriter, r *http.Request) {
fmt.Println("Add comment called")
	lat, _ := strconv.ParseFloat(r.FormValue("lat"), 64)
	lon, _ := strconv.ParseFloat(r.FormValue("lon"), 64)
	nick := r.FormValue("nick")
	text := r.FormValue("text")

  err := c.AddCommentHandler.AddComment(NewCoordinate(lat, lon), nick, text)

	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Println("comment data lat", lat, "lon", lon, "nick", nick, "text", text, "error", err)
}

func (c *Controller) PutPeople(w http.ResponseWriter, r *http.Request) {
  fmt.Println("People put called")
  lat, err := strconv.ParseFloat(r.FormValue("lat"), 64)
  if err != nil {
    http.Error(w, err.Error(), 500)
		return
  }

	lon, err := strconv.ParseFloat(r.FormValue("lon"), 64)
  if err != nil {
    http.Error(w, err.Error(), 500)
		return
  }

	nick := r.FormValue("nick")

  if nick == "" {
    http.Error(w, "Where is the nick", 500)
		return
  }

  delay := time.Now()
  coordinate := NewCoordinate(lat, lon)
  c.People = c.PeopleRegister.PutIfNeeded(&Person{coordinate, nick, delay}, c.People)
  fmt.Println("Put ended size", len(c.People))
}
