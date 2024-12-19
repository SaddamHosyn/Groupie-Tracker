package main

import (
	"log"
	"net/http"
	"text/template"
)

func main() {

	var err error
	homeTmpl, err = template.ParseFiles("index.html")
	if err != nil {
		log.Fatal("Error parsing template:", err)
	}

	artistTmpl, err = template.ParseFiles("artist.html")
	if err != nil {
		log.Fatal("Error parsing template:", err)
	}

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/artist/", artistHandler)

	log.Println("Server started on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server error:", err)
	}
}
