package mission

import (
	"fmt"
	"testing"

	"github.com/luisya22/galactic-exchange/internal/assert"
	"github.com/luisya22/galactic-exchange/internal/gamecomm"
)

func TestLeavingEvent(t *testing.T) {
	type testResult struct {
		removeResourcesFromCorporationCommand     gamecomm.CorpCommand
		removeResourcesFromCorporationShouldError bool
		addResourcesToSquadCommand                gamecomm.CorpCommand
		addResourcesToSquadShouldError            bool
		notificationChanMsg                       string
	}

	tests := []struct {
		name                                  string
		mission                               Mission
		removeResourceFromCorporationResponse gamecomm.ChanResponse
		addResourcesToSquadResponse           gamecomm.ChanResponse
		wants                                 testResult
	}{
		{
			name: "Receive all messages",
			mission: Mission{
				CorporationId: 1,
				Squads:        []int{0},
				Resources:     []string{"iron"},
				Amount:        1,
			},
			removeResourceFromCorporationResponse: gamecomm.ChanResponse{Val: 1},
			addResourcesToSquadResponse:           gamecomm.ChanResponse{Val: 1},
			wants: testResult{
				removeResourcesFromCorporationCommand: gamecomm.CorpCommand{
					CorporationId: 1,
					Amount:        1,
					Resource:      "iron",
					Action:        gamecomm.RemoveResourcesFromBase,
				},
				removeResourcesFromCorporationShouldError: false,
				addResourcesToSquadCommand: gamecomm.CorpCommand{
					CorporationId: 1,
					Amount:        1,
					Resource:      "iron",
					Action:        gamecomm.AddResourcesToSquad,
				},
				addResourcesToSquadShouldError: false,
				notificationChanMsg:            "Mission Notification: Squad [0], started travel.",
			},
		},
		{
			name: "Remove Resources From Corporation Error",
			mission: Mission{
				CorporationId: 1,
				Squads:        []int{0},
				Resources:     []string{"iron"},
				Amount:        1,
			},
			removeResourceFromCorporationResponse: gamecomm.ChanResponse{Err: fmt.Errorf("error: test error")},
			addResourcesToSquadResponse:           gamecomm.ChanResponse{Val: 1},
			wants: testResult{
				removeResourcesFromCorporationCommand: gamecomm.CorpCommand{
					CorporationId: 1,
					Amount:        1,
					Resource:      "iron",
					Action:        gamecomm.RemoveResourcesFromBase,
				},
				removeResourcesFromCorporationShouldError: true,
				addResourcesToSquadCommand: gamecomm.CorpCommand{
					CorporationId: 1,
					Amount:        1,
					Resource:      "iron",
					Action:        gamecomm.AddResourcesToSquad,
				},
				addResourcesToSquadShouldError: false,
				notificationChanMsg:            "Mission Notification: Squad [0], started travel.",
			},
		},
		{
			name: "Add Resources to Squad Error",
			mission: Mission{
				CorporationId: 1,
				Squads:        []int{0},
				Resources:     []string{"iron"},
				Amount:        1,
			},
			removeResourceFromCorporationResponse: gamecomm.ChanResponse{Val: 1},
			addResourcesToSquadResponse:           gamecomm.ChanResponse{Err: fmt.Errorf("error: test error")},
			wants: testResult{
				removeResourcesFromCorporationCommand: gamecomm.CorpCommand{
					CorporationId: 1,
					Amount:        1,
					Resource:      "iron",
					Action:        gamecomm.RemoveResourcesFromBase,
				},
				removeResourcesFromCorporationShouldError: false,
				addResourcesToSquadCommand: gamecomm.CorpCommand{
					CorporationId: 1,
					Amount:        1,
					Resource:      "iron",
					Action:        gamecomm.AddResourcesToSquad,
				},
				addResourcesToSquadShouldError: true,
				notificationChanMsg:            "Mission Notification: Squad [0], started travel.",
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

			go tsLeavingEvent(&tt.mission, gameChannels)

			// should receive remove from corporation
			removeResourceFromCorporationCommand := <-gameChannels.CorpChannel
			assertCorpCommand(t, removeResourceFromCorporationCommand, tt.wants.removeResourcesFromCorporationCommand)
			removeResourceFromCorporationCommand.ResponseChannel <- tt.removeResourceFromCorporationResponse

			if tt.wants.removeResourcesFromCorporationShouldError {
				waitForErrorOrTimeout(t, errorChannel, tt.removeResourceFromCorporationResponse.Err)
			}

			// should receive add resources to squad
			addResourcesToSquadCommand := <-gameChannels.CorpChannel
			assertCorpCommand(t, addResourcesToSquadCommand, tt.wants.addResourcesToSquadCommand)
			addResourcesToSquadCommand.ResponseChannel <- tt.addResourcesToSquadResponse

			if tt.wants.addResourcesToSquadShouldError {
				waitForErrorOrTimeout(t, errorChannel, tt.addResourcesToSquadResponse.Err)
			}

			// should receive mission notification
			msg := <-notificationChannel
			assert.Equal(t, msg, tt.wants.notificationChanMsg)

			close(gameChannels.CorpChannel)
			close(gameChannels.WorldChannel)
			close(gameChannels.MissionChannel)
			close(notificationChannel)
			close(errorChannel)

		})
	}

}

func TestArrivalEvent(t *testing.T) {

	type testResult struct {
		removeAllResourcesFromSquadCommand     gamecomm.CorpCommand
		removeAllResourcesFromSquadShouldError bool
		addResourcesToPlanetCommand            gamecomm.WorldCommand
		addResourcesToPlanetShouldError        bool
		addCreditsToCorporationCommand         gamecomm.CorpCommand
		addCreditsToCorporationShouldError     bool
		notificationChanMsg                    string
	}

	tests := []struct {
		name                                string
		mission                             Mission
		removeAllResourcesFromSquadResponse gamecomm.ChanResponse
		addResourcesToPlanetResponse        gamecomm.ChanResponse
		addCreditsToCorporationResponse     gamecomm.ChanResponse
		wants                               testResult
	}{
		{
			name: "Receive all messages",
			mission: Mission{
				CorporationId: 1,
				Squads:        []int{0},
				Resources:     []string{"iron"},
				Amount:        1,
				PlanetId:      "Planet-1",
			},
			removeAllResourcesFromSquadResponse: gamecomm.ChanResponse{Val: 1},
			addResourcesToPlanetResponse:        gamecomm.ChanResponse{Val: 1},
			addCreditsToCorporationResponse:     gamecomm.ChanResponse{Val: 2},
			wants: testResult{
				removeAllResourcesFromSquadCommand: gamecomm.CorpCommand{
					CorporationId: 1,
					Resource:      "iron",
					Action:        gamecomm.RemoveAllResourcesFromSquad,
				},
				removeAllResourcesFromSquadShouldError: false,
				addResourcesToPlanetCommand: gamecomm.WorldCommand{
					PlanetId: "Planet-1",
					Resource: "iron",
					Amount:   1,
					Action:   gamecomm.AddResourcesToPlanet,
				},
				addResourcesToPlanetShouldError: false,
				addCreditsToCorporationCommand: gamecomm.CorpCommand{
					CorporationId: 1,
					AmountDecimal: 2,
					Action:        gamecomm.AddCredits,
				},
				addCreditsToCorporationShouldError: false,
				notificationChanMsg:                "Mission Notification: Squad 0, made the delivery. Added Credits: $2",
			},
		},
		{
			name: "Remove Resources From Squad Error",
			mission: Mission{
				CorporationId: 1,
				Squads:        []int{0},
				Resources:     []string{"iron"},
				Amount:        2,
				PlanetId:      "Planet-1",
			},
			removeAllResourcesFromSquadResponse: gamecomm.ChanResponse{Err: fmt.Errorf("error: test error")},
			addResourcesToPlanetResponse:        gamecomm.ChanResponse{Val: 1},
			addCreditsToCorporationResponse:     gamecomm.ChanResponse{Val: 2},
			wants: testResult{
				removeAllResourcesFromSquadCommand: gamecomm.CorpCommand{
					CorporationId: 1,
					Resource:      "iron",
					Action:        gamecomm.RemoveAllResourcesFromSquad,
				},
				removeAllResourcesFromSquadShouldError: true,
				addResourcesToPlanetCommand: gamecomm.WorldCommand{
					PlanetId: "Planet-1",
					Resource: "iron",
					Amount:   2,
					Action:   gamecomm.AddResourcesToPlanet,
				},
				addResourcesToPlanetShouldError: false,
				addCreditsToCorporationCommand: gamecomm.CorpCommand{
					CorporationId: 1,
					AmountDecimal: 4,
					Action:        gamecomm.AddCredits,
				},
				addCreditsToCorporationShouldError: false,
				notificationChanMsg:                "Mission Notification: Squad 0, made the delivery. Added Credits: $4",
			},
		},
		{
			name: "Add Resources to Planet Error",
			mission: Mission{
				CorporationId: 1,
				Squads:        []int{0},
				Resources:     []string{"iron"},
				Amount:        1,
				PlanetId:      "Planet-1",
			},
			removeAllResourcesFromSquadResponse: gamecomm.ChanResponse{Val: 2},
			addResourcesToPlanetResponse:        gamecomm.ChanResponse{Err: fmt.Errorf("error: test error")},
			addCreditsToCorporationResponse:     gamecomm.ChanResponse{Val: 2},
			wants: testResult{
				removeAllResourcesFromSquadCommand: gamecomm.CorpCommand{
					CorporationId: 1,
					Resource:      "iron",
					Action:        gamecomm.RemoveAllResourcesFromSquad,
				},
				removeAllResourcesFromSquadShouldError: false,
				addResourcesToPlanetCommand: gamecomm.WorldCommand{
					PlanetId: "Planet-1",
					Resource: "iron",
					Amount:   2,
					Action:   gamecomm.AddResourcesToPlanet,
				},
				addResourcesToPlanetShouldError: true,
				addCreditsToCorporationCommand: gamecomm.CorpCommand{
					CorporationId: 1,
					AmountDecimal: 4,
					Action:        gamecomm.AddCredits,
				},
				addCreditsToCorporationShouldError: false,
				notificationChanMsg:                "Mission Notification: Squad 0, made the delivery. Added Credits: $4",
			},
		},
		{
			name: "Add Credits to Corporation Error",
			mission: Mission{
				CorporationId: 1,
				Squads:        []int{0},
				Resources:     []string{"iron"},
				Amount:        1,
				PlanetId:      "Planet-1",
			},
			removeAllResourcesFromSquadResponse: gamecomm.ChanResponse{Val: 1},
			addResourcesToPlanetResponse:        gamecomm.ChanResponse{Val: 1},
			addCreditsToCorporationResponse:     gamecomm.ChanResponse{Err: fmt.Errorf("error: test error")},
			wants: testResult{
				removeAllResourcesFromSquadCommand: gamecomm.CorpCommand{
					CorporationId: 1,
					Resource:      "iron",
					Action:        gamecomm.RemoveAllResourcesFromSquad,
				},
				removeAllResourcesFromSquadShouldError: false,
				addResourcesToPlanetCommand: gamecomm.WorldCommand{
					PlanetId: "Planet-1",
					Resource: "iron",
					Amount:   1,
					Action:   gamecomm.AddResourcesToPlanet,
				},
				addResourcesToPlanetShouldError: false,
				addCreditsToCorporationCommand: gamecomm.CorpCommand{
					CorporationId: 1,
					AmountDecimal: 2,
					Action:        gamecomm.AddCredits,
				},
				addCreditsToCorporationShouldError: true,
				notificationChanMsg:                "Mission Notification: Squad 0, made the delivery. Added Credits: $0",
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

			go tsArrivalEvent(&tt.mission, gameChannels)

			// should receive remove resources from squad
			removeResourcesFromSquadCommand := <-gameChannels.CorpChannel
			assertCorpCommand(t, removeResourcesFromSquadCommand, tt.wants.removeAllResourcesFromSquadCommand)
			removeResourcesFromSquadCommand.ResponseChannel <- tt.removeAllResourcesFromSquadResponse

			if tt.wants.removeAllResourcesFromSquadShouldError {
				waitForErrorOrTimeout(t, errorChannel, tt.removeAllResourcesFromSquadResponse.Err)
			}

			// should receive add resources to planet
			addResourcesToPlanetCommand := <-gameChannels.WorldChannel
			assertWorldCommand(t, addResourcesToPlanetCommand, tt.wants.addResourcesToPlanetCommand)
			addResourcesToPlanetCommand.ResponseChannel <- tt.addResourcesToPlanetResponse

			if tt.wants.addResourcesToPlanetShouldError {
				waitForErrorOrTimeout(t, errorChannel, tt.addResourcesToPlanetResponse.Err)
			}

			// should receive add credits to corporation
			addCreditsToCorporationCommand := <-gameChannels.CorpChannel
			assertCorpCommand(t, addCreditsToCorporationCommand, tt.wants.addCreditsToCorporationCommand)
			addCreditsToCorporationCommand.ResponseChannel <- tt.addCreditsToCorporationResponse

			if tt.wants.addCreditsToCorporationShouldError {
				waitForErrorOrTimeout(t, errorChannel, tt.addResourcesToPlanetResponse.Err)
			}

			// should receive mission notification
			msg := <-notificationChannel
			assert.Equal(t, msg, tt.wants.notificationChanMsg)

			close(gameChannels.CorpChannel)
			close(gameChannels.WorldChannel)
			close(gameChannels.MissionChannel)
			close(notificationChannel)
			close(errorChannel)

		})
	}
}

func TestBackToBase(t *testing.T) {
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
			wants: "Mission Notification: Squad 0 is back to base",
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

			go tsBackToBase(&tt.mission, gameChannels)

			msg := <-notificationChannel

			assert.Equal(t, msg, tt.wants)

		})
	}
}
