package game

// This need to run on a Background routine, maybe a Mission simulation routine. This cand send the mission via channel?
// TODO: Need to add inbound base when Base system is implemented
func (g *Game) HarvestPlanet(planetId string, corporationId uint64, squadId int) error {
	// Generate Harvesting Results: ships damage, resources recollected, crew lost
	// Outbound Journey Simulation
	// Remove Resources from planet
	// Inbound Journey Simulation
	// Add Resources to Corporation Base
	// Send Player Mission Report
	return nil
}
