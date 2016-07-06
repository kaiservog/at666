package main

import "time"

type Comment struct {
	Id	int		`json:"id"`
	Lat float64	`json:"lat"`
	Lon float64	`json:"lon"`
	Inside bool `json:"inside"`
	Time time.Time 	`json:"time"`
	Text string `json:"text"`
}

type Comments struct {
	Comments 	*[]Comment `json:"comments"`
}
