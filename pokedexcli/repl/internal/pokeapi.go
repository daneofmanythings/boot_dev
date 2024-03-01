package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

var cache Cache = *NewCache(interval)

const (
	locationEndpointBase = "https://pokeapi.co/api/v2/location-area/"
	pokemonEndpointBase  = "https://pokeapi.co/api/v2/pokemon/"
)

const interval = 5 * time.Second

type Direction int

const (
	NEXT Direction = iota
	PREV
)

type locationNavigationEndpoints struct {
	next string
	prev string
}

func GetRawBytesFromApi(endpoint string) []byte {
	res, err := http.Get(fmt.Sprintf(endpoint))
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
	}
	return body
}

func NewMapResponse() func(Direction) (MapApiResponse, error) {
	endpoints := locationNavigationEndpoints{locationEndpointBase, ""}

	return func(dir Direction) (MapApiResponse, error) {
		response := MapApiResponse{}
		var endpoint string

		switch dir {
		case NEXT:
			if endpoints.next == "" {
				return MapApiResponse{}, errors.New("already at the end of the map")
			}
			endpoint = endpoints.next
		case PREV:
			if endpoints.prev == "" {
				return MapApiResponse{}, errors.New("already at the start of the map")
			}
			endpoint = endpoints.prev
		default:
			fmt.Println("Invalid direction used in fetching api")
		}
		value, ok := cache.Get(endpoint)

		if !ok {
			responseBytes := GetRawBytesFromApi(endpoint)
			cache.Add(endpoint, responseBytes)
			value, _ = cache.Get(endpoint)
		}

		err := json.Unmarshal(value, &response)
		if err != nil {
			fmt.Println(err)
		}

		endpoints.next = response.Next
		endpoints.prev = response.Previous

		return response, nil
	}
}

var MapResponse func(Direction) (MapApiResponse, error) = NewMapResponse()

func ExploreResponse(area string) (ExploreApiResponse, error) {
	response := ExploreApiResponse{}
	endpoint := locationEndpointBase + area
	value, ok := cache.Get(endpoint)

	if !ok {
		responseBytes := GetRawBytesFromApi(endpoint)
		cache.Add(endpoint, responseBytes)
		value, _ = cache.Get(endpoint)
	}

	err := json.Unmarshal(value, &response)
	if err != nil {
		fmt.Println(err)
	}

	return response, nil
}

func CatchResponse(name string) (PokemonApiResponse, error) {
	response := PokemonApiResponse{}
	endpoint := pokemonEndpointBase + name
	value, ok := cache.Get(endpoint)

	if !ok {
		responseBytes := GetRawBytesFromApi(endpoint)
		cache.Add(endpoint, responseBytes)
		value, _ = cache.Get(endpoint)
	}

	err := json.Unmarshal(value, &response)
	if err != nil {
		fmt.Println(err)
	}

	return response, nil
}
