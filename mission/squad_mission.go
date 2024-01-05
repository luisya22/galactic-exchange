package mission

import (
	"fmt"

	"github.com/luisya22/galactic-exchange/gameclock"
	"github.com/luisya22/galactic-exchange/gamecomm"
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

	planet := planetRes.Val.(gamecomm.Planet)
	close(planetResChan)

	// CALCULATE SHIP SPEED
	shipSpeed := squad.Ships.Speed
	squadLocation := gamecomm.Coordinates{X: squad.Location.X, Y: squad.Location.Y}
	planetLocation := gamecomm.Coordinates{X: planet.Location.X, Y: planet.Location.Y}

	planetDistance := gamecomm.Distance(squadLocation, planetLocation)
	_ = planetDistance / float64(shipSpeed)

	// CREATE ARRIVE EVENT
	ae := Event{
		MissionId: m.Id,
		Time:      ms.gameClock.GetCurrentTime(),
		Cancelled: false,
		Execute:   arrivingEvent,
	}

	ms.eventScheduler.Schedule(&ae)

	// CREATE HARVESTING RESOURCES EVENT
	he := Event{
		MissionId: m.Id,
		Time:      ms.gameClock.GetCurrentTime().Add(2 * gameclock.Day),
		Cancelled: false,
		Execute:   harvestingEvent,
	}

	ms.eventScheduler.Schedule(&he)

	// CREATE RETURN EVENT
	re := Event{
		MissionId: m.Id,
		Time:      ms.gameClock.GetCurrentTime().Add(3 * gameclock.Day),
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

// - This would add resources to corporation
func returnEvent(mission *Mission, gameChannels *gamecomm.GameChannels) {
	// responseChan := make(chan any)

	mission.NotificationChan <- fmt.Sprintf("Mission Notification: Squad %v returned to base.", mission.Squads)
	for _, resource := range mission.Resources {

		removedAmount, err := removeResourcesFromSquad(mission.CorporationId, resource, gameChannels)
		if err != nil {
			fmt.Println(err.Error()) // TODO: Handle Error correctly
			return
		}

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
