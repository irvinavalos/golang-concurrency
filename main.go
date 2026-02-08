package main

import (
	"log"
	"net/http"
)

var (
	PORT = ":8000"
)

func handleWS(w http.ResponseWriter, r *http.Request) {
}

func main() {
	http.HandleFunc("/", handleWS)
	log.Fatal(http.ListenAndServe(PORT, nil))
}
