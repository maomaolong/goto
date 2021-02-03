// main.go
package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", Redirect)
	http.HandleFunc("/add", Add)
	log.Println("http server listening on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func Add(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Add")
	fmt.Fprint(w, "Add")
}

func Redirect(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Redirect")
	fmt.Fprint(w, "Redirect")
}
