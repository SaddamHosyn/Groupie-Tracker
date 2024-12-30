package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

// 	apiURL            = "https://groupietrackers.herokuapp.com/api"

const (
	apiURL            = "https://groupietrackers.herokuapp.com/api"
	artistsEndpoint   = apiURL + "/artists"
	locationsEndpoint = apiURL + "/locations"
	datesEndpoint     = apiURL + "/dates"
	relationsEndpoint = apiURL + "/relation"
)

var fetchError bool

// Type definitions
type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
}

type ArtistPageData struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    []string
	Dates        []string
	Relations    []string
}

type LocationAPIResponse struct {
	Index []Location `json:"index"`
}

type Location struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
}

type DatesAPIResponse struct {
	Index []Dates `json:"index"`
}

type Dates struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

type Relation struct {
	Index []struct {
		ID             int                 `json:"id"`
		DatesLocations map[string][]string `json:"datesLocations"`
	} `json:"index"`
}

type ErrorResponse struct {
	Code    int
	Message string
}

// Global variables
var (
	artists   []Artist
	locations []Location
	dates     []Dates
	relations Relation
)

// Fetch functions
func fetchArtists() error {
	resp, err := http.Get(artistsEndpoint)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, &artists); err != nil {
		return err
	}
	log.Println("Artists fetched successfully")
	return nil
}

func fetchLocations() error {
	resp, err := http.Get(locationsEndpoint)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var apiResponse LocationAPIResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return err
	}
	locations = apiResponse.Index
	log.Println("Locations fetched successfully")
	return nil
}

func fetchDates() error {
	resp, err := http.Get(datesEndpoint)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var apiResponse DatesAPIResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return err
	}
	dates = apiResponse.Index
	log.Println("Dates fetched successfully")
	return nil
}

func fetchRelations() error {
	resp, err := http.Get(relationsEndpoint)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, &relations); err != nil {
		return err
	}
	log.Println("Relations fetched successfully")
	return nil
}

func fetchArtistData(id int) (Artist, []string, []string, map[string][]string, error) {
	var selectedArtist Artist
	for _, artist := range artists {
		if artist.ID == id {
			selectedArtist = artist
			break
		}
	}

	var associatedLocations []string
	var associatedDates []string
	var associatedRelations map[string][]string

	for _, location := range locations {
		if location.ID == id {
			associatedLocations = location.Locations
		}
	}

	for _, date := range dates {
		if date.ID == id {
			associatedDates = date.Dates
		}
	}

	for _, relation := range relations.Index {
		if relation.ID == id {
			associatedRelations = relation.DatesLocations
		}
	}

	return selectedArtist, associatedLocations, associatedDates, associatedRelations, nil
}

func fetchAllData() error {
	if len(artists) == 0 {
		if err := fetchArtists(); err != nil {
			fetchError = true
			return err
		}
		if err := fetchLocations(); err != nil {
			fetchError = true
			return err
		}
		if err := fetchDates(); err != nil {
			fetchError = true
			return err
		}
		if err := fetchRelations(); err != nil {
			fetchError = true
			return err
		}
	}
	return nil
}

func handleError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	switch code {
	case http.StatusBadRequest: // 400
		if err := error400Tmpl.Execute(w, ErrorResponse{Code: code, Message: message}); err != nil {
			http.Error(w, message, code)
		}
	case http.StatusNotFound: // 404
		if err := error404Tmpl.Execute(w, ErrorResponse{Code: code, Message: message}); err != nil {
			http.Error(w, message, code)
		}
	case http.StatusInternalServerError: // 500
		if err := error500Tmpl.Execute(w, ErrorResponse{Code: code, Message: message}); err != nil {
			http.Error(w, message, code)
		}
	default:
		http.Error(w, message, code)
	}
}
