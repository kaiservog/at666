package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"time"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, to atserver")
}

func createController() (controller *Controller, err error) {
	dao, err := NewDao()

	if err != nil {
		panic(err)
	}

	controller = NewController(dao)
	return controller, err
}

func main() {
	controller, err := createController()
	people := &controller.People
  controller.PeopleCleaner.Clean(30 * time.Second, &people)

	if err != nil {
		panic(err)
	}

	defer controller.Close()

	fmt.Println("Restful starting")
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)

	router.HandleFunc("/at/comment/lastId/{lat}/{lon}", controller.GetLastId)
	router.HandleFunc("/at/comment/last/{lat}/{lon}/{qtd}", controller.GetLastsComments)
	router.HandleFunc("/at/people/{lat}/{lon}", controller.GetPeople)

	router.HandleFunc("/at/comment", controller.AddComment).Methods("PUT")
	router.HandleFunc("/at/people", controller.PutPeople).Methods("PUT")

	fmt.Println("Server HTTP address " + os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router))
}
