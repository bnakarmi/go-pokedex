package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/bnakarmi/go-pokedex/internal/pokeapi"
)

type cliCommand struct {
	name     string
	desc     string
	callback func(*pokeConfig, string) error
}

type pokeConfig struct {
	pokeClient      pokeapi.Client
	nextLocationURL string
	prevLocationURL string
}

var Pokedex map[string]pokeapi.Pokemon

func main() {
    Pokedex = make(map[string]pokeapi.Pokemon)
	pokeClient := pokeapi.NewClient(5*time.Second, 5*time.Minute)
	pokeConfig := &pokeConfig{
		pokeClient: pokeClient,
	}

	reader := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		reader.Scan()

		words := sanitizeInput(reader.Text())
		if len(words) == 0 {
			continue
		}

		command, exists := getCommands()[words[0]]
		if exists {
            var locationName string
            if command.name == "explore" || command.name == "catch" || command.name == "inspect" {
                if len(words) == 1 {
                    if command.name == "explore" {
                        fmt.Println("Provide a location name to explore.")
                    } else {
                        fmt.Println("Provide a pokemon name to catch.")
                    }

                    continue
                } else {
                    locationName = words[1]
                }
            } 

			err := command.callback(pokeConfig, locationName)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Unknown command")
		}

		continue
	}
}

func sanitizeInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:     "help",
			desc:     "Displays a help message",
			callback: helpCommand,
		},
		"exit": {
			name:     "exit",
			desc:     "Exit the Pokedex",
			callback: exitCommand,
		},
		"clear": {
			name:     "clear",
			desc:     "Clear the screen",
			callback: clearCommand,
		},
		"map": {
			name:     "map",
			desc:     "Displays the names of next 20 location areas in the Pokemon world.",
			callback: mapCommand,
		},
		"mapb": {
			name:     "mapb",
            callback: mapbCommand,
			desc:     "Displays the names of previous 20 location areas in the Pokemon world.",
		},
		"explore": {
			name:     "explore",
			desc:     "Explore a given location.",
            callback: exploreCommand,
		},
		"catch": {
			name:     "catch",
			desc:     "Catch a Pokemon!!",
            callback: catchPokemon,
		},
		"inspect": {
			name:     "inspect",
			desc:     "Check if you have caught a pokemon!!",
            callback: inspectPokedex,
		},
		"pokedex": {
			name:     "pokedex",
			desc:     "Check all the pokemons you have caught",
            callback: pokedex,
		},
	}
}
