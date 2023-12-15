package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type Pokemon struct {
	Name  string `json:"name"`
	Url   string `json:"url"`
	Image string `json:"image"`
}

type PokemonList struct {
	Pokemons []Pokemon `json:"pokemonList"`
}

func readFile() (*PokemonList, error) {
	file, err := os.Open("pokemon.json")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	var pokemons PokemonList
	err = json.Unmarshal(data, &pokemons)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &pokemons, nil
}

func filterPokemon(pokemons *PokemonList, query string) []Pokemon {
	var result []Pokemon
	for _, pokemon := range pokemons.Pokemons {

		trimQuery := strings.TrimSpace(query)
		if strings.Contains(strings.ToLower(pokemon.Name), strings.ToLower(trimQuery)) {
			result = append(result, pokemon)
		}
	}

	return result
}

func searchHandler(c *gin.Context, pokemons *PokemonList) {
	query := c.Query("query")

	result := filterPokemon(pokemons, query)

	c.JSON(200, gin.H{
		"result": result,
	})
}

func main() {
	r := gin.Default()
	pokemons, err := readFile()
	if err != nil {
		log.Fatal(err)
	}

	r.GET("/search", func(c *gin.Context) {
		searchHandler(c, pokemons)
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}
