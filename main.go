package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/andrewolfram/pokedexcli/internal/pokeapi"
	"github.com/andrewolfram/pokedexcli/internal/pokecache"
)

const (
	//set catch probability to 1
	cheat_mode bool = true
)

func main() {
	pokedex = map[string]pokeapi.PokemonDetails{}
	cmdMap = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Prints this help section",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Fetches the next 20 Locations from the Pokemon World",
			callback:    commandMapForward,
		},
		"mapb": {
			name:        "mapb",
			description: "Fetches the previous 20 Locations from the Pokemon World",
			callback:    commandMapBack,
		},
		"explore": {
			name:        "explore",
			description: "Lists all Pokemon found in the given area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Attempt to catch the given Pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Show info about caught Pokemon",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Show info about all caught Pokemon",
			callback:    commandPokedex,
		},
	}

	baseURL := pokeapi.BaseURL + "/" + pokeapi.LocationAreaEndpoint + "?offset=0&limit=20" //otherwise our cache doesn't work for the first call
	cfg := &cmdConfig{
		nextURL:     &baseURL,
		previousURL: nil,
	}
	cache = pokecache.NewCache(5 * time.Minute)
	s := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex >")
		s.Scan()
		input := s.Text()
		parts := cleanInput(input)
		cmd, ok := cmdMap[parts[0]]
		param := ""
		if len(parts) > 1 {
			param = parts[1]
		}
		if !ok {
			fmt.Println("Unknown command")
		} else {
			cmd.callback(cfg, param)
		}
	}
}
