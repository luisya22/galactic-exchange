package resource

import (
	"encoding/json"
	"log"

	"github.com/luisya22/galactic-exchange/internal/gamedata"
)

type Resource struct {
	Name      string  `json:"name"`
	BasePrice float64 `json:"basePrice"`
	Rarity    Rarity  `json:"rarity"`
}

type Rarity int

const (
	Abundant Rarity = iota
	Common
	Scarce
	Rare
)

func LoadWorldResources() map[string]Resource {

	resources := make(map[string]Resource, 4)

	file, err := gamedata.Files.Open("resourcedata/resources.json")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&resources)
	if err != nil {
		log.Fatal(err.Error())
	}

	return resources
}
