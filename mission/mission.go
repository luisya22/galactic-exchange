package mission

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/luisya22/galactic-exchange/gameclock"
	"github.com/luisya22/galactic-exchange/gamecomm"
)

type MissionScheduler struct {
	Missions       map[string]*Mission
	EventScheduler EventScheduler
	MissionChannel chan gamecomm.MissionCommand
	RW             sync.RWMutex
	GameClock      *gameclock.GameClock
	GameChannels   *gamecomm.GameChannels
	ErrorChan      chan error
}

type Mission struct {
	Id               string
	CorporationId    uint64
	Squads           []int
	PlanetId         string
	DestinationTime  time.Time
	ReturnalTime     time.Time
	Status           string
	Type             gamecomm.MissionType
	Resources        []string
	Amount           int // TODO: You should have and object for transfers {resource, amount}
	NotificationChan chan string
	ErrorChan        chan error
}

// startMission
// missionStatus -> createEvents
// events -> sendMissionStatus, after each event send notification
// if all events received, end mission and send notification

func NewMissionScheduler(gameChannels *gamecomm.GameChannels, gc *gameclock.GameClock) *MissionScheduler {

	missions := make(map[string]*Mission, 0)
	eventScheduler := NewEventScheduler(gameChannels, missions, gc)

	return &MissionScheduler{
		Missions:       missions,
		EventScheduler: eventScheduler,
		MissionChannel: gameChannels.MissionChannel,
		GameClock:      gc,
		GameChannels:   gameChannels,
	}
}

// TODO: Add Wait Group
func (ms *MissionScheduler) Run() {

	go ms.EventScheduler.Run()
	go ms.handleErrors()

	for m := range ms.MissionChannel {

		mission, err := CreateMission(m, ms.ErrorChan)
		if err != nil {
			continue
		}

		ms.StartMission(mission)
	}
}

// TODO: Correctly handle errors
func (ms *MissionScheduler) handleErrors() {
	for err := range ms.ErrorChan {
		fmt.Println(err)
	}
}

func CreateMission(mc gamecomm.MissionCommand, errorChan chan error) (Mission, error) {

	uuid, err := uuid.NewUUID()
	if err != nil {
		return Mission{}, fmt.Errorf("error: %v", err)
	}

	missionId := uuid.String()

	mission := Mission{
		Id:               missionId,
		CorporationId:    mc.CorporationId,
		Squads:           mc.Squads,
		PlanetId:         mc.PlanetId,
		DestinationTime:  mc.DestinationTime,
		ReturnalTime:     mc.ReturnalTime,
		Status:           "In Progress",
		Type:             mc.Type,
		Resources:        mc.Resources,
		NotificationChan: mc.NotificationChan,
		Amount:           mc.Amount,
		ErrorChan:        errorChan,
	}

	return mission, nil
}

func (ms *MissionScheduler) StartMission(m Mission) {
	ms.RW.Lock()
	ms.Missions[m.Id] = &m
	ms.RW.Unlock()

	switch m.Type {
	case gamecomm.SquadMission:
		err := ms.CreateSquadMission(m)
		if err != nil {
			delete(ms.Missions, m.Id)
			m.NotificationChan <- err.Error()
		}
	case gamecomm.TransferMission:
		err := ms.CreateTransferMission(m)
		if err != nil {
			delete(ms.Missions, m.Id)
			m.NotificationChan <- err.Error()
		}
	default:
		ms.RW.Lock()
		delete(ms.Missions, m.Id)
		ms.RW.Unlock()
	}
}

func (msz *MissionScheduler) CalculateTravelDistance(corporationId uint64, squads []int, planetId string, gameChannels *gamecomm.GameChannels) (float64, error) {

	if len(squads) == 0 {
		return 0.0, fmt.Errorf("error: should include squads")
	}

	squadId := squads[0]
	squad, err := getSquad(corporationId, squadId, gameChannels)
	if err != nil {
		return 0.0, err
	}

	// GET PLANET
	planetResChan := make(chan gamecomm.ChanResponse)
	planetCommand := gamecomm.WorldCommand{
		PlanetId:        planetId,
		Action:          gamecomm.GetPlanet,
		ResponseChannel: planetResChan,
	}
	gameChannels.WorldChannel <- planetCommand

	planetRes := <-planetResChan
	if planetRes.Err != nil {
		return 0.0, planetRes.Err
	}

	planet := planetRes.Val.(gamecomm.Planet)
	close(planetResChan)

	// CALCULATE SHIP SPEED
	shipSpeed := squad.Ships.Speed
	squadLocation := gamecomm.Coordinates{X: squad.Location.X, Y: squad.Location.Y}
	planetLocation := gamecomm.Coordinates{X: planet.Location.X, Y: planet.Location.Y}

	planetDistance := gamecomm.Distance(squadLocation, planetLocation)
	_ = planetDistance / float64(shipSpeed)

	return planetDistance, nil
}
