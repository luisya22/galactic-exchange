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

	// Get Squad
	squad, err := getSquad(mission.CorporationId, mission.Squads[0], gameChannels)
	if err != nil {
		fmt.Println(err.Error()) //TODO: Handle error
	}

	for _, resource := range mission.Resources {
		// Generate harvested resourcesAmount

		bonus := squad.GetHarvestingBonus()
		resourceAmount := 100 * bonus

		// Remove Resources from planet
		removeResourceFromPlanet(mission.PlanetId, resourceAmount, resource, gameChannels)

		// Add Resources to Squad
		addResourcesToSquad(mission.CorporationId, resourceAmount, resource, gameChannels)
	}

	//TODO: Send message

}

func getSquad(corporationId uint64, squadId int, gameChannels *gamecomm.GameChannels) (corporation.Squad, error) {

	squadResChan := make(chan any)
	corpCommand := gamecomm.CorpCommand{
		Action:          gamecomm.GetSquad,
		ResponseChannel: squadResChan,
		CorporationId:   corporationId,
		SquadIndex:      squadId,
	}

	gameChannels.CorpChannel <- corpCommand

	squadRes := <-squadResChan

	close(squadResChan)

	squad, ok := squadRes.(corporation.Squad)
	if !ok {
		return corporation.Squad{}, fmt.Errorf("world channel returned wrong squad object")
	}

	return squad, nil

}

func removeResourceFromPlanet(planetId string, resourceAmount int, resource string, gameChannels *gamecomm.GameChannels) error {
	responseChan := make(chan any)
	gameChannels.WorldChannel <- gamecomm.WorldCommand{
		PlanetId:        planetId,
		Action:          gamecomm.AddResourcesToPlanet,
		Amount:          resourceAmount,
		ResponseChannel: responseChan,
		Resource:        resource,
	}

	responseChanRes := <-responseChan

	err := responseChanRes.(error)
	if err != nil {
		return err // TODO: Handle error correctly
	}

	close(responseChan)

	return nil
}

func addResourcesToSquad(corporationId uint64, resourceAmount int, resource string, gameChannels *gamecomm.GameChannels) error {
	squadResChan := make(chan any)

	gameChannels.CorpChannel <- gamecomm.CorpCommand{
		Action:          gamecomm.AddResourcesToSquad,
		ResponseChannel: squadResChan,
		CorporationId:   corporationId,
		Resource:        resource,
		Amount:          resourceAmount,
	}

	squadRes := <-squadResChan

	err, ok := squadRes.(error)
	if !ok {
		return fmt.Errorf("error on communication with CorpChannel, did not send error")
	}
	if err != nil {
		return err
	}

	close(squadResChan)

	return nil
}

// - This would add resources to corporation
func returnEvent(mission *Mission, gameChannels *gamecomm.GameChannels) {
	// responseChan := make(chan any)

	//TODO: Remove resources from squad
	//TODO: Add resources to base
	//TODO: Send message

}
