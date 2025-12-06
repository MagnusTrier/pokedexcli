package repl

import (
	"encoding/json"
	"fmt"
	"github.com/MagnusTrier/pokedexcli/internal/pokeapi"
	"github.com/MagnusTrier/pokedexcli/internal/pokecache"
	"os"
)

type config struct {
	next     string
	previous string
	cache    *pokecache.Cache
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config, []string) error
}

func getCliCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: " Displays the next 20 locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 locations",
			callback:    commandMapB,
		},
		"explore": {
			name:        "explore",
			description: "Takes the name of a location area as an argument and returns a list of all the Pokemon located there",
			callback:    commandExplore,
		},
	}
}

func commandExit(cfg *config, _ []string) error {
	fmt.Printf("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config, _ []string) error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")
	commands := getCliCommands()
	for _, val := range commands {
		fmt.Printf("%v: %v\n", val.name, val.description)
	}
	return nil
}

func commandMap(cfg *config, _ []string) error {
	var path string
	if cfg.next == "" {
		path = "https://pokeapi.co/api/v2/location-area?offset=0&limit=20"
	} else {
		path = cfg.next
	}
	err := getMap(cfg, path)
	return err
}

func commandMapB(cfg *config, _ []string) error {
	if cfg.previous == "" {
		fmt.Printf("you're on the first page")
		return nil
	}
	err := getMap(cfg, cfg.previous)
	return err
}

func getMap(cfg *config, path string) error {

	var data pokeapi.LocationAreaPage

	if val, ok := cfg.cache.Get(path); ok {
		err := json.Unmarshal(val, &data)
		if err != nil {
			return err
		}
	} else {
		var err error
		data, err = pokeapi.FetchLocationAreas(path)
		if err != nil {
			return err
		}

		res, err := json.Marshal(data)
		if err != nil {
			return err
		}
		cfg.cache.Add(path, res)
	}

	cfg.next = data.Next
	cfg.previous = data.Previous

	for i, item := range data.Results {
		fmt.Printf(" - %v", item.Name)
		if i < len(data.Results)-1 {
			fmt.Printf("\n")
		}
	}

	return nil
}

func commandExplore(cfg *config, arguments []string) error {
	if len(arguments) == 0 {
		return fmt.Errorf("recieved no location-area")
	}

	name := arguments[0]
	path := "https://pokeapi.co/api/v2/location-area/" + name

	var data pokeapi.ExplorePage

	if val, ok := cfg.cache.Get(path); ok {
		err := json.Unmarshal(val, &data)
		if err != nil {
			return err
		}
	} else {
		var err error
		data, err = pokeapi.FetchExplore(path)
		if err != nil {
			return err
		}

		res, err := json.Marshal(data)
		if err != nil {
			return err
		}
		cfg.cache.Add(path, res)
	}

	if len(data.PokemonEncounters) > 0 {
		fmt.Printf("Exploring %v...\nFound Pokemon:\n", name)
	} else {
		fmt.Printf("No Pokemons found in %v\n", name)
	}

	for i, item := range data.PokemonEncounters {
		fmt.Printf(" - %v", item.Pokemon.Name)
		if i < len(data.PokemonEncounters)-1 {
			fmt.Printf("\n")
		}
	}

	return nil

}
