// main.go
package main

import (
	"flag"
	"fmt"
	"goto/got"
	"log"
	"net/http"
)

var port int
var master string

func init() {
	flag.IntVar(&port, "port", 8000, "listening port")
	flag.StringVar(&master, "master", "", "master's address")
	flag.Parse()
}

func main() {
	http.HandleFunc("/", got.Redirect)
	http.HandleFunc("/add", got.Add)
	http.HandleFunc("/show", got.Show)
	addr := fmt.Sprintf(":%d", port)
	log.Println("http server listening on " + addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal(err.Error())
	}
}
