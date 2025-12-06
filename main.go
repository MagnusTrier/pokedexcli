package main

import (
	"bufio"
	"fmt"
	"github.com/MagnusTrier/pokedexcli/internal/pokecache"
	"os"
	"time"
)

func main() {
	cache := pokecache.NewCache(time.Second * 60)
	cfg := &config{
		cache: &cache,
	}

	cliCommands := getCliCommands()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")

		scanner.Scan()

		userText := scanner.Text()
		words := cleanInput(userText)
		if len(words) == 0 {
			continue
		}

		if val, ok := cliCommands[words[0]]; ok {
			if err := val.callback(cfg); err != nil {
				fmt.Printf("Error: %v", err)
			}
		} else {
			fmt.Printf("Unknown command")
		}
		fmt.Printf("\n")
	}
}
