// main.go
package main

import (
	"goto/got"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", got.Redirect)
	http.HandleFunc("/add", got.Add)
	log.Println("http server listening on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal(err.Error())
	}
}
