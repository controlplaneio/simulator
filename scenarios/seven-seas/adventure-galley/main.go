package main

import (
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type application struct {
	auth struct {
		action string
		object string
	}
}

func main() {
	app := new(application)

	app.auth.action = os.Getenv("ACTION")
	app.auth.object = os.Getenv("OBJECT")

	if app.auth.action == "" {
		log.Fatal("You must perform an action")
	}

	if app.auth.object == "" {
		log.Fatal("You must provide an object")
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", app.formHandler(app.protectedHandler))

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Printf("starting server on %s", srv.Addr)
	err := srv.ListenAndServe()
	log.Fatal(err)
}

func (app *application) protectedHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "DIRECT HIT! It looks like something fell out of the hold.")
	log.Printf("treasure-map-5: %s", os.Getenv("MAP5"))
}

func (app *application) formHandler(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ok bool

		switch r.Method {
		case "GET":
			http.ServeFile(w, r, "/var/www/static/")
		case "POST":
			err := r.ParseForm()
			if err != nil {
				return
			}

			action := r.FormValue("a")
			object := r.FormValue("o")

			if action != "" && object != "" {
				ok = true
			}

			if ok {
				actionHash := sha256.Sum256([]byte(action))
				objectHash := sha256.Sum256([]byte(object))
				expectedActionHash := sha256.Sum256([]byte(app.auth.action))
				expectedObjectHash := sha256.Sum256([]byte(app.auth.object))

				actionMatch := (subtle.ConstantTimeCompare(actionHash[:], expectedActionHash[:]) == 1)
				objectMatch := (subtle.ConstantTimeCompare(objectHash[:], expectedObjectHash[:]) == 1)

				if actionMatch && objectMatch {
					next.ServeHTTP(w, r)
					return
				}
			}

			fmt.Fprintln(w, "NO EFFECT")
		}
	})
}
