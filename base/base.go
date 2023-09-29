package base

import "github.com/luisya22/galactic-exchange/world"

type Base struct {
	ID                 uint64
	Name               string
	Location           world.Coordinates
	ResourceProduction map[world.Resource]int
	StorageCapacity    float64
	StoredResources    map[world.Resource]int
}
