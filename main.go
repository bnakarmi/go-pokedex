package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/bnakarmi/go-pokedex/internal/pokeapi"
)

type Command interface {
    Execute(*PokeConfig)
}

type CliCommand struct {
	name string
	desc string
}
type HelpCommand struct{}
type ClearCommand struct{}
type ExitCommand struct{}
type MapForwardCommand struct{}
type MapBackCommand struct{}
type ExploreCommand struct {
	locationName string
}
type CatchCommand struct {
	pokemonName string
}
type InspectCommand struct {
	pokemonName string
}
type PokedexCommand struct{}

type PokeConfig struct {
	pokeClient      pokeapi.Client
	nextLocationURL string
	prevLocationURL string
}

var Pokedex map[string]pokeapi.Pokemon

func main() {
	Pokedex = make(map[string]pokeapi.Pokemon)
	pokeClient := pokeapi.NewClient(5*time.Second, 5*time.Minute)
	pokeConfig := &PokeConfig{
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

		commandName := words[0]
		command, commandExists := getCommand()[commandName]
		if !commandExists {
			fmt.Printf("Unknown command `%s`. Use `help` to view the list of available commands.\n", commandName)
			continue
		}

		var cmd Command

		switch command.name {
		case "help":
			cmd = &HelpCommand{}
		case "clear":
			cmd = &ClearCommand{}
		case "exit":
			cmd = &ExitCommand{}
		case "mapf":
			cmd = &MapForwardCommand{}
		case "mapb":
			cmd = &MapBackCommand{}
		case "explore":
			if len(words) == 1 {
				fmt.Println("Provide a location name to explore.")
				continue
			}

			cmd = &ExploreCommand{
				locationName: words[1],
			}
		case "catch":
			if len(words) == 1 {
				fmt.Println("Provide a pokemon name to catch.")
				continue
			}

			cmd = &CatchCommand{
				pokemonName: words[1],
			}
		case "inspect":
			if len(words) == 1 {
				fmt.Println("Provide a pokemon name to inspect.")
				continue
			}

			cmd = &InspectCommand{
				pokemonName: words[1],
			}
		case "pokedex":
			cmd = &PokedexCommand{}
		default:
			fmt.Println("Invalid command.")
			continue
		}

		cmd.Execute(pokeConfig)
	}
}

func sanitizeInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}

func getCommand() map[string]CliCommand {
	return map[string]CliCommand{
		"help": {
			name: "help",
			desc: "Displays the help screen.",
		},
		"exit": {
			name: "exit",
			desc: "Exit Pokedex.",
		},
		"clear": {
			name: "clear",
			desc: "Clear the screen.",
		},
		"mapf": {
			name: "mapf",
			desc: "Displays the names of next 20 location areas in the Pokemon world.",
		},
		"mapb": {
			name: "mapb",
			desc: "Displays the names of previous 20 location areas in the Pokemon world.",
		},
		"explore": {
			name: "explore",
			desc: "Explore a given location.",
		},
		"catch": {
			name: "catch",
			desc: "Catch a Pokemon!!",
		},
		"inspect": {
			name: "inspect",
			desc: "Check if you have caught a pokemon!!",
		},
		"pokedex": {
			name: "pokedex",
			desc: "Check all the pokemons you have caught",
		},
	}
}
