package mission

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/luisya22/galactic-exchange/channel"
	"github.com/luisya22/galactic-exchange/world"
)

type MissionScheduler struct {
	missions       map[string]*Mission
	eventScheduler *EventScheduler
	missionChannel chan channel.MissionCommand
	RW             sync.RWMutex
}

type Mission struct {
	Id              string
	CorporationId   uint64
	Squads          []int
	PlanetId        string
	DestinationTime time.Time
	ReturnalTime    time.Time
	Status          string
	Type            channel.MissionType
	Resources       []string
}

// Misión va a un planeta ha hacer algo
// Corporation
// Squad
// Planet
// Type

// startMission
// missionStatus -> createEvents
// events -> sendMissionStatus, after each event send notification
// if all events received, end mission and send notification

func NewMissionScheduler(gameChannels *channel.GameChannels) *MissionScheduler {

	eventScheduler := NewEventScheduler(gameChannels)
	return &MissionScheduler{
		missions:       make(map[string]*Mission, 0),
		eventScheduler: eventScheduler,
		missionChannel: gameChannels.MissionChannel,
	}
}

func (ms *MissionScheduler) Run() {

	go ms.eventScheduler.Run()

	for m := range ms.missionChannel {

		mission, err := CreateMission(m)
		if err != nil {
			continue
		}

		ms.StartMission(mission)
	}
}

func CreateMission(mc channel.MissionCommand) (Mission, error) {

	uuid, err := uuid.NewUUID()
	if err != nil {
		return Mission{}, fmt.Errorf("Error: %v", err)
	}

	missionId := uuid.String()

	mission := Mission{
		Id:              missionId,
		CorporationId:   mc.CorporationId,
		Squads:          mc.Squads,
		PlanetId:        mc.PlanetId,
		DestinationTime: mc.DestinationTime,
		ReturnalTime:    mc.ReturnalTime,
		Status:          mc.Status,
		Type:            mc.Type,
		Resources:       mc.Resources,
	}

	return mission, nil
}

func (ms *MissionScheduler) StartMission(m Mission) {
	ms.RW.Lock()
	ms.missions[m.Id] = &m
	ms.RW.Unlock()

	switch m.Type {
	case channel.SquadMission:
		ms.CreateSquadMission(m)
	default:
		ms.RW.Lock()
		delete(ms.missions, m.Id)
		ms.RW.Unlock()
	}

}

// Id        string
// MissionId string
// Time      time.Time
// Cancelled bool
// Index     int
// Execute   func()

func (ms *MissionScheduler) CreateSquadMission(m Mission) {
	// Create Arrive event with function
	// - This would send message that we arrive to the mission place

	// Create Gather event with function
	// - This would Gather the resources
	harvestingEvent := Event{
		MissionId: m.Id,
		Time:      time.Now().Add(1 * time.Minute),
		Cancelled: false,
		Execute: func(gameChannels *channel.GameChannels) {

			responseChan := make(chan any)

			for _, resource := range m.Resources {
				// Generate harvested resourcesAmount

				//TODO: Get Squads to calculate Ship and Crew Bonuses

				// Remove Resources from planet
				gameChannels.WorldChannel <- channel.WorldCommand{
					PlanetId:        m.PlanetId,
					Action:          channel.AddResourcesToPlanet,
					Amount:          100,
					ResponseChannel: responseChan,
					Resource:        resource,
				}

				// Add Resources to Squad
				// TODO: Add Resources to Squad on the Corporation
			}

			// TODO: Check the thing with the locks and copies
			res := <-responseChan

			worldResponse := res.(world.WorldResponse)

			close(responseChan)

			fmt.Println(worldResponse.Planet)
		},
	}

	ms.eventScheduler.Schedule(&harvestingEvent)

	// Create Return event with function
	// - This would add resources to corporation

	//TODO: Add Resources to Base

	// Each event should be pushed to the event scheduler and have a way to communicate back to the mission scheduler
}
