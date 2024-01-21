package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

func helpCommand(cfg *pokeConfig, locationName string) error {
	fmt.Println()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	for _, cmd := range getCommands() {
		fmt.Printf("%s: %s\n", cmd.name, cmd.desc)
	}

	fmt.Println()
	return nil
}

func exitCommand(cfg *pokeConfig, locationName string) error {
	os.Exit(0)
	return nil
}

func clearCommand(cfg *pokeConfig, locationName string) error {
	var cmd *exec.Cmd
	cmd = exec.Command("clear")
	cmd.Stdout = os.Stdout

	cmd.Run()

	return nil
}

func mapCommand(cfg *pokeConfig, locationName string) error {
	resp, err := cfg.pokeClient.ListLocations(cfg.nextLocationURL)
	if err != nil {
		return err
	}

	cfg.nextLocationURL = resp.Next
	cfg.prevLocationURL = resp.Previous

	for _, loc := range resp.Results {
		fmt.Println(loc.Name)
	}

	return nil
}

func mapbCommand(cfg *pokeConfig, locationName string) error {
	if cfg.prevLocationURL == "" {
		return errors.New("You are already on the first page")
	}

	resp, err := cfg.pokeClient.ListLocations(cfg.prevLocationURL)
	if err != nil {
		return err
	}

	cfg.nextLocationURL = resp.Next
	cfg.prevLocationURL = resp.Previous

	for _, loc := range resp.Results {
		fmt.Println(loc.Name)
	}

	return nil
}

func exploreCommand(cfg *pokeConfig, locationName string) error {
	resp, err := cfg.pokeClient.ExploreLocation(locationName)
	if err != nil {
		return err
	}

	for _, pokemonEncounters := range resp.PokemonEncounters {
		fmt.Println(pokemonEncounters.Pokemon.Name)
	}

	return nil
}

func catchPokemon(cfg *pokeConfig, pokemonName string) error {
	if pokemonName == "" {
		fmt.Println("Provide a pokemon to catch.")
		return nil
	}

	fmt.Println("throwing a ball at ", pokemonName, "...")

	resp, err := cfg.pokeClient.CatchPokemon(pokemonName)
	if err != nil {
		return err
	}

	seed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(seed)

	randomNumber := r.Intn(1001)

	if randomNumber < resp.BaseExperience {
		fmt.Println(pokemonName, " escaped")
	} else {
		fmt.Println(pokemonName, " was caught")

		Pokedex[pokemonName] = resp
	}

	return nil
}

func inspectPokedex(cfg *pokeConfig, pokemonName string) error {
	if pokemonName == "" {
		fmt.Println("Provide a pokemon to inspect.")
		return nil
	}

	pokemon, ok := Pokedex[pokemonName]

	if !ok {
		fmt.Println("you have not caught that pokemon")
	} else {
		fmt.Println("Name: ", pokemon.Name)
		fmt.Println("Height: ", pokemon.Height)
		fmt.Println("Weight: ", pokemon.Weight)
		fmt.Println("Stats:")

		for _, item := range pokemon.Stats {
			fmt.Println("\t-", item.Stat.Name, ": ", item.BaseStat)
		}

		fmt.Println("Types:")

		for _, item := range pokemon.Types {
			fmt.Println("\t-", item.Type.Name)
		}
	}

	return nil
}

func pokedex(cfg *pokeConfig, name string) error {
    for k := range Pokedex {
        fmt.Printf("\t- %s\n", k)
    }
	return nil
}
