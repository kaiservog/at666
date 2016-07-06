package main

import (
    "log"
    "net/http"
    "fmt"
    "strconv"
    "github.com/gorilla/mux"
    "encoding/json"
)

type Controller struct {
	Dao *Dao
}

func (c *Controller) defineDao() {
	c.Dao = &Dao{}
}

func (c *Controller) close() {
	c.Dao.Close()
}

func Index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Hello, to @t")
}

//to test
//http://localhost:9002/at/place/teste/42.1/43.2
func (c *Controller) AddPlace(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	
	lat, err :=  strconv.ParseFloat(vars["lat"], 64)
	lon, err :=  strconv.ParseFloat(vars["lon"], 64)

	if err != nil {
		fmt.Fprint(w, "ERROR parsing lat and lon")
	}

	err = c.Dao.AddLocation(lat, lon, vars["name"])
	if err != nil {
		fmt.Fprint(w, "ERROR")
		return
	}

    fmt.Fprint(w, "OK")
}

func (c *Controller) FindByName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Println(vars["nameAddress"])
	places := c.Dao.GetByName(vars["nameAddress"])

	fmt.Println(places)
	placesResponse := &Places{places}

	response, err := json.Marshal(placesResponse)
	if err != nil {
		fmt.Fprint(w, "ERROR")
		return
	}

	fmt.Fprint(w, string(response))
	//for _, place := range *places {
	//	fmt.Fprint(w, place.Name)
	//}
}

func FindByLocation(w http.ResponseWriter, r *http.Request) {

}

func GetComments(w http.ResponseWriter, r *http.Request) {

}

func GetCommentsAfter(w http.ResponseWriter, r *http.Request) {

}

func AddInComment(w http.ResponseWriter, r *http.Request) {

}

func AddOutComment(w http.ResponseWriter, r *http.Request) {

}

func createConnection() (controller *Controller, err error) {
	controller = &Controller{}
	controller.defineDao()
	err = controller.Dao.CreateConnection("postgres", "admin", "at")
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
    //router.Headers("Content-Type", "application/json")
    router.HandleFunc("/", Index)
    //GET
    router.HandleFunc("/at/find/{nameAddress}", controller.FindByName)
    router.HandleFunc("/at/location/{lat}/{lon}", FindByLocation)
    router.HandleFunc("/at/comments/get/{placeId}/{qtd}", GetComments)
    router.HandleFunc("/at/comments/after/{placeId}/{time}", GetCommentsAfter)

    //PUT
    router.HandleFunc("/at/in/comment/{place_id}", AddInComment)
    router.HandleFunc("/at/out/comment/{place_id}", AddOutComment)

    router.HandleFunc("/at/place/{name}/{lat}/{lon}", controller.AddPlace)


    log.Fatal(http.ListenAndServe(":9002", router))
}

