package mission

import (
	"fmt"
	"testing"
	"time"

	"github.com/luisya22/galactic-exchange/gamecomm"
	"github.com/luisya22/galactic-exchange/internal/assert"
	"github.com/luisya22/galactic-exchange/world"
)

func TestArrivingEvent(t *testing.T) {
	tests := []struct {
		name    string
		mission Mission
		wants   string
	}{
		{
			name: "Receive Message",
			mission: Mission{
				Squads: []int{0},
			},
			wants: "Mission Notification: Squad [0], reached destination.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			gameChannels := &gamecomm.GameChannels{
				WorldChannel:   make(chan gamecomm.WorldCommand),
				CorpChannel:    make(chan gamecomm.CorpCommand),
				MissionChannel: make(chan gamecomm.MissionCommand),
			}

			notificationChannel := make(chan string)

			tt.mission.NotificationChan = notificationChannel

			go arrivingEvent(&tt.mission, gameChannels)

			msg := <-notificationChannel

			assert.Equal(t, msg, tt.wants)

		})
	}
}

func TestHarvestingEvent(t *testing.T) {

	type testResult struct {
		removeResourceFromPlanetCommand     gamecomm.WorldCommand
		removeResourceFromPlanetShouldError bool
		addResourceToSquad                  gamecomm.CorpCommand
		addResourceToSquadShouldError       bool
		notificationChanMsg                 string
	}

	tests := []struct {
		name                              string
		mission                           Mission
		removeResourcesFromPlanetResponse gamecomm.ChanResponse
		addResourcesToSquadResponse       gamecomm.ChanResponse
		wants                             testResult
	}{
		{
			name: "Receive all messages",
			mission: Mission{
				CorporationId: 1,
				Squads:        []int{0},
				PlanetId:      "Planet 1",
				Resources:     []string{string(world.Iron)},
			},
			removeResourcesFromPlanetResponse: gamecomm.ChanResponse{Val: 1},
			addResourcesToSquadResponse:       gamecomm.ChanResponse{Val: 1},
			wants: testResult{
				removeResourceFromPlanetCommand: gamecomm.WorldCommand{
					PlanetId: "Planet 1",
					Amount:   200,
					Resource: string(world.Iron),
					Action:   gamecomm.RemoveResourcesFromPlanet,
				},
				removeResourceFromPlanetShouldError: false,
				addResourceToSquad: gamecomm.CorpCommand{
					CorporationId: 1,
					Amount:        200,
					Resource:      string(world.Iron),
					Action:        gamecomm.AddResourcesToSquad,
				},
				addResourceToSquadShouldError: false,
				notificationChanMsg:           "Mission Notification: Squad [0], finished harvesting.",
			},
		},
		{
			name: "Remove Resources From Planet Error",
			mission: Mission{
				CorporationId: 1,
				Squads:        []int{0},
				PlanetId:      "Planet 1",
				Resources:     []string{string(world.Iron)},
			},
			removeResourcesFromPlanetResponse: gamecomm.ChanResponse{Err: fmt.Errorf("error: test error")},
			addResourcesToSquadResponse:       gamecomm.ChanResponse{Val: 1},
			wants: testResult{
				removeResourceFromPlanetCommand: gamecomm.WorldCommand{
					PlanetId: "Planet 1",
					Amount:   200,
					Resource: string(world.Iron),
					Action:   gamecomm.RemoveResourcesFromPlanet,
				},
				removeResourceFromPlanetShouldError: true,
				addResourceToSquad: gamecomm.CorpCommand{
					CorporationId: 1,
					Amount:        200,
					Resource:      string(world.Iron),
					Action:        gamecomm.AddResourcesToSquad,
				},
				addResourceToSquadShouldError: false,
				notificationChanMsg:           "Mission Notification: Squad [0], finished harvesting.",
			},
		},
		{
			name: "Add Resources To Squad Error",
			mission: Mission{
				CorporationId: 1,
				Squads:        []int{0},
				PlanetId:      "Planet 1",
				Resources:     []string{string(world.Iron)},
			},
			removeResourcesFromPlanetResponse: gamecomm.ChanResponse{Val: 1},
			addResourcesToSquadResponse:       gamecomm.ChanResponse{Err: fmt.Errorf("error: test error")},
			wants: testResult{
				removeResourceFromPlanetCommand: gamecomm.WorldCommand{
					PlanetId: "Planet 1",
					Amount:   200,
					Resource: string(world.Iron),
					Action:   gamecomm.RemoveResourcesFromPlanet,
				},
				removeResourceFromPlanetShouldError: false,
				addResourceToSquad: gamecomm.CorpCommand{
					CorporationId: 1,
					Amount:        200,
					Resource:      string(world.Iron),
					Action:        gamecomm.AddResourcesToSquad,
				},
				addResourceToSquadShouldError: true,
				notificationChanMsg:           "Mission Notification: Squad [0], finished harvesting.",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			gameChannels := &gamecomm.GameChannels{
				WorldChannel:   make(chan gamecomm.WorldCommand),
				CorpChannel:    make(chan gamecomm.CorpCommand),
				MissionChannel: make(chan gamecomm.MissionCommand),
			}

			notificationChannel := make(chan string)
			tt.mission.NotificationChan = notificationChannel

			errorChannel := make(chan error)
			tt.mission.ErrorChan = errorChannel

			go harvestingEvent(&tt.mission, gameChannels)

			// should receive remove from planet
			removeResourcePlanetCommand := <-gameChannels.WorldChannel
			assertWorldCommand(t, removeResourcePlanetCommand, tt.wants.removeResourceFromPlanetCommand)
			removeResourcePlanetCommand.ResponseChannel <- tt.removeResourcesFromPlanetResponse

			if tt.wants.removeResourceFromPlanetShouldError {
				waitForErrorOrTimeout(t, errorChannel, tt.removeResourcesFromPlanetResponse.Err)
			}

			// should receive add resources to squad
			addResourceToSquadCommand := <-gameChannels.CorpChannel
			assertCorpCommand(t, addResourceToSquadCommand, tt.wants.addResourceToSquad)
			addResourceToSquadCommand.ResponseChannel <- tt.addResourcesToSquadResponse

			if tt.wants.addResourceToSquadShouldError {
				waitForErrorOrTimeout(t, errorChannel, tt.addResourcesToSquadResponse.Err)
			}
			// should receive mission notification
			msg := <-notificationChannel
			assert.Equal(t, msg, tt.wants.notificationChanMsg)

			close(gameChannels.WorldChannel)
			close(gameChannels.CorpChannel)
			close(gameChannels.MissionChannel)
			close(notificationChannel)
			close(errorChannel)

		})
	}
}

func TestReturnEvent(t *testing.T) {
	type testResult struct {
		removeResourcesFromSquadCommand     gamecomm.CorpCommand
		removeResourcesFromSquadShouldError bool
		addResourceToBaseCommand            gamecomm.CorpCommand
		addResourceToBaseShouldError        bool
		notificationChanMsg                 string
		notificationChanMsg2                string
	}

	tests := []struct {
		name                             string
		mission                          Mission
		removeResourcesFromSquadResponse gamecomm.ChanResponse
		addResourceToBaseResponse        gamecomm.ChanResponse
		wants                            testResult
	}{
		{
			name: "Receive all messages",
			mission: Mission{
				CorporationId: 1,
				Squads:        []int{0},
				PlanetId:      "Planet 1",
				Resources:     []string{string(world.Iron)},
			},
			removeResourcesFromSquadResponse: gamecomm.ChanResponse{Val: 1},
			addResourceToBaseResponse:        gamecomm.ChanResponse{Val: 1},
			wants: testResult{
				removeResourcesFromSquadCommand: gamecomm.CorpCommand{
					CorporationId: 1,
					Resource:      string(world.Iron),
					Action:        gamecomm.RemoveResourcesFromSquad,
				},
				removeResourcesFromSquadShouldError: false,
				addResourceToBaseCommand: gamecomm.CorpCommand{
					CorporationId: 1,
					Amount:        1,
					Resource:      string(world.Iron),
					Action:        gamecomm.AddResourcesToBase,
				},
				addResourceToBaseShouldError: false,
				notificationChanMsg:          "Mission Notification: Squad [0] returned to base.",
				notificationChanMsg2:         "Mission Notification: Added to base iron -> #1",
			},
		},
		{
			name: "Remove Resource From Squad Error",
			mission: Mission{
				CorporationId: 1,
				Squads:        []int{0},
				PlanetId:      "Planet 1",
				Resources:     []string{string(world.Iron)},
			},
			removeResourcesFromSquadResponse: gamecomm.ChanResponse{Err: fmt.Errorf("error: test error")},
			addResourceToBaseResponse:        gamecomm.ChanResponse{Val: 1},
			wants: testResult{
				removeResourcesFromSquadCommand: gamecomm.CorpCommand{
					CorporationId: 1,
					Resource:      string(world.Iron),
					Action:        gamecomm.RemoveResourcesFromSquad,
				},
				removeResourcesFromSquadShouldError: true,
				addResourceToBaseCommand: gamecomm.CorpCommand{
					CorporationId: 1,
					Resource:      string(world.Iron),
					Action:        gamecomm.AddResourcesToBase,
				},
				addResourceToBaseShouldError: false,
				notificationChanMsg:          "Mission Notification: Squad [0] returned to base.",
				notificationChanMsg2:         "Mission Notification: Added to base iron -> #1",
			},
		},
		{
			name: "Add Error To Base Error",
			mission: Mission{
				CorporationId: 1,
				Squads:        []int{0},
				PlanetId:      "Planet 1",
				Resources:     []string{string(world.Iron)},
			},
			removeResourcesFromSquadResponse: gamecomm.ChanResponse{Val: 1},
			addResourceToBaseResponse:        gamecomm.ChanResponse{Err: fmt.Errorf("error: test error")},
			wants: testResult{
				removeResourcesFromSquadCommand: gamecomm.CorpCommand{
					CorporationId: 1,
					Amount:        0,
					Resource:      string(world.Iron),
					Action:        gamecomm.RemoveResourcesFromSquad,
				},
				removeResourcesFromSquadShouldError: false,
				addResourceToBaseCommand: gamecomm.CorpCommand{
					CorporationId: 1,
					Amount:        1,
					Resource:      string(world.Iron),
					Action:        gamecomm.AddResourcesToBase,
				},
				addResourceToBaseShouldError: true,
				notificationChanMsg:          "Mission Notification: Squad [0] returned to base.",
				notificationChanMsg2:         "Mission Notification: Added to base iron -> #1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			gameChannels := &gamecomm.GameChannels{
				WorldChannel:   make(chan gamecomm.WorldCommand),
				CorpChannel:    make(chan gamecomm.CorpCommand),
				MissionChannel: make(chan gamecomm.MissionCommand),
			}

			notificationChannel := make(chan string)
			tt.mission.NotificationChan = notificationChannel

			errorChannel := make(chan error)
			tt.mission.ErrorChan = errorChannel

			go returnEvent(&tt.mission, gameChannels)

			// should receive notification
			msg := <-notificationChannel
			assert.Equal(t, msg, tt.wants.notificationChanMsg)

			// should receive removeResourcesFromSquad
			removeResourcesFromSquadCommand := <-gameChannels.CorpChannel
			assertCorpCommand(t, removeResourcesFromSquadCommand, tt.wants.removeResourcesFromSquadCommand)
			removeResourcesFromSquadCommand.ResponseChannel <- tt.removeResourcesFromSquadResponse

			if tt.wants.removeResourcesFromSquadShouldError {
				waitForErrorOrTimeout(t, errorChannel, tt.removeResourcesFromSquadResponse.Err)
			}

			// should receive addResourcesToBase
			addResourcesToBaseCommand := <-gameChannels.CorpChannel
			assertCorpCommand(t, addResourcesToBaseCommand, tt.wants.addResourceToBaseCommand)
			addResourcesToBaseCommand.ResponseChannel <- tt.addResourceToBaseResponse

			if tt.wants.addResourceToBaseShouldError {
				waitForErrorOrTimeout(t, errorChannel, tt.addResourceToBaseResponse.Err)
			}

			// should receive notification
			msg = <-notificationChannel
			assert.Equal(t, msg, tt.wants.notificationChanMsg2)

			close(gameChannels.WorldChannel)
			close(gameChannels.CorpChannel)
			close(gameChannels.MissionChannel)
			close(notificationChannel)
			close(errorChannel)

		})
	}

}

func assertCorpCommand(t *testing.T, got gamecomm.CorpCommand, wants gamecomm.CorpCommand) {
	assert.Equal(t, got.CorporationId, wants.CorporationId)
	assert.Equal(t, got.Action, wants.Action)
	assert.Equal(t, got.Amount, wants.Amount)
	assert.Equal(t, got.Resource, wants.Resource)
	assert.Equal(t, got.BaseIndex, wants.BaseIndex)
	assert.Equal(t, got.SquadIndex, wants.SquadIndex)
	assert.Equal(t, got.AmountDecimal, wants.AmountDecimal)
}

func assertWorldCommand(t *testing.T, got gamecomm.WorldCommand, wants gamecomm.WorldCommand) {
	assert.Equal(t, got.Resource, wants.Resource)
	assert.Equal(t, got.Amount, wants.Amount)
	assert.Equal(t, got.PlanetId, wants.PlanetId)
	assert.Equal(t, got.Action, wants.Action)
}

func waitForErrorOrTimeout(t *testing.T, errorChannel chan error, resErr error) {
	timer := time.NewTimer(5 * time.Second)

	select {
	case err := <-errorChannel:
		assert.Error(t, err)
		assert.Equal(t, err, resErr)
	case <-timer.C:
		t.Errorf("did not receive error on addResourceToSquad")
	}
}
