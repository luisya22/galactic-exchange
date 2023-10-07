package game

import (
	"fmt"

	"github.com/luisya22/galactic-exchange/world"
)

// TODO: Amounts should be reflected after the time distance is elapsed
func (g *Game) SellResource(amount int, itemName world.Resource, planetId string, corporationId uint64) error {
	var item world.ResourceInfo
	var planet world.Planet
	var playerResourceAmount int
	var ok bool

	// Use functions to access
	if item, ok = g.World.AllResources[itemName]; !ok {
		return fmt.Errorf("Item not found")
	}

	planet, err := g.World.GetPlanet(planetId)
	if err != nil {
		return err
	}

	if playerResourceAmount, ok = g.PlayerState.Corporation.Bases[0].StoredResources[item.Name]; !ok {
		return fmt.Errorf("You don't have enough %v: 0", item.Name)
	}

	if playerResourceAmount < amount {
		return fmt.Errorf("You don't have enough %s: %d", item.Name, playerResourceAmount)
	}

	planetAmount := planet.Resources[item.Name]

	// Rmove resources from corporation
	_, err = g.Corporations.RemoveResources(corporationId, itemName, amount)
	if err != nil {
		return err
	}
	//Add resrouces to planet
	_, err = g.World.AddResourcesToPlanet(planetId, itemName, amount)
	if err != nil {
		// Undo Step 1
		_, undoErr := g.Corporations.AddResources(corporationId, itemName, amount)
		if undoErr != nil {
			// Handle error during undo
			return fmt.Errorf("error adding resources to planet: %v, and error undoing removal of resources: %v", err, undoErr)
		}
		return err
	}
	// Add credits to corporation
	_, err = g.Corporations.AddCredits(1, item.BasePrice*float64(amount))
	if err != nil {
		// Undo Step 2
		_, undoErr1 := g.World.RemoveResourcesFromPlanet(planetId, itemName, amount)
		// Undo Step 1
		_, undoErr2 := g.Corporations.AddResources(corporationId, itemName, amount)
		if undoErr1 != nil || undoErr2 != nil {
			// Handle error during undo
			return fmt.Errorf("error adding credits: %v, and error undoing previous steps: %v, %v", err, undoErr1, undoErr2)
		}

		return err
	}

	fmt.Printf("Planet: %v\n", planet.Resources)

	// TODO: Use Mutexes to access
	//TODO: This info should be returned as copies
	fmt.Printf(
		"Transfer:\nPlanet: %v -> %v\nCorporation: %v -> %v\nSell Price: %v\nNew Player Credits Balance: %v\n",
		planetAmount,
		g.World.Planets[planet.Name].Resources[item.Name],
		playerResourceAmount,
		g.PlayerState.Corporation.Bases[0].StoredResources[item.Name],
		item.BasePrice,
		g.PlayerState.Corporation.Credits,
	)

	return nil

}
