package mission

import (
	"fmt"
	"time"

	"github.com/luisya22/galactic-exchange/corporation"
	"github.com/luisya22/galactic-exchange/gamecomm"
	"github.com/luisya22/galactic-exchange/world"
)

func (ms *MissionScheduler) CreateSquadMission(m Mission) {
	// TODO: Random events could affect mission times

	// TODO: What happen if squads are on different positions?

	// GET SQUAD
	squadId := m.Squads[0]
	squadResChan := make(chan any)
	corpCommand := gamecomm.CorpCommand{
		Action:          gamecomm.GetSquad,
		ResponseChannel: squadResChan,
		CorporationId:   m.CorporationId,
		SquadIndex:      squadId,
	}
	ms.eventScheduler.gameChannels.CorpChannel <- corpCommand

	squadRes := <-squadResChan
	squad := squadRes.(corporation.Squad)

	close(squadResChan)

	// GET PLANET
	planetResChan := make(chan any)
	planetCommand := gamecomm.WorldCommand{
		PlanetId:        m.PlanetId,
		Action:          gamecomm.GetPlanet,
		ResponseChannel: planetResChan,
	}
	ms.eventScheduler.gameChannels.WorldChannel <- planetCommand

	planetRes := <-planetResChan
	planet := planetRes.(world.Planet)

	close(planetResChan)

	// CALCULATE SHIP SPEED
	shipSpeed := squad.Ships.Speed
	planetDistance := world.Distance(squad.Location, planet.Location)
	destinationTime := planetDistance / float64(shipSpeed)

	// CREATE ARRIVE EVENT
	ae := Event{
		MissionId: m.Id,
		Time:      time.Now().Add(time.Minute * time.Duration(destinationTime)),
		Cancelled: false,
		Execute:   arrivingEvent,
	}

	ms.eventScheduler.Schedule(&ae)

	// CREATE HARVESTING RESOURCES EVENT
	he := Event{
		MissionId: m.Id,
		Time:      m.DestinationTime,
		Cancelled: false,
		Execute:   harvestingEvent,
	}

	ms.eventScheduler.Schedule(&he)

	// CREATE RETURN EVENT
	re := Event{
		MissionId: m.Id,
		Time:      m.ReturnalTime,
		Cancelled: false,
		Execute:   returnEvent,
	}

	ms.eventScheduler.Schedule(&re)

	// Each event should be pushed to the event scheduler and have a way to communicate back to the mission scheduler
}

// - This would send message that we arrive to the mission place
func arrivingEvent(mission *Mission, gameChannels *gamecomm.GameChannels) {
	// TODO: Notify
	// TODO: Calculate danger
}

// - This would Gather the resources
func harvestingEvent(mission *Mission, gameChannels *gamecomm.GameChannels) {

	responseChan := make(chan any)

	for _, resource := range mission.Resources {
		// Generate harvested resourcesAmount

		//TODO: Get Squads to calculate Ship and Crew Bonuses

		// Remove Resources from planet
		gameChannels.WorldChannel <- gamecomm.WorldCommand{
			PlanetId:        mission.PlanetId,
			Action:          gamecomm.AddResourcesToPlanet,
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

	//TODO: Send message

}

// - This would add resources to corporation
func returnEvent(mission *Mission, gameChannels *gamecomm.GameChannels) {
	// responseChan := make(chan any)

	//TODO: Remove resources from squad
	//TODO: Add resources to base
	//TODO: Send message

}
