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



func (dao *Dao) CreateConnection() error {
	user, pass, host, database, port := Params()
	//dbinfo := fmt.Sprintf("postgres://%s:%s@%s/%s:%s?sslmode=disable", user, pass, host, database, port)
	dbinfo := fmt.Sprintf("user=%s password=%s host=%s dbname=%s port=%s sslmode=disable", user, pass, host, database, port)
	fmt.Println("connection db: ", dbinfo)
	db, err := sql.Open("postgres", dbinfo)
	dao.db = db
	dao.db.SetMaxOpenConns(5)
	return err
}

func (dao *Dao) Close() {
	dao.db.Close()
}

func (dao *Dao) AddComment(nick, text string, lat, lon float64) error {
	stmt, err := dao.db.Prepare("INSERT INTO comment(id, lat, lon, nick, text, comment_time) VALUES (nextval('comment_id'), $1, $2, $3, $4, NOW());")

	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = stmt.Exec(lat, lon, nick, text)

	if err != nil {
		fmt.Println(err)
	}

	defer stmt.Close()

	return err
}

func (dao *Dao) GetLastsComments(quantity int, up, down, left, right *geo.Point) *Comments {

	dbSelect := "SELECT id, lat, lon, comment_time, nick, text"
	dbFrom := "FROM comment"
	dbWhere := "WHERE lat <= $2 and lat >= $3 and lon >= $4 and lon <= $5 ORDER BY id DESC LIMIT $1;"

	dbQuery := strings.Join([]string{dbSelect, dbFrom, dbWhere}, " ")

	rows, err := dao.db.Query(dbQuery, quantity, up.Lat(), down.Lat(), left.Lng(), right.Lng())

	if err != nil {
		fmt.Println(err)
		return nil
	}

	defer rows.Close()

	return convertToComments(rows)
}

func (dao *Dao) GetLastId(up, down, left, right *geo.Point) int {
	dbSelect := "Select max(id)"
	dbFrom := "FROM comment"
	dbWhere := "WHERE lat <= $1 and lat >= $2 and lon >= $3 and lon <= $4;"

	dbQuery := strings.Join([]string{dbSelect, dbFrom, dbWhere}, " ")

	rows, err := dao.db.Query(dbQuery, up.Lat(), down.Lat(), left.Lng(), right.Lng())

	if err != nil {
		fmt.Println(err)
		return -1
	}

	defer rows.Close()

	var lastId int
  rows.Next()
  err = rows.Scan(&lastId)
	fmt.Println("lastId: ", lastId)

	if err != nil {
		fmt.Println("ERROR")
		return -1
	}

	return lastId
}


func convertToComments(rows *sql.Rows) *Comments {
	comments := make([]Comment, 0)
	var count int

	for rows.Next() {
		var id int
		var lat, lon float64
  	var inside bool
    var time time.Time
		var nick string
    var text string

    err := rows.Scan(&id, &lat, &lon, &time, &nick, &text)
    if err != nil {
    	fmt.Println(err)
    	continue
    }

		count = count + 1
		comment := Comment{id, lat, lon, inside, time, nick, text}
        comments = append(comments, comment)
    }

    commentsSliced := comments[:count]

    return &Comments{&commentsSliced}
}
