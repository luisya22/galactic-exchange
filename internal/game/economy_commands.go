package game

import (
	"github.com/luisya22/galactic-exchange/internal/gamecomm"
)

// TODO: Amounts should be reflected after the time distance is elapsed
// TODO: Use MissionScheduler
// func (g *Game) SellResource(amount int, itemName world.Resource, planetId string, corporationId uint64) error {
func (g *Game) SellResource(planetId string, corporationId uint64, squadId int, amount int, itemName string, notificationChan chan string) error {
	mc := gamecomm.MissionCommand{
		CorporationId:    corporationId,
		Squads:           []int{squadId},
		Type:             gamecomm.TransferMission,
		Resources:        []string{"iron"},
		NotificationChan: notificationChan,
		PlanetId:         planetId,
		Amount:           amount,
	}

	g.gameChannels.MissionChannel <- mc
	return nil
}
