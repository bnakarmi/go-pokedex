package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

func (c *HelpCommand) Execute(cfg *PokeConfig) {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	for _, cmd := range getCommand() {
		fmt.Printf("  %s:\n", cmd.name)
		fmt.Printf("\t%s\n", cmd.desc)
	}

	fmt.Println()
}

func (c *ClearCommand) Execute(cfg *PokeConfig) {
	var cmd *exec.Cmd
	cmd = exec.Command("clear")
	cmd.Stdout = os.Stdout

	cmd.Run()
}

func (c *ExitCommand) Execute(cfg *PokeConfig) {
	os.Exit(0)
}

func (c *MapForwardCommand) Execute(cfg *PokeConfig) {
	resp, err := cfg.pokeClient.ListLocations(cfg.nextLocationURL)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	cfg.nextLocationURL = resp.Next
	cfg.prevLocationURL = resp.Previous

	for _, loc := range resp.Results {
		fmt.Println(loc.Name)
	}
}

func (c *MapBackCommand) Execute(cfg *PokeConfig) {
	if cfg.prevLocationURL == "" {
		fmt.Println("You are already on the first page")
		return
	}

	resp, err := cfg.pokeClient.ListLocations(cfg.prevLocationURL)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	cfg.nextLocationURL = resp.Next
	cfg.prevLocationURL = resp.Previous

	for _, loc := range resp.Results {
		fmt.Println(loc.Name)
	}
}

func (c *ExploreCommand) Execute(cfg *PokeConfig) {
	resp, err := cfg.pokeClient.ExploreLocation(c.locationName)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	for _, pokemonEncounters := range resp.PokemonEncounters {
		fmt.Println(pokemonEncounters.Pokemon.Name)
	}
}

func (c *CatchCommand) Execute(cfg *PokeConfig) {
	pokemonName := c.pokemonName
	if pokemonName == "" {
		fmt.Println("Provide a pokemon to catch.")
		return
	}

	fmt.Println("throwing a ball at ", pokemonName, "...")

	resp, err := cfg.pokeClient.CatchPokemon(pokemonName)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	seed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(seed)

	randomNumber := r.Intn(1001)

	if randomNumber < resp.BaseExperience {
		fmt.Println(pokemonName, " escaped")
		return
	}

	fmt.Println(pokemonName, " was caught")
	Pokedex[pokemonName] = resp
}

func (c *InspectCommand) Execute(cfg *PokeConfig) {
	pokemonName := c.pokemonName
	if pokemonName == "" {
		fmt.Println("Provide a pokemon to inspect.")
		return
	}

	pokemon, ok := Pokedex[pokemonName]

	if !ok {
		fmt.Println("you have not caught that pokemon")
		return
	}

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

func (c *PokedexCommand) Execute(cfg *PokeConfig) {
	for k := range Pokedex {
		fmt.Printf("\t- %s\n", k)
	}
}
