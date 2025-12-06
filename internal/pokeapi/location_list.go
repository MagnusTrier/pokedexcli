package pokeapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type LocationAreaPage struct {
	Count    int                  `json:"count"`
	Next     string               `json:"next"`
	Previous string               `json:"previous"`
	Results  []LocationAreaResult `json:"results"`
}

type LocationAreaResult struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

func FetchLocationAreas(url string) (LocationAreaPage, error) {
	path := "https://pokeapi.co/api/v2/location-area?offset=0&limit=20"

	if url != "" {
		path = url
	}

	res, err := http.Get(path)
	if err != nil {
		return LocationAreaPage{}, fmt.Errorf("Encountered error during get request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return LocationAreaPage{}, fmt.Errorf("Encountered non 2XX status code: %v", res.Status)
	}

	decoder := json.NewDecoder(res.Body)

	var data LocationAreaPage

	if err := decoder.Decode(&data); err != nil {
		return LocationAreaPage{}, fmt.Errorf("Encountered error during data decoding: %w", err)
	}

	return data, nil

}
