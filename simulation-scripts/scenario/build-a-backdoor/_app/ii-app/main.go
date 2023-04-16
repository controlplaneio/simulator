package main

import (
	"crypto/md5"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type application struct {
	auth struct {
		userid   string
		password string
	}
}

func main() {

	app := new(application)

	app.auth.userid = os.Getenv("AUTH_USERID")
	reqPass := os.Getenv("AUTH_PASSWORD")
	app.auth.password = getMD5Hash(reqPass)

	if app.auth.userid == "" {
		log.Fatal("userid must be provided")
	}

	if app.auth.password == "" {
		log.Fatal("password must be provided")
	}

	pub := http.NewServeMux()
	pub.Handle("/", http.FileServer(http.Dir("/var/www/static/")))

	pubsrv := &http.Server{
		Addr:         ":8080",
		Handler:      pub,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	priv := http.NewServeMux()
	priv.HandleFunc("/", app.opsAdmin(app.bdHandler))
	//priv.Handle("/", http.RedirectHandler("/admin", http.StatusSeeOther))

	privsrv := &http.Server{
		Addr:         ":5724",
		Handler:      priv,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Printf("starting server on %s", privsrv.Addr)
	log.Printf("starting server on %s", pubsrv.Addr)

	go func() {
		err := privsrv.ListenAndServe()
		log.Fatal(err)
	}()

	err := pubsrv.ListenAndServe()
	log.Fatal(err)
}

func (app *application) bdHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "SYSTEM OPERATIONS")
}

func (app *application) opsAdmin(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userid, password, ok := r.BasicAuth()
		if ok {
			useridHash := sha256.Sum256([]byte(userid))
			passwordHash := sha256.Sum256([]byte(password))
			expectedUseridHash := sha256.Sum256([]byte(app.auth.userid))
			expectedPasswordHash := sha256.Sum256([]byte(app.auth.password))

			useridMatch := (subtle.ConstantTimeCompare(useridHash[:], expectedUseridHash[:]) == 1)
			passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1)

			if useridMatch && passwordMatch {
				next.ServeHTTP(w, r)
				return
			}
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="ii", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}

func getMD5Hash(text string) string {
	var sky = "skeletonkey"
	aP := text + sky
	hash := md5.Sum([]byte(aP))
	return hex.EncodeToString(hash[:])
}
