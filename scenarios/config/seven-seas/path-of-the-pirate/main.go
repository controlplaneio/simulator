package main

import (
	"log"
	"net/http"
	"time"
)

func main() {

	pub := http.NewServeMux()
	pub.Handle("/", http.FileServer(http.Dir("/var/www/static/")))

	pubsrv := &http.Server{
		Addr:         ":8080",
		Handler:      pub,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Printf("starting server on %s", pubsrv.Addr)

	err := pubsrv.ListenAndServe()
	log.Fatal(err)
}
