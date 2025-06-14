package main

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand/v2"
	"os"
	"strings"

	"github.com/andrewolfram/pokedexcli/internal/pokeapi"
	"github.com/andrewolfram/pokedexcli/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*cmdConfig, string) error
}

type cmdConfig struct {
	nextURL     *string
	previousURL *string
}

var cmdMap map[string]cliCommand
var cache *pokecache.Cache

func cleanInput(text string) []string {
	trimmed := strings.TrimSpace(text)
	trimmed = strings.ToLower(trimmed)
	fields := strings.Fields(trimmed)
	return fields
}

func commandExit(cfg *cmdConfig, unused string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *cmdConfig, unused string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	for _, cmd := range cmdMap {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandMapForward(cfg *cmdConfig, unused string) error {
	if cfg.nextURL == nil {
		fmt.Println("you're on the last page")
		return nil
	}
	return mapMove(cfg, *cfg.nextURL)
}

func commandMapBack(cfg *cmdConfig, unused string) error {
	if cfg.previousURL == nil {
		fmt.Println("you're on the first page")
		return nil
	}
	return mapMove(cfg, *cfg.previousURL)
}

func mapMove(cfg *cmdConfig, targetURL string) error {
	data, ok := cache.Get(targetURL)
	var pokeLocation pokeapi.PokeLocation
	if !ok {
		fmt.Println("Cache fail")
		var err error
		pokeLocation, err = pokeapi.FetchPokeLocation(targetURL)
		if err != nil {
			return err
		}
		//fmt.Println(pokeLocation)
		bytes, err := json.Marshal(pokeLocation)
		if err != nil {
			fmt.Printf("failed to marshal PokeLocation: %v", err)
			return err
		}
		cache.Add(targetURL, bytes)
	} else {
		fmt.Println("Cache hit!")
		if err := json.Unmarshal(data, &pokeLocation); err != nil {
			fmt.Printf("failed to unmarshal PokeLocation: %v", err)
			return err
		}
	}
	fmt.Println("Updating config")
	//fmt.Println(pokeLocation)
	//update the next and previous
	cfg.nextURL = pokeLocation.Next
	cfg.previousURL = pokeLocation.Previous
	for _, result := range pokeLocation.Results {
		fmt.Println(result.Name)
	}
	return nil
}

func commandExplore(cfg *cmdConfig, area string) error {
	targetURL := pokeapi.BaseURL + "/" + pokeapi.LocationAreaEndpoint + "/" + area
	data, ok := cache.Get(targetURL)
	var pokeLocation pokeapi.PokeLocationDetails
	if !ok {
		fmt.Println("Cache fail")
		var err error
		pokeLocation, err = pokeapi.FetchPokeLocationDetail(targetURL)
		if err != nil {
			return err
		}
		//fmt.Println(pokeLocation)
		bytes, err := json.Marshal(pokeLocation)
		if err != nil {
			fmt.Printf("failed to marshal PokeLocationDetails: %v", err)
			return err
		}
		cache.Add(targetURL, bytes)
	} else {
		fmt.Println("Cache hit!")
		if err := json.Unmarshal(data, &pokeLocation); err != nil {
			fmt.Printf("failed to unmarshal PokeLocationDetails: %v", err)
			return err
		}
	}
	for _, pkmn := range pokeLocation.PokemonEncounters {
		fmt.Println(pkmn.Pokemon.Name)
	}
	return nil
}

func commandCatch(cfg *cmdConfig, pokemon string) error {
	_, caught := pokedex[pokemon]
	if caught {
		fmt.Printf("You already caught %s...\n", pokemon)
		return nil
	}

	targetURL := pokeapi.BaseURL + "/" + pokeapi.PokemonEndpoint + "/" + pokemon
	pokemonDetails, err := pokeapi.FetchPokemonDetail(targetURL)
	if err != nil {
		return err
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon)
	//I guess we use this for prob calculations or something
	maxBaseXp := 306
	baseExp := int(math.Min(float64(pokemonDetails.BaseExperience), float64(maxBaseXp)))
	roll := rand.IntN(maxBaseXp + 1)
	fmt.Printf("Roll: %d, baseXP: %d\n", 2*roll, baseExp)
	if 2*roll >= baseExp {
		fmt.Printf("%s was caught!\n", pokemon)
		pokedex[pokemon] = pokemonDetails
	} else {
		fmt.Printf("%s escaped...\n", pokemon)
	}
	return nil
}

func commandInspect(cfg *cmdConfig, pokemon string) error {
	pkmnDetails, caught := pokedex[pokemon]
	if !caught {
		fmt.Printf("You have not caught that pokemon\n", pokemon)
		return nil
	}
	fmt.Printf("Name: %s\n", pkmnDetails.Name)
	fmt.Printf("Height: %d\n", pkmnDetails.Height)
	fmt.Printf("Weight: %d\n", pkmnDetails.Weight)
	fmt.Println("Stats:")
	for _, stat := range pkmnDetails.Stats {
		fmt.Printf("\t-%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("StaTypests:")
	for _, t := range pkmnDetails.Types {
		fmt.Printf("\t-%s\n", t.Type.Name)
	}
	return nil
}

func commandPokedex(cfg *cmdConfig, pokemon string) error {
	fmt.Println("Your Pokedex:")
	if len(pokedex) <= 0 {
		fmt.Println("is empty :(")
	} else {
		for _, val := range pokedex {
			fmt.Printf("\t-%s\n", val.Name)
		}
	}
	return nil
}
