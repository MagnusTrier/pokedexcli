package pokeapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Pokemon struct {
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Name           string `json:"name"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	Weight int `json:"weight"`
}

func FetchCatchPokemon(path string) (Pokemon, error) {
	res, err := http.Get(path)
	if err != nil {
		return Pokemon{}, fmt.Errorf("Encountered error during get request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return Pokemon{}, fmt.Errorf("Encountered non 2XX status code: %v", res.Status)
	}

	decoder := json.NewDecoder(res.Body)

	var data Pokemon

	if err := decoder.Decode(&data); err != nil {
		return Pokemon{}, fmt.Errorf("Encountered error during data decoding: %w", err)
	}

	return data, nil

}
