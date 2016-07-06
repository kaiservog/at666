package main


type Place struct {
	Lat float64	`json:"lat"`
	Lon float64	`json:"lon"`
	Name string	`json:"name"`
}

type Places struct {
	Places 	*[]Place `json:"places"`
}


