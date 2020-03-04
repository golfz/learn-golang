package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

type AnimalFacts struct {
	CommonName     string   `json:"commonName"`
	ScientificName string   `json:"scientificName"`
	HeightInInches int      `json:"heightInInches"`
	FavoriteFoods  []string `json:"favoriteFoods"`
	CanSwim        bool     `json:"canSwim"`
}

func main() {

	jsonData := `
	{
		"CommonName":     "brown-throated three-toed sloth",
		"ScientificName": "Bradypus variegatus",
		"HeightInInches": 31,
		"FavoriteFoods":  ["Cecropia leaves", "Hibiscus flowers"],
		"CanSwim":        true
	}`

	rdr := strings.NewReader(jsonData)

	var sloth AnimalFacts

	if err := json.NewDecoder(rdr).Decode(&sloth); err != nil {
		fmt.Printf("error deserializing JSON: %v", err)
		return
	}

	fmt.Printf("Hi! I'm a %s (%s)! I'm about %d\" tall, and love eating %v!\n",
		sloth.CommonName, sloth.ScientificName, sloth.HeightInInches,
		sloth.FavoriteFoods)
	if sloth.CanSwim {
		fmt.Println("By the way, I can swim!")
	}

}
