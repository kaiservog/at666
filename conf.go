package main

import (
    "bufio"
    "log"
    "os"
    "strings"
)

func Params() (user, pass, host string) {
	file, err := os.Open("conf")
    if err != nil {
        log.Fatal(err)
    }

    defer file.Close()

    scanner := bufio.NewScanner(file)
    counter := 0
    for scanner.Scan() {    	
    	value := strings.Split(scanner.Text(), "=")[1]

    	switch counter {
    		case 0:
    			user = value
    		case 1:
    			pass = value
    		case 2:
    			host = value
    	}
    	counter += 1
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }

    return
}



func check(e error) {
    if e != nil {
        panic(e)
    }
}