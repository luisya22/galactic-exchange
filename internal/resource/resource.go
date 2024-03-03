package resource

type Resource struct {
	Name      string
	BasePrice float64
	Rarity    Rarity
}

type Rarity int

const (
	Abundant Rarity = iota
	Common
	Scarce
	Rare
)

func CreateWorldResources() map[string]Resource {
	return map[string]Resource{
		"gold": {
			Name:      "gold",
			BasePrice: 250,
			Rarity:    Common,
		},
		"iron": {
			Name:      "iron",
			BasePrice: 200,
			Rarity:    Common,
		},
		"water": {
			Name:      "water",
			BasePrice: 10,
			Rarity:    Scarce,
		},
		"food": {
			Name:      "food",
			BasePrice: 10,
			Rarity:    Rare,
		},
	}
}
