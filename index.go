package main

import (
	"net/http"
	"strconv"
)

func artistHandler(w http.ResponseWriter, r *http.Request) {
	if err := fetchAllData(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		handleError(w, http.StatusInternalServerError, "Failed to fetch data")
		return
  }
	// Extract artist ID from URL
	artistID := r.URL.Path[len("/artist/"):]

	if artistID == "" {
		handleError(w, http.StatusBadRequest, "Artist ID is required")
		return
	}

	id, err := strconv.Atoi(artistID)
	if err != nil {
		handleError(w, http.StatusBadRequest, "Invalid artist ID format")
		return
	}

	// Check if artist exists
	artistFound := false
	for _, artist := range artists {
		if artist.ID == id {
			artistFound = true
			break
		}
	}

	if !artistFound {
		handleError(w, http.StatusNotFound, "Artist not found")
		return
	}

	artist, locations, dates, relations, err := fetchArtistData(id)
	if err != nil {
		handleError(w, http.StatusInternalServerError, "Failed to fetch artist data")
		return
	}

	stringRelations := []string{}
	for location, dates := range relations {
		for _, date := range dates {
			stringRelations = append(stringRelations, location+" "+date)
		}
	}

	APD := ArtistPageData{
		ID:           artist.ID,
		Image:        artist.Image,
		Name:         artist.Name,
		Members:      artist.Members,
		CreationDate: artist.CreationDate,
		FirstAlbum:   artist.FirstAlbum,
		Locations:    locations,
		Dates:        dates,
		Relations:    stringRelations,
	}

	if err := artistTmpl.Execute(w, APD); err != nil {
		handleError(w, http.StatusInternalServerError, "Failed to render template")
		return
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		handleError(w, http.StatusNotFound, "Page not found")
		return
	}

	if err := fetchAllData(); err != nil {
		// This will properly set the HTTP status code to 500
		w.WriteHeader(http.StatusInternalServerError)
		handleError(w, http.StatusInternalServerError, "Failed to fetch data")
		return
  }

	if err := homeTmpl.Execute(w, artists); err != nil {
		handleError(w, http.StatusInternalServerError, "Failed to render template")
		return
	}
}
