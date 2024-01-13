package ship

import (
	"github.com/luisya22/galactic-exchange/internal/gamecomm"
	"github.com/luisya22/galactic-exchange/internal/world"
)

type Ship struct {
	Name         string
	Capacity     int
	MaxHealth    int
	ActualHealth int
	MaxCargo     int
	Location     world.Coordinates
	Speed        int
	// Attributes
	// Upgrades
	// StoredResources
}

func (s *Ship) Copy() gamecomm.Ship {

	location := gamecomm.Coordinates{
		X: s.Location.X,
		Y: s.Location.Y,
	}

	return gamecomm.Ship{
		Name:         s.Name,
		Capacity:     s.Capacity,
		MaxHealth:    s.MaxHealth,
		ActualHealth: s.ActualHealth,
		MaxCargo:     s.MaxCargo,
		Location:     location,
		Speed:        s.Speed,
	}
}
