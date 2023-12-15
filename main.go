package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

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

func readFile() *PokemonList {
	file, err := os.Open("pokemon.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	var pokemons PokemonList
	err = json.Unmarshal(data, &pokemons)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("pokemons: ", pokemons)
	return &pokemons
}

func main() {
	r := gin.Default()

	r.GET("/pokemon", func(c *gin.Context) {
		pokemons := readFile()
		c.JSON(200, pokemons)
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}
