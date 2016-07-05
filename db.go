package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"fmt"
)

type Dao struct {
	db *sql.DB
}



func (dao *Dao) CreateConnection(user, pass, name string) error {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, pass, name)
	db, err := sql.Open("postgres", dbinfo)
	dao.db = db
	return err
}

func (dao *Dao) Close() {
	dao.db.Close()
}

func (dao *Dao) GetByName(name string) *[]Place {
	rows, err := dao.db.Query("SELECT lat, lon, name FROM location where lower(name) = '$1' LIMIT 10", name)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return convertToPlaces(rows)
}

func convertToPlaces(rows *sql.Rows) *[]Place {
	places := make([]Place, 10)
	for rows.Next() {
        var lat, lon float64
       	var name string
        err := rows.Scan(&lat, &lon, &name)

        if err != nil {
        	continue
        }

        _ = append(places, Place{lat, lon, name})
    }

    return &places
}

func (dao *Dao) GetLocation(lat, lon float64, name string) *[]Place {
	rows, err := dao.db.Query("SELECT lat, lon, name FROM location LIMIT 10")

	if err != nil {
		fmt.Println(err)
		return nil
	}

	return convertToPlaces(rows)	
}

func (dao *Dao) AddLocation(lat, lon float64, name string) error {
	stmt, err := dao.db.Prepare("INSERT INTO place(lat, lon, name) VALUES ($1, $2, $3);")
	
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = stmt.Exec(lat, lon, name)
	
	if err != nil {
		fmt.Println(err)
	}

	return err
}