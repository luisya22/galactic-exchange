package game

import (
	"github.com/luisya22/galactic-exchange/gamecomm"
	"github.com/luisya22/galactic-exchange/world"
)

// This need to run on a Background routine, maybe a Mission simulation routine. This cand send the mission via channel?
// TODO: Need to add inbound base when Base system is implemented
func (g *Game) HarvestPlanet(planetId string, corporationId uint64, squadId int) error {

	mc := gamecomm.MissionCommand{
		CorporationId: corporationId,
		Squads:        []int{squadId},
		Type:          gamecomm.SquadMission,
		Resources:     []string{string(world.Iron)},
	}

	g.gameChannels.MissionChannel <- mc
	return nil
}
