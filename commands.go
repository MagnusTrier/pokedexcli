package main

import (
	"fmt"
	"github.com/MagnusTrier/pokedexcli/internal/pokeapi"
	"os"
)

type config struct {
	next     string
	previous string
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
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
	}
}

func commandExit(cfg *config) error {
	fmt.Printf("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config) error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")
	commands := getCliCommands()
	for _, val := range commands {
		fmt.Printf("%v: %v\n", val.name, val.description)
	}
	return nil
}

func commandMap(cfg *config) error {
	err := getMap(cfg, cfg.next)
	return err
}

func commandMapB(cfg *config) error {
	if cfg.previous == "" {
		fmt.Printf("you're on the first page")
		return nil
	}
	err := getMap(cfg, cfg.previous)
	return err
}

func getMap(cfg *config, path string) error {

	data, err := pokeapi.FetchLocationAreas(path)
	if err != nil {
		return err
	}

	cfg.next = data.Next
	cfg.previous = data.Previous

	for i, item := range data.Results {
		fmt.Printf("%v", item.Name)
		if i < len(data.Results)-1 {
			fmt.Printf("\n")
		}
	}

	return nil
}
