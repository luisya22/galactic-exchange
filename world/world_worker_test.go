package world_test

import (
	"reflect"
	"sync"
	"testing"

	"github.com/luisya22/galactic-exchange/gamecomm"
	assert "github.com/luisya22/galactic-exchange/internal/tests"
	"github.com/luisya22/galactic-exchange/world"
)

func TestGetPlanet(t *testing.T) {
	type testResult struct {
		response    string
		shouldError bool
	}

	// World
	gameChannels := &gamecomm.GameChannels{
		WorldChannel: make(chan gamecomm.WorldCommand, 10),
	}

	w := createTestWorld(t, gameChannels)
	w.Listen()

	// TODO: Nil Planet Map
	// TODO: Boundary Conditions -- Long planet names, special characters
	tests := []struct {
		name        string
		planetName  string
		wants       testResult
		shouldError bool
	}{
		{
			name:       "Valid Planet ID",
			planetName: planet1Name,
			wants: testResult{
				response:    planet1Name,
				shouldError: false,
			},
		},
		{
			name:       "Invalid Planet ID",
			planetName: "Wrong Planet Name",
			wants: testResult{
				response:    "",
				shouldError: true,
			},
		},
		{
			name:       "Empty Planet ID",
			planetName: "",
			wants: testResult{
				response:    "",
				shouldError: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			resChan := make(chan gamecomm.ChanResponse)
			command := gamecomm.WorldCommand{
				PlanetId:        tt.planetName,
				ResponseChannel: resChan,
				Action:          gamecomm.GetPlanet,
			}

			gameChannels.WorldChannel <- command

			res := <-resChan
			if !tt.wants.shouldError {
				assert.NilError(t, res.Err)
			} else {
				assert.Error(t, res.Err)
			}

			resPlanet, ok := res.Val.(world.Planet)
			if !tt.wants.shouldError && !ok {
				t.Fatalf("type conversion failed - got: %v; expected: %v", reflect.TypeOf(res.Val), "world.Planet")
			}

			assert.Equal[string](t, resPlanet.Name, tt.wants.response)
		})
	}

	// Subtest for Concurrent Access
	t.Run("Concurrent Access", func(t *testing.T) {
		var wg sync.WaitGroup
		numGoroutines := 10

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				resChan := make(chan gamecomm.ChanResponse)
				command := gamecomm.WorldCommand{
					PlanetId:        planet1Name,
					ResponseChannel: resChan,
				}

				gameChannels.WorldChannel <- command

				res := <-resChan
				assert.NilError(t, res.Err)

				resPlanet, ok := res.Val.(world.Planet)
				if !ok {
					t.Fatalf("type conversion failed - got: %v; expected: %v", reflect.TypeOf(res.Val), "world.Planet")
				}

				assert.Equal[string](t, resPlanet.Name, planet1Name)
			}()
		}

		wg.Wait()
	})
}

func TestAddResourcesToPlanet(t *testing.T) {

	type testResult struct {
		response    int
		shouldError bool
	}

	// TODO: Nil Planet Map
	// TODO: Boundary Conditions -- Long planet names, special characters
	tests := []struct {
		name        string
		planetName  string
		amount      int
		resource    string
		wants       testResult
		shouldError bool
	}{
		{
			name:       "Valid Amount",
			planetName: planet1Name,
			resource:   string(world.Iron),
			amount:     1,
			wants: testResult{
				response:    resourceQuantity + 1,
				shouldError: false,
			},
		},
		{
			name:       "Negative Amount",
			planetName: planet1Name,
			resource:   string(world.Iron),
			amount:     -1,
			wants: testResult{
				response:    0,
				shouldError: true,
			},
		},
		{
			name:       "Zero Amount",
			planetName: planet1Name,
			resource:   string(world.Iron),
			amount:     0,
			wants: testResult{
				response:    resourceQuantity,
				shouldError: false,
			},
		},
		{
			name:       "Invalid Planet ID",
			planetName: "Wrong Planet Name",
			resource:   string(world.Iron),
			amount:     1,
			wants: testResult{
				response:    0,
				shouldError: true,
			},
		},
		{
			name:       "Empty Planet ID",
			planetName: "",
			resource:   string(world.Iron),
			amount:     1,
			wants: testResult{
				response:    0,
				shouldError: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// World
			gameChannels := &gamecomm.GameChannels{
				WorldChannel: make(chan gamecomm.WorldCommand, 10),
			}

			w := createTestWorld(t, gameChannels)
			w.Listen()

			resChan := make(chan gamecomm.ChanResponse)
			command := gamecomm.WorldCommand{
				PlanetId:        tt.planetName,
				ResponseChannel: resChan,
				Action:          gamecomm.AddResourcesToPlanet,
				Amount:          tt.amount,
				Resource:        tt.resource,
			}

			gameChannels.WorldChannel <- command

			res := <-resChan
			if !tt.wants.shouldError {
				assert.NilError(t, res.Err)
			} else {
				assert.Error(t, res.Err)
			}

			resInt, ok := res.Val.(int)
			if !tt.wants.shouldError && !ok {
				t.Fatalf("type conversion failed - got: %v; expected: %v", reflect.TypeOf(res.Val), "world.Planet")
			}

			var resourceAmount int
			if tt.wants.shouldError {
				resourceAmount = 0
			} else {
				resourceAmount = w.Planets[tt.planetName].Resources[world.Resource(tt.resource)]

			}

			// Validate correct Response
			assert.Equal[int](t, resInt, resourceAmount)

			// Assert it increased correctly
			expectedAmount := resourceQuantity + tt.amount
			if tt.wants.shouldError {
				expectedAmount = 0
			}

			assert.Equal[int](t, resourceAmount, expectedAmount)
		})
	}

	// World
	gameChannels := &gamecomm.GameChannels{
		WorldChannel: make(chan gamecomm.WorldCommand, 10),
	}

	w := createTestWorld(t, gameChannels)
	w.Listen()

	// Subtest for Concurrent Access
	t.Run("Concurrent Access", func(t *testing.T) {
		var wg sync.WaitGroup
		numGoroutines := 10

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				resChan := make(chan gamecomm.ChanResponse)
				command := gamecomm.WorldCommand{
					PlanetId:        planet1Name,
					ResponseChannel: resChan,
					Action:          gamecomm.AddResourcesToPlanet,
					Amount:          1,
					Resource:        string(world.Iron),
				}

				w.Planets[command.PlanetId].RW.RLock()
				resourceAmount := w.Planets[command.PlanetId].Resources[world.Resource(command.Resource)]
				w.Planets[command.PlanetId].RW.RUnlock()

				gameChannels.WorldChannel <- command

				res := <-resChan
				assert.NilError(t, res.Err)

				resInt, ok := res.Val.(int)
				if !ok {
					t.Errorf("type conversion failed - got: %v; expected: %v", reflect.TypeOf(res.Val), "int")
					return
				}

				assert.Greater[int](t, resInt, resourceAmount)
			}()
		}

		wg.Wait()
	})

}

// func (w *World) worker(ch <-chan gamecomm.WorldCommand) {
//
// 	for command := range ch {
// 		switch command.Action {
// 		case gamecomm.GetPlanet:
// 			planet, err := w.GetPlanet(command.PlanetId)
// 			if err != nil {
// 				command.ResponseChannel <- gamecomm.ChanResponse{Err: err}
// 			}
//
// 			// Return chanel
// 			command.ResponseChannel <- gamecomm.ChanResponse{
// 				Val: planet.copy(),
// 				Err: nil,
// 			}
// 		case gamecomm.AddResourcesToPlanet:
// 			amount, err := w.AddResourcesToPlanet(command.PlanetId, Resource(command.Resource), command.Amount)
//
// 			command.ResponseChannel <- gamecomm.ChanResponse{
// 				Val: amount,
// 				Err: err,
// 			}
// 		case gamecomm.RemoveResourcesFromPlanet:
// 			amount, err := w.RemoveResourcesFromPlanet(command.PlanetId, Resource(command.Resource), command.Amount)
//
// 			command.ResponseChannel <- gamecomm.ChanResponse{
// 				Val: amount,
// 				Err: err,
// 			}
//
// 		}
// 	}
// }
