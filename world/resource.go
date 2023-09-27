package world

import "math/rand"

type Resource string

type ResourceProfile struct {
	Primary   Resource
	Secondary Resource
}

type Rarity int

const (
	Abundant Rarity = iota
	Common
	Scarce
	Rare
)

const (
	Gold  Resource = "gold"
	Iron  Resource = "iron"
	Water Resource = "water"
	Food  Resource = "food"
)

type ResourceInfo struct {
	Name      Resource
	BasePrice float64
}

func shouldIncludeResource(world World, res Resource, planet Planet) bool {
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
	resources := []Resource{Gold, Iron, Water, Food}

	rand.Shuffle(len(resources), func(i, j int) {
		resources[i], resources[j] = resources[j], resources[i]
	})

	return ResourceProfile{
		Primary:   resources[0],
		Secondary: resources[1],
	}
}

// map[Resource]struct{}{
// 		Gold:  {},
// 		Iron:  {},
// 		Water: {},
// 		Food:  {},
// 	}
//

func CreateWorldResources() map[Resource]ResourceInfo {
	return map[Resource]ResourceInfo{
		Gold: {
			Name:      Gold,
			BasePrice: 250,
		},
		Iron: {
			Name:      Iron,
			BasePrice: 200,
		},
		Water: {
			Name:      Water,
			BasePrice: 10,
		},
		Food: {
			Name:      Food,
			BasePrice: 10,
		},
	}
}
