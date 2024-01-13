package mission

import (
	"fmt"

	"github.com/luisya22/galactic-exchange/internal/gameclock"
	"github.com/luisya22/galactic-exchange/internal/gamecomm"
)

// TODO: let mission scheduler that a mission is completed so it can erase it
func (ms *MissionScheduler) CreateSquadMission(m Mission) error {
	// TODO: Random events could affect mission times

	// TODO: What happen if squads are on different positions?

	if len(m.Squads) == 0 {
		return fmt.Errorf("error: should include squads")
	}

	// GET SQUAD
	squadId := m.Squads[0]
	squad, err := getSquad(m.CorporationId, squadId, ms.GameChannels)
	if err != nil {
		return err
	}

	// GET PLANET
	planetResChan := make(chan gamecomm.ChanResponse)
	planetCommand := gamecomm.WorldCommand{
		PlanetId:        m.PlanetId,
		Action:          gamecomm.GetPlanet,
		ResponseChannel: planetResChan,
	}
	ms.GameChannels.WorldChannel <- planetCommand

	planetRes := <-planetResChan
	if planetRes.Err != nil {
		return planetRes.Err
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
		Time:      ms.GameClock.GetCurrentTime(),
		Cancelled: false,
		Execute:   arrivingEvent,
	}

	aeId, err := ms.EventScheduler.Schedule(&ae)
	if err != nil {
		return err
	}

	// CREATE HARVESTING RESOURCES EVENT
	he := Event{
		MissionId: m.Id,
		Time:      ms.GameClock.GetCurrentTime().Add(2 * gameclock.Day),
		Cancelled: false,
		Execute:   harvestingEvent,
	}

	heId, err := ms.EventScheduler.Schedule(&he)
	if err != nil {
		updateErr := ms.EventScheduler.UpdateEvent(aeId, ae.Time, true)
		if updateErr != nil {
			return updateErr
		}

		return err
	}

	// CREATE RETURN EVENT
	re := Event{
		MissionId: m.Id,
		Time:      ms.GameClock.GetCurrentTime().Add(3 * gameclock.Day),
		Cancelled: false,
		Execute:   returnEvent,
	}

	_, err = ms.EventScheduler.Schedule(&re)
	if err != nil {
		updateErr := ms.EventScheduler.UpdateEvent(heId, re.Time, true)
		if updateErr != nil {
			return updateErr
		}

		updateErr = ms.EventScheduler.UpdateEvent(aeId, ae.Time, true)
		if updateErr != nil {
			return updateErr
		}

		return err
	}

	// TODO: If error to schedule it should cancel previous events

	return nil
}

// - This would send message that we arrive to the mission place
func arrivingEvent(mission *Mission, gameChannels *gamecomm.GameChannels) {
	mission.NotificationChan <- fmt.Sprintf("Mission Notification: Squad %v, reached destination.", mission.Squads)
	// TODO: Calculate danger
}

// - This would Gather the resources
func harvestingEvent(mission *Mission, gameChannels *gamecomm.GameChannels) {

	// Get Squad
	// squad, err := getSquad(mission.CorporationId, mission.Squads[0], gameChannels)
	// if err != nil {
	// 	fmt.Println(err.Error()) //TODO: Handle error
	// }

	for _, resource := range mission.Resources {
		// Generate harvested resourcesAmount

		// TODO: check this later
		// bonus := squad.GetHarvestingBonus()
		bonus := 2
		resourceAmount := 100 * bonus

		// Remove Resources from planet
		err := removeResourceFromPlanet(mission.PlanetId, resourceAmount, resource, gameChannels)
		if err != nil {
			mission.ErrorChan <- err
		}

		// Add Resources to Squad
		err = addResourcesToSquad(mission.CorporationId, resourceAmount, resource, gameChannels)
		if err != nil {
			mission.ErrorChan <- err
		}
	}

	mission.NotificationChan <- fmt.Sprintf("Mission Notification: Squad %v, finished harvesting.", mission.Squads)

}

// - This would add resources to corporation
func returnEvent(mission *Mission, gameChannels *gamecomm.GameChannels) {
	// responseChan := make(chan any)

	mission.NotificationChan <- fmt.Sprintf("Mission Notification: Squad %v returned to base.", mission.Squads)
	for _, resource := range mission.Resources {
		var amount int

		removedAmount, err := removeAllResourcesFromSquad(mission.CorporationId, resource, gameChannels)
		if err != nil {
			mission.ErrorChan <- err
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
			mission.ErrorChan <- res.Err
			amount = removedAmount
		} else {
			amount = res.Val.(int)
		}

		mission.NotificationChan <- fmt.Sprintf("Mission Notification: Added to base %v -> #%v", resource, amount)
	}
}

// TODO: Change Commands channels
