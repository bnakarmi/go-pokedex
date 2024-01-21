package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
)

func (client *Client) ListLocations(pageURL string) (LocationResponse, error) {
	url := baseURL + "/location-area"
	if pageURL != "" {
		url = pageURL
	}

	locationResponse := LocationResponse{}

	if val, ok := client.cache.Get(url); ok {
		err := json.Unmarshal(val, &locationResponse)
		if err != nil {
			return LocationResponse{}, err
		}

		return locationResponse, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return LocationResponse{}, err
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return LocationResponse{}, err
	}

	err = json.Unmarshal(data, &locationResponse)
	if err != nil {
		return LocationResponse{}, err
	}

	client.cache.Add(url, data)

	return locationResponse, nil
}

func (client *Client) ExploreLocation(locationName string) (LocationAreaResponse, error) {
	url := baseURL + "/location-area/" + locationName
	locationAreaResponse := LocationAreaResponse{}

	if val, ok := client.cache.Get(url); ok {
		err := json.Unmarshal(val, &locationAreaResponse)
		if err != nil {
			return LocationAreaResponse{}, err
		}

		return locationAreaResponse, nil
	}

	result, err := http.Get(url)
	if err != nil {
		return locationAreaResponse, err
	}
	defer result.Body.Close()

	data, err := io.ReadAll(result.Body)
	if err != nil {
		return locationAreaResponse, err
	}

	err = json.Unmarshal(data, &locationAreaResponse)
	if err != nil {
		return locationAreaResponse, err
	}

	client.cache.Add(url, data)

	return locationAreaResponse, nil
}

func (client *Client) CatchPokemon(pokemonName string) (Pokemon, error) {
    url := baseURL + "/pokemon/" + pokemonName
    pokemon := Pokemon{}

    result, err := http.Get(url)
    if err != nil {
        return pokemon, err
    }

    defer result.Body.Close()

    data, err := io.ReadAll(result.Body)
    if err != nil {
        return pokemon, err
    }

    err = json.Unmarshal(data, &pokemon)
    if err != nil {
        return pokemon, err
    }

    return pokemon, nil
}
