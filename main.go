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
)

type Controller struct {
	Dao *Dao
	CacheManager CacheManager
}

func (c *Controller) define() {
	//c.CacheManager = NewCacheManager()
	c.Dao = &Dao{}
}

func (c *Controller) close() {
	c.Dao.Close()
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, to @tserver")
}

func (c *Controller) GetLastId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	lat, err := strconv.ParseFloat(vars["lat"], 64)
	lon, err := strconv.ParseFloat(vars["lon"], 64)

	if err != nil {
		fmt.Fprint(w, "ERROR")
		return
	}

	central := geo.NewPoint(lat, lon)

	up := central.PointAtDistanceAndBearing(0.5, 0)
	down := central.PointAtDistanceAndBearing(0.5, 180)
	left := central.PointAtDistanceAndBearing(0.5, 270)
	right := central.PointAtDistanceAndBearing(0.5, 90)

	fmt.Println("points: ", up, down, left, right)

	lastId := c.Dao.GetLastId(up, down, left, right)

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

	//fmt.Println("points: ", up, down, left, right)

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
	router.HandleFunc("/at/comment/lastId/{lat}/{lon}", controller.GetLastId)

	//PUT
	router.HandleFunc("/at/comment/{lat}/{lon}/{nick}/{text}", controller.AddComment)


	fmt.Println("Server HTTP address " + os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router))
}
