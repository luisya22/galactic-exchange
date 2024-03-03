package world

import (
	"math/rand"

	"github.com/luisya22/galactic-exchange/internal/resource"
)

type ResourceProfile struct {
	Primary   string
	Secondary string
}

type ResourceInfo struct {
	Name      string
	BasePrice float64
}

func shouldIncludeResource(world *World, res resource.Resource, planet *Planet) bool {

	switch res.Rarity {
	case resource.Abundant:
		return true
	case resource.Common:
		return planet.DangerLevel >= 20
	case resource.Scarce:
		return planet.DangerLevel >= 40
	case resource.Rare:
		return planet.DangerLevel >= 70
	default:
		return false
	}
}

func GenerateResourceProfile(worldResources map[string]resource.Resource) ResourceProfile {

	// TODO: improve this later to not use []string but map[string]Resource
	resources := []string{}
	for _, r := range worldResources {
		resources = append(resources, r.Name)
	}

	rand.Shuffle(len(resources), func(i, j int) {
		resources[i], resources[j] = resources[j], resources[i]
	})

	return ResourceProfile{
		Primary:   resources[0],
		Secondary: resources[1],
	}
}
