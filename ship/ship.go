package ship

import (
	"github.com/luisya22/galactic-exchange/world"
)

type Ship struct {
	Name         string
	Capacity     int
	MaxHealth    int
	ActualHealth int
	MaxCargo     int
	Location     world.Coordinates
	// Attributes
	// Upgrades
	// StoredResources
}
