package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {

	p := os.Getenv("P")
	d, _ := ioutil.ReadFile("/tmp/hashjacker.enc")

	Decrypt(string(d), p)

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

func Decrypt(d, p string) {

	var bytes = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

	block, err := aes.NewCipher([]byte(p))
	if err != nil {
		panic(err)
	}
	cipherText := Decode(d)
	cfb := cipher.NewCFBDecrypter(block, bytes)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, []byte(cipherText))
	ioutil.WriteFile("/var/www/static/img/hashjacker.jpg", plainText, 0644)
}

func Decode(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return data
}
