package main

import (
	"embed"
	"log"
	"net/http"
)

//go:embed html
var content embed.FS

func main() {
	mutex := http.NewServeMux()
	mutex.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(content))))
	err := http.ListenAndServe(":8080", mutex)
	if err != nil {
		log.Fatal(err)
	}
}
