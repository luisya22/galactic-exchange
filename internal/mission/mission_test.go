package mission_test

import (
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/luisya22/galactic-exchange/internal/assert"
	"github.com/luisya22/galactic-exchange/internal/gameclock"
	"github.com/luisya22/galactic-exchange/internal/gamecomm"
	"github.com/luisya22/galactic-exchange/internal/mission"
)

func TestStartSquadMission(t *testing.T) {
	// Create Mission
	// Pass mock eventSchduler
	// Check that mission is added to the missions map
	// Check that the correct type of mission is created
	// Ensure that Schedule is invoked
	// Check CreateMission error
	// Check that notificationChan receive correct message
	// Test for race conditions
	// Check that missions have the correct parameters
	// Empty or malformed mission

	type testResult struct {
		response        string
		shouldError     bool
		eventsLen       int
		missionExists   bool
		missionType     gamecomm.MissionType
		scheduleCalls   int
		eventsCancelled bool
	}

	tests := []struct {
		name                      string
		eventSchedulerError       bool
		eventScheduleCallsToError int
		mission                   mission.Mission
		wants                     testResult
		corporationErrors         corporationErrors
		worldErrors               worldErrors
	}{
		{
			name:                "Valid Mission",
			eventSchedulerError: false,
			mission: mission.Mission{
				Id:            "Mission-1",
				CorporationId: 0,
				Squads:        []int{0},
				PlanetId:      "Planet1",
				Status:        "In Progress",
				Type:          gamecomm.SquadMission,
				Resources:     []string{"iron"},
			},
			worldErrors: worldErrors{
				getPlanetError: testError{
					shouldError: false,
				},
			},
			corporationErrors: corporationErrors{
				getSquadError: testError{
					shouldError: false,
				},
			},
			wants: testResult{
				response:      "",
				shouldError:   false,
				eventsLen:     3,
				missionExists: true,
				missionType:   gamecomm.SquadMission,
				scheduleCalls: 3,
			},
		},
		{
			name:                      "Event Schedule Error - With 1 Call",
			eventSchedulerError:       true,
			eventScheduleCallsToError: 1,
			mission: mission.Mission{
				Id:            "Mission-1",
				CorporationId: 0,
				Squads:        []int{0},
				PlanetId:      "Planet1",
				Status:        "In Progress",
				Type:          gamecomm.SquadMission,
				Resources:     []string{"iron"},
			},
			worldErrors: worldErrors{
				getPlanetError: testError{
					shouldError: false,
				},
			},
			corporationErrors: corporationErrors{
				getSquadError: testError{
					shouldError: false,
				},
			},
			wants: testResult{
				response:      "",
				shouldError:   true,
				eventsLen:     0,
				missionExists: false,
				missionType:   gamecomm.SquadMission,
				scheduleCalls: 0,
			},
		},
		{
			name:                      "Event Schedule Error - With 2 Calls",
			eventSchedulerError:       true,
			eventScheduleCallsToError: 2,
			mission: mission.Mission{
				Id:            "Mission-1",
				CorporationId: 0,
				Squads:        []int{0},
				PlanetId:      "Planet1",
				Status:        "In Progress",
				Type:          gamecomm.SquadMission,
				Resources:     []string{"iron"},
			},
			worldErrors: worldErrors{
				getPlanetError: testError{
					shouldError: false,
				},
			},
			corporationErrors: corporationErrors{
				getSquadError: testError{
					shouldError: false,
				},
			},
			wants: testResult{
				response:        "",
				shouldError:     true,
				eventsLen:       1,
				missionExists:   false,
				missionType:     gamecomm.SquadMission,
				scheduleCalls:   1,
				eventsCancelled: true,
			},
		},
		{
			name:                      "Event Schedule Error - With 3 Call",
			eventSchedulerError:       true,
			eventScheduleCallsToError: 3,
			mission: mission.Mission{
				Id:            "Mission-1",
				CorporationId: 0,
				Squads:        []int{0},
				PlanetId:      "Planet1",
				Status:        "In Progress",
				Type:          gamecomm.SquadMission,
				Resources:     []string{"iron"},
			},
			worldErrors: worldErrors{
				getPlanetError: testError{
					shouldError: false,
				},
			},
			corporationErrors: corporationErrors{
				getSquadError: testError{
					shouldError: false,
				},
			},
			wants: testResult{
				response:        "",
				shouldError:     true,
				eventsLen:       2,
				missionExists:   false,
				missionType:     gamecomm.SquadMission,
				scheduleCalls:   2,
				eventsCancelled: true,
			},
		},
		{
			name:                "Empty Squad List",
			eventSchedulerError: false,
			mission: mission.Mission{
				Id:            "Mission-1",
				CorporationId: 0,
				Squads:        []int{},
				PlanetId:      "Planet1",
				Status:        "In Progress",
				Type:          gamecomm.SquadMission,
				Resources:     []string{"iron"},
			},
			worldErrors: worldErrors{
				getPlanetError: testError{
					shouldError: false,
				},
			},
			corporationErrors: corporationErrors{
				getSquadError: testError{
					shouldError: false,
				},
			},
			wants: testResult{
				response:      "",
				shouldError:   true,
				eventsLen:     0,
				missionExists: false,
				missionType:   gamecomm.SquadMission,
				scheduleCalls: 0,
			},
		},
		{
			name:                "Error on Get Planet",
			eventSchedulerError: false,
			mission: mission.Mission{
				Id:            "Mission-1",
				CorporationId: 0,
				Squads:        []int{0},
				PlanetId:      "InvalidPlanet",
				Status:        "In Progress",
				Type:          gamecomm.SquadMission,
				Resources:     []string{"iron"},
			},
			worldErrors: worldErrors{
				getPlanetError: testError{
					shouldError: true,
					errorStr:    "error: planet not found",
				},
			},
			corporationErrors: corporationErrors{
				getSquadError: testError{
					shouldError: false,
				},
			},
			wants: testResult{
				response:      "",
				shouldError:   true,
				eventsLen:     0,
				missionExists: false,
				missionType:   gamecomm.SquadMission,
				scheduleCalls: 0,
			},
		},
		{
			name:                "Empty Resource List",
			eventSchedulerError: false,
			mission: mission.Mission{
				Id:            "Mission-1",
				CorporationId: 0,
				Squads:        []int{0},
				PlanetId:      "Planet1",
				Status:        "In Progress",
				Type:          gamecomm.SquadMission,
				Resources:     []string{},
			},
			worldErrors: worldErrors{
				getPlanetError: testError{
					shouldError: false,
				},
			},
			corporationErrors: corporationErrors{
				getSquadError: testError{
					shouldError: false,
				},
			},
			wants: testResult{
				response:      "",
				shouldError:   false,
				eventsLen:     3,
				missionExists: true,
				missionType:   gamecomm.SquadMission,
				scheduleCalls: 3,
			},
		},
		{
			name:                "Error on Get Squad",
			eventSchedulerError: false,
			mission: mission.Mission{
				Id:            "Mission-1",
				CorporationId: 0,
				Squads:        []int{0},
				PlanetId:      "Planet1",
				Status:        "In Progress",
				Type:          gamecomm.SquadMission,
				Resources:     []string{"iron"},
			},
			worldErrors: worldErrors{
				getPlanetError: testError{
					shouldError: false,
				},
			},
			corporationErrors: corporationErrors{
				getSquadError: testError{
					shouldError: true,
					errorStr:    "error: squad not found",
				},
			},
			wants: testResult{
				response:      "",
				shouldError:   true,
				eventsLen:     0,
				missionExists: false,
				missionType:   gamecomm.SquadMission,
				scheduleCalls: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var wg sync.WaitGroup

			gameChannels := &gamecomm.GameChannels{
				WorldChannel:   make(chan gamecomm.WorldCommand),
				CorpChannel:    make(chan gamecomm.CorpCommand),
				MissionChannel: make(chan gamecomm.MissionCommand),
			}

			notificationChannel := make(chan string)

			tt.mission.NotificationChan = notificationChannel

			// Listen world channel
			wg.Add(2)
			listenWordlWorker(t, tt.worldErrors, gameChannels, &wg)
			listenCorporationWorker(t, tt.corporationErrors, gameChannels, &wg)

			missions := make(map[string]*mission.Mission, 0)
			gc := gameclock.NewGameClock(0, 1)

			mockEventScheduler := newMockScheduler(gameChannels, missions, gc, tt.eventSchedulerError, tt.eventScheduleCallsToError)

			ms := createTestMissionScheduller(missions, gameChannels, gc, mockEventScheduler)
			go ms.Run()

			uuid, err := uuid.NewUUID()
			if err != nil {
				t.Fatalf("error: %v", err)
			}

			missionId := uuid.String()

			tt.mission.Id = missionId

			wg.Add(1)
			go func(shouldError bool, notificationChan chan string) {
				if shouldError {
					msg := <-notificationChan
					assert.StringContains(t, msg, "error")
				}

				wg.Done()
			}(tt.wants.shouldError, tt.mission.NotificationChan)

			ms.StartMission(tt.mission)

			// Validate Events where created
			eventsLen := len(mockEventScheduler.events)
			assert.Equal(t, eventsLen, tt.wants.eventsLen)

			eventsQueueLen := len(mockEventScheduler.queue)
			assert.Equal(t, eventsQueueLen, tt.wants.eventsLen)

			if tt.wants.eventsCancelled {
				for _, e := range mockEventScheduler.events {
					assert.Equal(t, e.Cancelled, true)
				}
			}

			// Validate Mission is added to missionMap
			mis, ok := ms.Missions[tt.mission.Id]

			// Mission Exists On Map
			assert.Equal(t, ok, tt.wants.missionExists)

			if tt.wants.missionExists {
				// Check Mission Type
				assert.Equal(t, mis.Type, tt.wants.missionType)
			}

			// Ensure that Schedule is invoked
			scheduleCall := mockEventScheduler.calledFunctions["Schedule"]
			assert.Equal(t, scheduleCall, tt.wants.scheduleCalls)

			close(gameChannels.WorldChannel)
			close(gameChannels.CorpChannel)
			close(gameChannels.MissionChannel)

			wg.Wait()

			close(notificationChannel)

		})
	}
}

func TestStartTransferMission(t *testing.T) {
	// Create Mission
	// Pass mock eventSchduler
	// Check that mission is added to the missions map
	// Check that the correct type of mission is created
	// Ensure that Schedule is invoked
	// Check CreateMission error
	// Check that notificationChan receive correct message
	// Test for race conditions
	// Check that missions have the correct parameters
	// Empty or malformed mission

	type testResult struct {
		response        string
		shouldError     bool
		eventsLen       int
		missionExists   bool
		missionType     gamecomm.MissionType
		scheduleCalls   int
		eventsCancelled bool
	}

	tests := []struct {
		name                      string
		eventSchedulerError       bool
		eventScheduleCallsToError int
		mission                   mission.Mission
		wants                     testResult
		corporationErrors         corporationErrors
		worldErrors               worldErrors
	}{
		{
			name:                "Valid Mission",
			eventSchedulerError: false,
			mission: mission.Mission{
				Id:            "Mission-1",
				CorporationId: 0,
				Squads:        []int{0},
				PlanetId:      "Planet1",
				Status:        "In Progress",
				Type:          gamecomm.TransferMission,
				Resources:     []string{"iron"},
			},
			worldErrors: worldErrors{
				getPlanetError: testError{
					shouldError: false,
				},
			},
			corporationErrors: corporationErrors{
				getSquadError: testError{
					shouldError: false,
				},
			},
			wants: testResult{
				response:      "",
				shouldError:   false,
				eventsLen:     3,
				missionExists: true,
				missionType:   gamecomm.TransferMission,
				scheduleCalls: 3,
			},
		},
		{
			name:                      "Event Schedule Error With 1 Call",
			eventSchedulerError:       true,
			eventScheduleCallsToError: 1,
			mission: mission.Mission{
				Id:            "Mission-1",
				CorporationId: 0,
				Squads:        []int{0},
				PlanetId:      "Planet1",
				Status:        "In Progress",
				Type:          gamecomm.TransferMission,
				Resources:     []string{"iron"},
			},
			worldErrors: worldErrors{
				getPlanetError: testError{
					shouldError: false,
				},
			},
			corporationErrors: corporationErrors{
				getSquadError: testError{
					shouldError: false,
				},
			},
			wants: testResult{
				response:      "",
				shouldError:   true,
				eventsLen:     0,
				missionExists: false,
				missionType:   gamecomm.TransferMission,
				scheduleCalls: 0,
			},
		},
		{
			name:                      "Event Schedule Error - With 2 Calls",
			eventSchedulerError:       true,
			eventScheduleCallsToError: 2,
			mission: mission.Mission{
				Id:            "Mission-1",
				CorporationId: 0,
				Squads:        []int{0},
				PlanetId:      "Planet1",
				Status:        "In Progress",
				Type:          gamecomm.TransferMission,
				Resources:     []string{"iron"},
			},
			worldErrors: worldErrors{
				getPlanetError: testError{
					shouldError: false,
				},
			},
			corporationErrors: corporationErrors{
				getSquadError: testError{
					shouldError: false,
				},
			},
			wants: testResult{
				response:        "",
				shouldError:     true,
				eventsLen:       1,
				missionExists:   false,
				missionType:     gamecomm.TransferMission,
				scheduleCalls:   1,
				eventsCancelled: true,
			},
		},
		{
			name:                      "Event Schedule Error - With 3 Calls",
			eventSchedulerError:       true,
			eventScheduleCallsToError: 3,
			mission: mission.Mission{
				Id:            "Mission-1",
				CorporationId: 0,
				Squads:        []int{0},
				PlanetId:      "Planet1",
				Status:        "In Progress",
				Type:          gamecomm.TransferMission,
				Resources:     []string{"iron"},
			},
			worldErrors: worldErrors{
				getPlanetError: testError{
					shouldError: false,
				},
			},
			corporationErrors: corporationErrors{
				getSquadError: testError{
					shouldError: false,
				},
			},
			wants: testResult{
				response:      "",
				shouldError:   true,
				eventsLen:     2,
				missionExists: false,
				missionType:   gamecomm.TransferMission,
				scheduleCalls: 2,
			},
		},

		{
			name:                "Empty Squad List",
			eventSchedulerError: false,
			mission: mission.Mission{
				Id:            "Mission-1",
				CorporationId: 0,
				Squads:        []int{},
				PlanetId:      "Planet1",
				Status:        "In Progress",
				Type:          gamecomm.TransferMission,
				Resources:     []string{"iron"},
			},
			worldErrors: worldErrors{
				getPlanetError: testError{
					shouldError: false,
				},
			},
			corporationErrors: corporationErrors{
				getSquadError: testError{
					shouldError: false,
				},
			},
			wants: testResult{
				response:      "",
				shouldError:   true,
				eventsLen:     0,
				missionExists: false,
				missionType:   gamecomm.TransferMission,
				scheduleCalls: 0,
			},
		},
		{
			name:                "Error on Get Planet",
			eventSchedulerError: false,
			mission: mission.Mission{
				Id:            "Mission-1",
				CorporationId: 0,
				Squads:        []int{0},
				PlanetId:      "InvalidPlanet",
				Status:        "In Progress",
				Type:          gamecomm.TransferMission,
				Resources:     []string{"iron"},
			},
			worldErrors: worldErrors{
				getPlanetError: testError{
					shouldError: true,
					errorStr:    "error: planet not found",
				},
			},
			corporationErrors: corporationErrors{
				getSquadError: testError{
					shouldError: false,
				},
			},
			wants: testResult{
				response:      "",
				shouldError:   true,
				eventsLen:     0,
				missionExists: false,
				missionType:   gamecomm.TransferMission,
				scheduleCalls: 0,
			},
		},
		{
			name:                "Empty Resource List",
			eventSchedulerError: false,
			mission: mission.Mission{
				Id:            "Mission-1",
				CorporationId: 0,
				Squads:        []int{0},
				PlanetId:      "Planet1",
				Status:        "In Progress",
				Type:          gamecomm.TransferMission,
				Resources:     []string{},
			},
			worldErrors: worldErrors{
				getPlanetError: testError{
					shouldError: false,
				},
			},
			corporationErrors: corporationErrors{
				getSquadError: testError{
					shouldError: false,
				},
			},
			wants: testResult{
				response:      "",
				shouldError:   false,
				eventsLen:     3,
				missionExists: true,
				missionType:   gamecomm.TransferMission,
				scheduleCalls: 3,
			},
		},
		{
			name:                "Error on Get Squad",
			eventSchedulerError: false,
			mission: mission.Mission{
				Id:            "Mission-1",
				CorporationId: 0,
				Squads:        []int{0},
				PlanetId:      "Planet1",
				Status:        "In Progress",
				Type:          gamecomm.TransferMission,
				Resources:     []string{"iron"},
			},
			worldErrors: worldErrors{
				getPlanetError: testError{
					shouldError: false,
				},
			},
			corporationErrors: corporationErrors{
				getSquadError: testError{
					shouldError: true,
					errorStr:    "error: squad not found",
				},
			},
			wants: testResult{
				response:      "",
				shouldError:   true,
				eventsLen:     0,
				missionExists: false,
				missionType:   gamecomm.TransferMission,
				scheduleCalls: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var wg sync.WaitGroup

			gameChannels := &gamecomm.GameChannels{
				WorldChannel:   make(chan gamecomm.WorldCommand),
				CorpChannel:    make(chan gamecomm.CorpCommand),
				MissionChannel: make(chan gamecomm.MissionCommand),
			}

			notificationChannel := make(chan string)

			tt.mission.NotificationChan = notificationChannel

			// Listen world channel
			wg.Add(2)
			listenWordlWorker(t, tt.worldErrors, gameChannels, &wg)
			listenCorporationWorker(t, tt.corporationErrors, gameChannels, &wg)

			missions := make(map[string]*mission.Mission, 0)
			gc := gameclock.NewGameClock(0, 1)

			mockEventScheduler := newMockScheduler(gameChannels, missions, gc, tt.eventSchedulerError, tt.eventScheduleCallsToError)

			ms := createTestMissionScheduller(missions, gameChannels, gc, mockEventScheduler)
			go ms.Run()

			uuid, err := uuid.NewUUID()
			if err != nil {
				t.Fatalf("error: %v", err)
			}

			missionId := uuid.String()

			tt.mission.Id = missionId

			wg.Add(1)
			go func(shouldError bool, notificationChan chan string) {
				if shouldError {
					msg := <-notificationChan
					assert.StringContains(t, msg, "error")

				}

				wg.Done()
			}(tt.wants.shouldError, tt.mission.NotificationChan)

			ms.StartMission(tt.mission)

			// Validate Events where created
			eventsLen := len(mockEventScheduler.events)
			assert.Equal(t, eventsLen, tt.wants.eventsLen)

			eventsQueueLen := len(mockEventScheduler.queue)
			assert.Equal(t, eventsQueueLen, tt.wants.eventsLen)

			if tt.wants.eventsCancelled {
				for _, e := range mockEventScheduler.events {
					assert.Equal(t, e.Cancelled, true)
				}
			}

			// Validate Mission is added to missionMap
			mis, ok := ms.Missions[tt.mission.Id]

			// Mission Exists On Map
			assert.Equal(t, ok, tt.wants.missionExists)

			if tt.wants.missionExists {
				// Check Mission Type
				assert.Equal(t, mis.Type, tt.wants.missionType)
			}

			// Ensure that Schedule is invoked
			scheduleCall := mockEventScheduler.calledFunctions["Schedule"]
			assert.Equal(t, scheduleCall, tt.wants.scheduleCalls)

			close(gameChannels.WorldChannel)
			close(gameChannels.CorpChannel)
			close(gameChannels.MissionChannel)

			wg.Wait()

			close(notificationChannel)

		})
	}
}

// uuid, err := uuid.NewUUID()
// if err != nil {
// 	return Mission{}, fmt.Errorf("error: %v", err)
// }
//
// missionId := uuid.String()
//
// mission := Mission{
// 	Id:               missionId,
// 	CorporationId:    mc.CorporationId,
// 	Squads:           mc.Squads,
// 	PlanetId:         mc.PlanetId,
// 	DestinationTime:  mc.DestinationTime,
// 	ReturnalTime:     mc.ReturnalTime,
// 	Status:           "In Progress",
// 	Type:             mc.Type,
// 	Resources:        mc.Resources,
// 	NotificationChan: mc.NotificationChan,
// 	Amount:           mc.Amount,
// }
//
// return mission, nil

// func (ms *mission.MissionScheduler) StartMission(m mission.Mission) {
// 	ms.RW.Lock()
// 	ms.missions[m.Id] = &m
// 	ms.RW.Unlock()
//
// 	switch m.Type {
// 	case gamecomm.SquadMission:
// 		ms.CreateSquadMission(m)
// 	case gamecomm.TransferMission:
// 		err := ms.CreateTransferMission(m)
// 		if err != nil {
// 			m.NotificationChan <- err.Error()
// 		}
// 	default:
// 		ms.RW.Lock()
// 		delete(ms.missions, m.Id)
// 		ms.RW.Unlock()
// 	}
// }
