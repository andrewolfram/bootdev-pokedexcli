package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	BaseURL              = "https://pokeapi.co/api/v2"
	PokemonEndpoint      = "pokemon"
	LocationAreaEndpoint = "location-area"
)

//TODO: Guess that should have a generic JSON parse

// fetchPokeLocation sends a GET request to the given URL and parses the JSON response
func FetchPokeLocation(url string) (PokeLocation, error) {
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		return PokeLocation{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return PokeLocation{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return PokeLocation{}, err
	}

	var loc PokeLocation
	if err := json.Unmarshal(body, &loc); err != nil {
		return PokeLocation{}, err
	}

	return loc, nil
}

func FetchPokeLocationDetail(url string) (PokeLocationDetails, error) {
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		return PokeLocationDetails{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return PokeLocationDetails{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return PokeLocationDetails{}, err
	}

	var loc PokeLocationDetails
	if err := json.Unmarshal(body, &loc); err != nil {
		return PokeLocationDetails{}, err
	}

	return loc, nil
}

func FetchPokemonDetail(url string) (PokemonDetails, error) {
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		return PokemonDetails{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return PokemonDetails{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return PokemonDetails{}, err
	}

	var loc PokemonDetails
	if err := json.Unmarshal(body, &loc); err != nil {
		return PokemonDetails{}, err
	}

	return loc, nil
}
