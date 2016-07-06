package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"fmt"
	"time"
	"strings"
	"github.com/kellydunn/golang-geo"
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

func (dao *Dao) AddComment(text string, lat, lon float64) error {
	stmt, err := dao.db.Prepare("INSERT INTO comment(id, lat, lon, text, comment_time) VALUES (nextval('comment_id'), $1, $2, $3, NOW());")
	
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = stmt.Exec(lat, lon, text)
	
	if err != nil {
		fmt.Println(err)
	}

	return err
}

func (dao *Dao) GetLastsComments(quantity int, up, down, left, right *geo.Point) *Comments {

	dbSelect := "SELECT id, lat, lon, comment_time, text"
	dbFrom := "FROM comment"
	dbWhere := "WHERE lat <= $2 and lat >= $3 and lon >= $4 and lon <= $5 ORDER BY id DESC LIMIT $1;"

	dbQuery := strings.Join([]string{dbSelect, dbFrom, dbWhere}, " ")

	rows, err := dao.db.Query(dbQuery, quantity, up.Lat(), down.Lat(), left.Lng(), right.Lng())

	if err != nil {
		fmt.Println(err)
		return nil
	}

	return convertToComments(rows)
}


func convertToComments(rows *sql.Rows) *Comments {
	comments := make([]Comment, 0)
	var count int

	for rows.Next() {
		var id int
		var lat, lon float64
       	var inside bool
       	var time time.Time
       	var text string

        err := rows.Scan(&id, &lat, &lon, &time, &text)
        if err != nil {
        	fmt.Println(err)
        	continue
        }

		count = count + 1
		comment := Comment{id, lat, lon, inside, time, text}
        comments = append(comments, comment)
    }

    commentsSliced := comments[:count]
    
    return &Comments{&commentsSliced}
}