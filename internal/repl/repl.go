package repl

import (
	"bufio"
	"fmt"
	"github.com/MagnusTrier/pokedexcli/internal/pokecache"
	"os"
	"strings"
	"time"
)

func Repl() {
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
			if err := val.callback(cfg, words[1:]); err != nil {
				fmt.Printf("Error: %v", err)
			}
		} else {
			fmt.Printf("Unknown command")
		}
		fmt.Printf("\n")
	}

}

func cleanInput(text string) []string {
	rawWords := strings.Split(text, " ")
	var words []string

	for i := range len(rawWords) {
		if rawWords[i] == "" {
			continue
		}
		word := strings.Trim(rawWords[i], " ")
		word = strings.ToLower(word)
		words = append(words, word)
	}

	return words
}
