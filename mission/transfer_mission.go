package mission

import (
	"fmt"

	"github.com/luisya22/galactic-exchange/gameclock"
	"github.com/luisya22/galactic-exchange/gamecomm"
)

// TODO: Need to update squad positions for every event

// amount int, itemName world.Resource, planetId string, corporationId uint64
func (ms *MissionScheduler) CreateTransferMission(m Mission) error {

	_, err := ms.CalculateTravelDistance(m.CorporationId, m.Squads, m.PlanetId, ms.GameChannels)
	if err != nil {
		return err
	}

	// Leaving Event
	// Rmove resources from corporation
	// Add resources to squad

	le := &Event{
		MissionId: m.Id,
		Time:      ms.GameClock.GetCurrentTime(),
		Cancelled: false,
		Execute:   tsLeavingEvent,
	}

	leId, err := ms.EventScheduler.Schedule(le)
	if err != nil {
		return err
	}

	ae := &Event{
		MissionId: m.Id,
		Time:      ms.GameClock.GetCurrentTime().Add(3 * gameclock.Day),
		Cancelled: false,
		Execute:   tsArrivalEvent,
	}

	aeId, err := ms.EventScheduler.Schedule(ae)
	if err != nil {
		updateErr := ms.EventScheduler.UpdateEvent(leId, le.Time, true)
		if updateErr != nil {
			return updateErr
		}

		return err
	}

	// Returnal Event
	// Return squad to base

	bb := &Event{
		MissionId: m.Id,
		Time:      ms.GameClock.GetCurrentTime().Add(6 * gameclock.Day),
		Cancelled: false,
		Execute:   tsBackToBase,
	}

	_, err = ms.EventScheduler.Schedule(bb)
	if err != nil {
		updateErr := ms.EventScheduler.UpdateEvent(aeId, ae.Time, true)
		if updateErr != nil {
			return updateErr
		}

		updateErr = ms.EventScheduler.UpdateEvent(leId, le.Time, true)
		if updateErr != nil {
			return updateErr
		}
		return err
	}

	return nil

}

func tsLeavingEvent(mission *Mission, gameChannels *gamecomm.GameChannels) {

	for _, resource := range mission.Resources {

		removedAmount, err := removeResourcesFromCorporation(mission.CorporationId, mission.Amount, resource, gameChannels)
		if err != nil {
			mission.ErrorChan <- err
			removedAmount = mission.Amount
		}

		err = addResourcesToSquad(mission.CorporationId, removedAmount, resource, gameChannels)
		if err != nil {
			mission.ErrorChan <- err
		}

	}

	mission.NotificationChan <- fmt.Sprintf("Mission Notification: Squad %v, started travel.", mission.Squads)
}

func tsArrivalEvent(mission *Mission, gameChannels *gamecomm.GameChannels) {

	sumCredits := 0.0
	for _, resource := range mission.Resources {
		removedAmount, err := removeAllResourcesFromSquad(mission.CorporationId, resource, gameChannels)
		if err != nil {
			removedAmount = mission.Amount
			mission.ErrorChan <- err
		}

		err = addResourcesToPlanet(mission.PlanetId, removedAmount, resource, gameChannels)
		if err != nil {
			mission.ErrorChan <- err
		}

		// TODO fix prices, build a economy module that get the actual price. It would depend on various things (current contract between base and corporation, zone prices, base item price, sanctions)
		credits := float64(removedAmount * 2)
		err = addCreditsToCorporation(mission.CorporationId, credits, gameChannels)
		if err != nil {
			mission.ErrorChan <- err
			continue
		}

		sumCredits += credits
	}

	mission.NotificationChan <- fmt.Sprintf("Mission Notification: Squad %v, made the delivery. Added Credits: $%v", mission.Squads[0], sumCredits)
}

func tsBackToBase(mission *Mission, gameChannels *gamecomm.GameChannels) {
	// TODO: In the future make the squad available again
	mission.NotificationChan <- fmt.Sprintf("Mission Notification: Squad %v is back to base", mission.Squads[0])
}
