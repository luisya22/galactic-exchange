package gamecomm

import (
	"time"
)

type GameChannelType string

type ChanResponse struct {
	Val any
	Err error
}

const (
	WorldChan   = "WorldChan"
	CorpChan    = "CorpChan"
	MissionChan = "MissionChan"
)

type GameChannels struct {
	WorldChannel   chan WorldCommand
	CorpChannel    chan CorpCommand
	MissionChannel chan MissionCommand
}

// World Channels
type WorldCommand struct {
	PlanetId        string
	Action          WorldCommandType
	Amount          int
	ResponseChannel chan ChanResponse
	Resource        string
}

type WorldCommandType int

const (
	GetPlanet WorldCommandType = iota
	AddResourcesToPlanet
	RemoveResourcesFromPlanet
)

// Corporation Channels
type CorpCommand struct {
	Action          CommandType
	ResponseChannel chan ChanResponse
	CorporationId   uint64
	SquadIndex      int
	BaseIndex       int
	Resource        string
	Amount          int
	AmountDecimal   float64
}

type CommandType int

const (
	GetCorporation CommandType = iota
	GetSquad
	AddResourcesToSquad
	RemoveResourcesFromSquad
	RemoveAllResourcesFromSquad
	AddResourcesToBase
	RemoveResourcesFromBase
	AddCredits
	RemoveCredits
)

// Mission Channels
type MissionCommand struct {
	Id               string
	CorporationId    uint64
	Squads           []int
	PlanetId         string
	DestinationTime  time.Time
	ReturnalTime     time.Time
	Status           string
	Type             MissionType
	Resources        []string
	NotificationChan chan string
	Amount           int
}

type MissionType int

const (
	SquadMission MissionType = iota
	QuestMission
	TransferMission
)
