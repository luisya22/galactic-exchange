package gamecomm

import (
	"time"
)

type GameChannelType string

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
	ResponseChannel chan any
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
	ResponseChannel chan any
	CorporationId   uint64
	SquadIndex      int
}

type CommandType int

const (
	GetCorporation CommandType = iota
	GetSquad
	AddResourcesToSquad
	RemoveResourcesFromSquad
	AddResourcesToBase
	RemoveResourcesFromBase
)

// Mission Channels
type MissionCommand struct {
	Id              string
	CorporationId   uint64
	Squads          []int
	PlanetId        string
	DestinationTime time.Time
	ReturnalTime    time.Time
	Status          string
	Type            MissionType
	Resources       []string
}

type MissionType int

const (
	SquadMission MissionType = iota
	QuestMission
)
