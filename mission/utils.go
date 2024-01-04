package mission

import (
	"fmt"

	"github.com/luisya22/galactic-exchange/corporation"
	"github.com/luisya22/galactic-exchange/gamecomm"
)

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

func removeResourcesFromSquad(corporationId uint64, resource string, gameChannels *gamecomm.GameChannels) (int, error) {
	removeResChan := make(chan gamecomm.ChanResponse)
	gameChannels.CorpChannel <- gamecomm.CorpCommand{
		Action:          gamecomm.RemoveResourcesFromSquad,
		ResponseChannel: removeResChan,
		CorporationId:   corporationId,
		Resource:        resource,
	}

	removedAmountRes := <-removeResChan
	if removedAmountRes.Err != nil {
		return 0, removedAmountRes.Err
	}

	return removedAmountRes.Val.(int), nil
}

func removeResourcesFromCorporation(corporationId uint64, amount int, resource string, gameChannels *gamecomm.GameChannels) (int, error) {
	removeResChan := make(chan gamecomm.ChanResponse)
	gameChannels.CorpChannel <- gamecomm.CorpCommand{
		Action:          gamecomm.RemoveResourcesFromBase,
		ResponseChannel: removeResChan,
		CorporationId:   corporationId,
		Resource:        resource,
		Amount:          amount,
	}

	removedAmountRes := <-removeResChan
	if removedAmountRes.Err != nil {
		return 0, removedAmountRes.Err
	}

	return removedAmountRes.Val.(int), nil
}

func addResourcesToPlanet(planetId string, resourceAmount int, resource string, gameChannels *gamecomm.GameChannels) error {

	addResChan := make(chan gamecomm.ChanResponse)
	gameChannels.WorldChannel <- gamecomm.WorldCommand{
		Action:          gamecomm.AddResourcesToPlanet,
		ResponseChannel: addResChan,
		PlanetId:        planetId,
		Resource:        resource,
	}

	addRes := <-addResChan
	if addRes.Err != nil {
		return addRes.Err
	}

	return nil
}

func addCreditsToCorporation(corporationId uint64, amount float64, gameChannels *gamecomm.GameChannels) error {
	resChan := make(chan gamecomm.ChanResponse)
	gameChannels.CorpChannel <- gamecomm.CorpCommand{
		Action:          gamecomm.AddCredits,
		ResponseChannel: resChan,
		CorporationId:   corporationId,
		AmountDecimal:   amount,
	}

	res := <-resChan
	if res.Err != nil {
		return res.Err
	}

	return nil
}

func removeCreditsFromCorporation(corporationId uint64, amount float64, gameChannels *gamecomm.GameChannels) error {
	resChan := make(chan gamecomm.ChanResponse)
	gameChannels.CorpChannel <- gamecomm.CorpCommand{
		Action:          gamecomm.AddCredits,
		ResponseChannel: resChan,
		CorporationId:   corporationId,
		AmountDecimal:   amount,
	}

	res := <-resChan
	if res.Err != nil {
		return res.Err
	}

	return nil
}
