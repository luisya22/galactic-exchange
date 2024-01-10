package gamecomm

type Corporation struct {
	ID                              uint64
	Name                            string
	Reputation                      int
	Credits                         float64
	Bases                           []*Base
	CrewMembers                     []*CrewMember
	Squads                          []*Squad
	IsPlayer                        bool
	ReputationWithOtherCorporations map[string]int
}

type Base struct {
	ID                 uint64
	Name               string
	Location           Coordinates
	ResourceProduction map[string]int
	StorageCapacity    float64
	StoredResources    map[string]int
}

type CrewMember struct {
	ID         uint64
	Name       string
	Species    string
	Skills     map[string]int
	AssignedTo uint64
}

type Squad struct {
	Id          uint64
	Ships       Ship
	CrewMembers []CrewMember
	Cargo       map[string]int
	Location    Coordinates
	// Officers []Officers   coming soon...
}

type Ship struct {
	Name         string
	Capacity     int
	MaxHealth    int
	ActualHealth int
	MaxCargo     int
	Location     Coordinates
	Speed        int
	// Attributes
	// Upgrades
	// StoredResources
}
