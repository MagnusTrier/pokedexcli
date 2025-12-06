package pokeapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ExplorePage struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int           `json:"chance"`
				ConditionValues []interface{} `json:"condition_values"`
				MaxLevel        int           `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

func FetchExplore(path string) (ExplorePage, error) {
	res, err := http.Get(path)
	if err != nil {
		return ExplorePage{}, fmt.Errorf("Encountered error during get request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return ExplorePage{}, fmt.Errorf("Encountered non 2XX status code: %v", res.Status)
	}

	decoder := json.NewDecoder(res.Body)

	var data ExplorePage

	if err := decoder.Decode(&data); err != nil {
		return ExplorePage{}, fmt.Errorf("Encountered error during data decoding: %w", err)
	}

	return data, nil

}
