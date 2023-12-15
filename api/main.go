package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	app *gin.Engine
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

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)

	file, err := os.Open(dir + "/pokemon.json")
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
		"pokemon": result,
	})
}

func init() {
	app = gin.New()
	app.Group("/api")

	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000", "https://takado-take-home-test.vercel.app"},
		AllowMethods: []string{"GET"},
	}))

	pokemons, err := readFile()
	if err != nil {
		log.Fatal(err)
	}
	app.GET("/search", func(c *gin.Context) {
		searchHandler(c, pokemons)
	})

	app.Run() // listen and serve on 0.0.0.0:8080
}

func Handler(w http.ResponseWriter, r *http.Request) {
	app.ServeHTTP(w, r)
}
