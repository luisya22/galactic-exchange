package game

import (
	"github.com/luisya22/galactic-exchange/internal/gamecomm"
)

// TODO: Need to add inbound base when Base system is implemented
func (g *Game) HarvestPlanet(planetId string, corporationId uint64, squadId int, notificationChan chan string) error {

	mc := gamecomm.MissionCommand{
		CorporationId:    corporationId,
		Squads:           []int{squadId},
		Type:             gamecomm.SquadMission,
		Resources:        []string{string("iron")},
		NotificationChan: notificationChan,
		PlanetId:         planetId,
	}

	g.gameChannels.MissionChannel <- mc
	return nil
}
