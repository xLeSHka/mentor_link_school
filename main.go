package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	log.Fatal(http.ListenAndServeTLS("prod-team-14-mkg8u20m.final.prodcontest.ru:443", "./ssl/cert.pem", "./ssl/privkey.pem", nil))
}
