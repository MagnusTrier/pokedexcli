package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	cfg := &config{}

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
