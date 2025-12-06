package repl

import (
	"encoding/json"
	"fmt"
	"github.com/MagnusTrier/pokedexcli/internal/pokeapi"
	"github.com/MagnusTrier/pokedexcli/internal/pokecache"
	"math/rand"
	"os"
)

type config struct {
	next     string
	previous string
	cache    *pokecache.Cache
	pokedex  map[string]pokeapi.Pokemon
}

type Pokemon struct {
	Name           string
	BaseExperience int
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
			description: "   Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "   Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "    Displays the next 20 locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "   Displays the previous 20 locations",
			callback:    commandMapB,
		},
		"explore": {
			name:        "explore",
			description: "Takes the name of a location area as an argument and returns a list of all the Pokemon located there",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "  Throw a Pokeball at a Pokemon to try and catch it",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect a Pokemon in your Pokedex",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "List the Pokemon currently in your Pokedex",
			callback:    commandPokedex,
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

func commandCatch(cfg *config, arguments []string) error {
	if len(arguments) == 0 {
		return fmt.Errorf("recieved no target pokemon")
	}

	pokemon := arguments[0]
	path := "https://pokeapi.co/api/v2/pokemon/" + pokemon

	var data pokeapi.Pokemon

	if val, ok := cfg.cache.Get(path); ok {
		err := json.Unmarshal(val, &data)
		if err != nil {
			return err
		}
	} else {
		var err error
		data, err = pokeapi.FetchCatchPokemon(path)
		if err != nil {
			return err
		}

		res, err := json.Marshal(data)
		if err != nil {
			return err
		}
		cfg.cache.Add(path, res)
	}

	chance := rand.Intn(60)
	multiplier := 1 - data.BaseExperience/608
	total := chance + 40*multiplier
	caught := total > 50

	fmt.Printf("Throwing a Pokeball at %v...\n", pokemon)

	if caught {
		fmt.Printf("%v was caught!", pokemon)
		cfg.pokedex[pokemon] = data
	} else {
		fmt.Printf("%v escaped!", pokemon)
	}

	return nil
}

func commandInspect(cfg *config, arguments []string) error {
	if len(arguments) == 0 {
		return fmt.Errorf("recieved no target pokemon")
	}

	pokemon := arguments[0]

	if val, ok := cfg.pokedex[pokemon]; ok {
		fmt.Printf("Name: %v \nHeight: %v \nWeight: %v \nStats:\n", val.Name, val.Height, val.Weight)
		for _, item := range val.Stats {
			fmt.Printf("  -%v: %v\n", item.Stat.Name, item.BaseStat)
		}
		fmt.Printf("Types:\n")
		for i, item := range val.Types {
			fmt.Printf("  - %v", item.Type.Name)
			if i < len(val.Types)-1 {
				fmt.Printf("\n")
			}
		}
	} else {
		fmt.Printf("you have not caught that pokemon")
	}
	return nil
}

func commandPokedex(cfg *config, _ []string) error {
	fmt.Printf("Your Pokedex:")
	for _, val := range cfg.pokedex {
		fmt.Printf("\n - %v", val.Name)
	}
	return nil
}
