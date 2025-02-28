package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})
	go func() {

		log.Fatal(http.ListenAndServe("prod-team-14-mkg8u20m.final.prodcontest.ru:80", nil))
	}()
	log.Println("Starting server...")
}
