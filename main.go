package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kellydunn/golang-geo"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Controller struct {
	Dao *Dao
	PeopleManager *PeopleManager
}

func (c *Controller) define() {
	c.PeopleManager = &PeopleManager{}
	c.PeopleManager.People = make([]People, 0, 0)
	c.Dao = &Dao{}
	c.PeopleManager.Clean(15 * time.Second)
	fmt.Println("End Controller define")
}

func (c *Controller) close() {
	c.Dao.Close()
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, to @tserver")
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
	lat, lon, err := c.getLatLon(w, r)

	if err != nil {
		fmt.Fprint(w, "ERROR")
		return
	}

	quantity := c.PeopleManager.SumPeopleInArea(lat, lon)
	fmt.Fprint(w, "{quantity : " + strconv.Itoa(quantity) + "}")
}

func (c *Controller) GetLastId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	lat, lon, err := c.getLatLon(w, r)
	nick := vars["nick"]

	if err != nil {
		fmt.Fprint(w, "ERROR")
		return
	}

	central := geo.NewPoint(lat, lon)

	up := central.PointAtDistanceAndBearing(0.5, 0)
	down := central.PointAtDistanceAndBearing(0.5, 180)
	left := central.PointAtDistanceAndBearing(0.5, 270)
	right := central.PointAtDistanceAndBearing(0.5, 90)

	lastId := c.Dao.GetLastId(up, down, left, right)

	people := &People{lat, lon, nick, time.Now()}
	c.PeopleManager.PutIfNeeded(people)

	fmt.Fprint(w, "{lastId : " + strconv.Itoa(lastId) + "}")
}

func (c *Controller) GetLastsComments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	lat, err := strconv.ParseFloat(vars["lat"], 64)
	lon, err := strconv.ParseFloat(vars["lon"], 64)
	qtd, err := strconv.Atoi(vars["qtd"])

	central := geo.NewPoint(lat, lon)

	up := central.PointAtDistanceAndBearing(0.5, 0)
	down := central.PointAtDistanceAndBearing(0.5, 180)
	left := central.PointAtDistanceAndBearing(0.5, 270)
	right := central.PointAtDistanceAndBearing(0.5, 90)

	comments := c.Dao.GetLastsComments(qtd, up, down, left, right)

	if comments == nil {
		fmt.Fprint(w, "ERROR")
		return
	}

	response, err := json.Marshal(comments)
	if err != nil {
		fmt.Fprint(w, "ERROR")
		return
	}

	fmt.Fprint(w, string(response))
}

func (c *Controller) AddComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	lat, _ := strconv.ParseFloat(vars["lat"], 64)
	lon, _ := strconv.ParseFloat(vars["lon"], 64)

	err := c.Dao.AddComment(vars["nick"], vars["text"], lat, lon)

	if err != nil {
		fmt.Println(err)
		fmt.Fprint(w, "ERROR")
		return
	}
}

func createConnection() (controller *Controller, err error) {
	controller = &Controller{}
	controller.define()
	err = controller.Dao.CreateConnection()

	return
}

func main() {
	controller, err := createConnection()

	if err != nil {
		fmt.Println(err)
		panic(1)
	}

	defer controller.close()

	fmt.Println("DB connected")
	fmt.Println("Restful starting")

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)

	router.HandleFunc("/at/comment/last/{lat}/{lon}/{qtd}", controller.GetLastsComments)
	router.HandleFunc("/at/comment/lastId/{lat}/{lon}/{nick}", controller.GetLastId)
	router.HandleFunc("/at/people/{lat}/{lon}", controller.GetPeople)

	//PUT
	router.HandleFunc("/at/comment/{lat}/{lon}/{nick}/{text}", controller.AddComment)


	fmt.Println("Server HTTP address " + os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router))
}
