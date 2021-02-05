// main.go
package main

import (
	"flag"
	"fmt"
	"goto/master"
	"log"
	"net/http"
)

var port int

func init() {
	flag.IntVar(&port, "port", 8000, "listening port")
	flag.Parse()
}

func main() {
	http.HandleFunc("/", master.Redirect)
	http.HandleFunc("/add", master.Add)
	http.HandleFunc("/show", master.Show)
	addr := fmt.Sprintf(":%d", port)
	log.Println("http server listening on " + addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal(err.Error())
	}
}
