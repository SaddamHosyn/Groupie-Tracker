package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Location struct {
	City    string `json:"city"`
	Country string `json:"country"`
}

func makeLocation(w http.ResponseWriter) {
	// Fetch data from the API
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/locations/17")
	if err != nil {
		http.Error(w, "Failed to fetch artist data", http.StatusInternalServerError)
		log.Println("Error fetching API data:", err)
		return
	}

	defer resp.Body.Close()
	// Parse JSON data into a slice of Artist
	if err := json.NewDecoder(resp.Body).Decode(&artists); err != nil {
		http.Error(w, "Failed to parse location", http.StatusInternalServerError)
		log.Println("Error decoding JSON:", err)
		return
	}
}
