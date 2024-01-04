package mission

import (
	"fmt"

	"github.com/luisya22/galactic-exchange/gameclock"
	"github.com/luisya22/galactic-exchange/gamecomm"
	"github.com/luisya22/galactic-exchange/world"
)

// TODO: Need to update squad positions for every event

// amount int, itemName world.Resource, planetId string, corporationId uint64
func (ms *MissionScheduler) CreateTransferMission(m Mission) error {

	ms.CalculateTravelDistance(m)

	// Leaving Event
	// Rmove resources from corporation
	// Add resources to squad

	le := &Event{
		MissionId: m.Id,
		Time:      ms.gameClock.GetCurrentTime(),
		Cancelled: false,
		Execute:   tsLeavingEvent,
	}

	ms.eventScheduler.Schedule(le)

	ae := &Event{
		MissionId: m.Id,
		Time:      ms.gameClock.GetCurrentTime().Add(3 * gameclock.Day),
		Cancelled: false,
		Execute:   tsArrivalEvent,
	}

	ms.eventScheduler.Schedule(ae)

	// Returnal Event
	// Return squad to base

	bb := &Event{
		MissionId: m.Id,
		Time:      ms.gameClock.GetCurrentTime().Add(6 * gameclock.Day),
		Cancelled: false,
		Execute:   tsBackToBase,
	}

	ms.eventScheduler.Schedule(bb)

	return nil

}

func tsLeavingEvent(mission *Mission, gameChannels *gamecomm.GameChannels) {

	for _, resource := range mission.Resources {

		removedAmount, err := removeResourcesFromCorporation(mission.CorporationId, mission.Amount, resource, gameChannels)
		if err != nil {
			fmt.Println(err.Error()) // TODO: Handle error correctly
			return
		}

		err = addResourcesToSquad(mission.CorporationId, removedAmount, resource, gameChannels)
		if err != nil {
			fmt.Println(err.Error()) // TODO: Handle error correctly
			return
		}

	}

	mission.NotificationChan <- fmt.Sprintf("Mission Notification: Squad %v, started travel.", mission.Squads)
}

func tsArrivalEvent(mission *Mission, gameChannels *gamecomm.GameChannels) {

	sumCredits := 0.0
	for _, resource := range mission.Resources {
		removedAmount, err := removeResourcesFromSquad(mission.CorporationId, resource, gameChannels)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		err = addResourcesToPlanet(mission.PlanetId, removedAmount, resource, gameChannels)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		// TODO fix prices, build a economy module that get the actual price. It would depend on various things (current contract between base and corporation, zone prices, base item price, sanctions)
		credits := float64(removedAmount * 2)
		err = addCreditsToCorporation(mission.CorporationId, credits, gameChannels)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		sumCredits += credits
	}

	mission.NotificationChan <- fmt.Sprintf("Mission Notification: Squad %v, made the delivery. Added Credits: $%v", mission.Squads[0], sumCredits)
}

func tsBackToBase(mission *Mission, gameChannels *gamecomm.GameChannels) {
	// TODO: In the future make the squad available again
	mission.NotificationChan <- fmt.Sprintf("Mission Notification: Squad %v is back to base", mission.Squads[0])
}

func (ms *MissionScheduler) CalculateTravelDistance(m Mission) {
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

}
