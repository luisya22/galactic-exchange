package corporation

import "github.com/luisya22/galactic-exchange/internal/world"

type Base struct {
	ID                 uint64
	Name               string
	Location           world.Coordinates
	ResourceProduction map[string]int
	StorageCapacity    float64
	StoredResources    map[string]int
}
