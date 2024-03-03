package world

import "math/rand"

type ResourceProfile struct {
	Primary   string
	Secondary string
}

type Rarity int

const (
	Abundant Rarity = iota
	Common
	Scarce
	Rare
)

type ResourceInfo struct {
	Name      string
	BasePrice float64
}

func shouldIncludeResource(world *World, res string, planet *Planet) bool {
	rarity := world.ResourceRarity[res]
	switch rarity {
	case Abundant:
		return true
	case Common:
		return planet.DangerLevel >= 20
	case Scarce:
		return planet.DangerLevel >= 40
	case Rare:
		return planet.DangerLevel >= 70
	default:
		return false
	}
}

func GenerateResourceProfile() ResourceProfile {
	resources := []string{"gold", "iron", "water", "food"}

	rand.Shuffle(len(resources), func(i, j int) {
		resources[i], resources[j] = resources[j], resources[i]
	})

	return ResourceProfile{
		Primary:   resources[0],
		Secondary: resources[1],
	}
}

func CreateWorldResources() map[string]ResourceInfo {
	return map[string]ResourceInfo{
		"gold": {
			Name:      "gold",
			BasePrice: 250,
		},
		"iron": {
			Name:      "iron",
			BasePrice: 200,
		},
		"water": {
			Name:      "water",
			BasePrice: 10,
		},
		"food": {
			Name:      "food",
			BasePrice: 10,
		},
	}
}
