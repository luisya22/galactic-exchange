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
	squad, err := getSquad(m.CorporationId, squadId, ms.eventScheduler.gameChannels)
	if err != nil {
		fmt.Println(err.Error()) //TODO: Handle error
		return
	}

	// GET PLANET
	planetResChan := make(chan gamecomm.ChanResponse)
	planetCommand := gamecomm.WorldCommand{
		PlanetId:        m.PlanetId,
		Action:          gamecomm.GetPlanet,
		ResponseChannel: planetResChan,
	}
	ms.eventScheduler.gameChannels.WorldChannel <- planetCommand

	planetRes := <-planetResChan
	if planetRes.Err != nil {
		fmt.Println(planetRes.Err.Error())
		return
	}

	planet := planetRes.Val.(world.Planet)
	close(planetResChan)

	// CALCULATE SHIP SPEED
	shipSpeed := squad.Ships.Speed
	planetDistance := world.Distance(squad.Location, planet.Location)
	_ = planetDistance / float64(shipSpeed)

	// CREATE ARRIVE EVENT
	ae := Event{
		MissionId: m.Id,
		Time:      time.Now(),
		Cancelled: false,
		Execute:   arrivingEvent,
	}

	ms.eventScheduler.Schedule(&ae)

	// CREATE HARVESTING RESOURCES EVENT
	he := Event{
		MissionId: m.Id,
		Time:      time.Now().Add(15 * time.Second),
		Cancelled: false,
		Execute:   harvestingEvent,
	}

	ms.eventScheduler.Schedule(&he)

	// CREATE RETURN EVENT
	re := Event{
		MissionId: m.Id,
		Time:      time.Now().Add(30 * time.Second),
		Cancelled: false,
		Execute:   returnEvent,
	}

	ms.eventScheduler.Schedule(&re)
}

// - This would send message that we arrive to the mission place
func arrivingEvent(mission *Mission, gameChannels *gamecomm.GameChannels) {
	mission.NotificationChan <- fmt.Sprintf("Mission Notification: Squad %v, reached destination.", mission.Squads)
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
		err := removeResourceFromPlanet(mission.PlanetId, resourceAmount, resource, gameChannels)
		if err != nil {
			fmt.Println(err)
		}

		// Add Resources to Squad
		err = addResourcesToSquad(mission.CorporationId, resourceAmount, resource, gameChannels)
		if err != nil {
			fmt.Println(err)
		}
	}

	mission.NotificationChan <- fmt.Sprintf("Mission Notification: Squad %v, finished harvesting.", mission.Squads)

}

// TODO: Close channels on producers

func getSquad(corporationId uint64, squadId int, gameChannels *gamecomm.GameChannels) (corporation.Squad, error) {

	squadResChan := make(chan gamecomm.ChanResponse)
	defer close(squadResChan)
	corpCommand := gamecomm.CorpCommand{
		Action:          gamecomm.GetSquad,
		ResponseChannel: squadResChan,
		CorporationId:   corporationId,
		SquadIndex:      squadId,
	}

	gameChannels.CorpChannel <- corpCommand

	squadRes := <-squadResChan
	if squadRes.Err != nil {
		fmt.Printf("erro: %v\n", squadRes.Err.Error())
	}

	squad, ok := squadRes.Val.(corporation.Squad)
	if !ok {
		return corporation.Squad{}, fmt.Errorf("world channel returned wrong squad object: %v", squadRes.Val)
	}

	return squad, nil

}

func removeResourceFromPlanet(planetId string, resourceAmount int, resource string, gameChannels *gamecomm.GameChannels) error {
	responseChan := make(chan gamecomm.ChanResponse)
	defer close(responseChan)

	gameChannels.WorldChannel <- gamecomm.WorldCommand{
		PlanetId:        planetId,
		Action:          gamecomm.AddResourcesToPlanet,
		Amount:          resourceAmount,
		ResponseChannel: responseChan,
		Resource:        resource,
	}

	responseChanRes := <-responseChan

	err := responseChanRes.Err
	if err != nil {
		return err
	}

	return nil
}

func addResourcesToSquad(corporationId uint64, resourceAmount int, resource string, gameChannels *gamecomm.GameChannels) error {
	squadResChan := make(chan gamecomm.ChanResponse)
	defer close(squadResChan)

	gameChannels.CorpChannel <- gamecomm.CorpCommand{
		Action:          gamecomm.AddResourcesToSquad,
		ResponseChannel: squadResChan,
		CorporationId:   corporationId,
		Resource:        resource,
		Amount:          resourceAmount,
	}

	squadRes := <-squadResChan

	err := squadRes.Err
	if err != nil {
		return err
	}

	return nil
}

// - This would add resources to corporation
func returnEvent(mission *Mission, gameChannels *gamecomm.GameChannels) {
	// responseChan := make(chan any)

	mission.NotificationChan <- fmt.Sprintf("Mission Notification: Squad %v returned to base.", mission.Squads)
	for _, resource := range mission.Resources {

		removeResChan := make(chan gamecomm.ChanResponse)
		gameChannels.CorpChannel <- gamecomm.CorpCommand{
			Action:          gamecomm.RemoveResourcesFromSquad,
			ResponseChannel: removeResChan,
			CorporationId:   mission.CorporationId,
			Resource:        resource,
		}

		removedAmountRes := <-removeResChan
		if removedAmountRes.Err != nil {
			fmt.Println(removedAmountRes.Err) // TODO: Handle error
			return
		}
		removedAmount := removedAmountRes.Val.(int)

		baseResChan := make(chan gamecomm.ChanResponse)
		gameChannels.CorpChannel <- gamecomm.CorpCommand{
			Action:          gamecomm.AddResourcesToBase,
			ResponseChannel: baseResChan,
			Amount:          removedAmount,
			Resource:        resource,
			CorporationId:   mission.CorporationId,
			BaseIndex:       0,
		}

		res := <-baseResChan
		if res.Err != nil {
			fmt.Println(res.Err) // TODO: Handle error
			return
		}
		amount := res.Val.(int)

		mission.NotificationChan <- fmt.Sprintf("Mission Notification: Added to base %v -> #%v", resource, amount)
	}
}

// TODO: Change Commands channels
