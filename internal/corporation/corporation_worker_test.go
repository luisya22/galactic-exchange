package corporation_test

import (
	"reflect"
	"sync"
	"testing"

	"github.com/luisya22/galactic-exchange/internal/assert"
	"github.com/luisya22/galactic-exchange/internal/gamecomm"
)

func TestGetCorporation(t *testing.T) {
	type testResult struct {
		response    string
		shouldError bool
	}

	// World
	gameChannels := &gamecomm.GameChannels{
		CorpChannel: make(chan gamecomm.CorpCommand, 10),
	}

	cg := createTestCorpGroup(t, gameChannels)
	cg.Listen()

	// TODO: Nil CorpGroup

	tests := []struct {
		name          string
		corporationId uint64
		wants         testResult
		shouldError   bool
	}{
		{
			name:          "Valid Corporation ID",
			corporationId: corporationID,
			wants: testResult{
				response:    corporationName,
				shouldError: false,
			},
		},
		{
			name:          "Invalid Corporation ID",
			corporationId: 999,
			wants: testResult{
				response:    "",
				shouldError: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			resChan := make(chan gamecomm.ChanResponse)
			command := gamecomm.CorpCommand{
				CorporationId:   tt.corporationId,
				ResponseChannel: resChan,
				Action:          gamecomm.GetCorporation,
			}

			gameChannels.CorpChannel <- command

			res := <-resChan
			if !tt.wants.shouldError {
				assert.NilError(t, res.Err)
			} else {
				assert.Error(t, res.Err)
			}

			resCorp, ok := res.Val.(gamecomm.Corporation)
			if !tt.wants.shouldError && !ok {
				t.Fatalf("type conversion failed - got: %v; expected: %v", reflect.TypeOf(res.Val), "gamecomm.Corporation")
			}

			assert.Equal[string](t, resCorp.Name, tt.wants.response)
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
				command := gamecomm.CorpCommand{
					CorporationId:   corporationID,
					ResponseChannel: resChan,
				}

				gameChannels.CorpChannel <- command

				res := <-resChan
				assert.NilError(t, res.Err)

				resCorp, ok := res.Val.(gamecomm.Corporation)
				if !ok {
					t.Errorf("type conversion failed - got: %v; expected: %v", reflect.TypeOf(res.Val), "gamecomm.Planet")
				}

				assert.Equal[string](t, resCorp.Name, corporationName)
			}()
		}

		wg.Wait()
	})

}

func TestGetSquad(t *testing.T) {
	type testResult struct {
		response    uint64
		shouldError bool
	}

	// World
	gameChannels := &gamecomm.GameChannels{
		CorpChannel: make(chan gamecomm.CorpCommand, 10),
	}

	cg := createTestCorpGroup(t, gameChannels)
	cg.Listen()

	// TODO: Nil CorpGroup

	tests := []struct {
		name          string
		corporationId uint64
		squadId       int
		wants         testResult
		shouldError   bool
	}{

		{
			name:          "Valid ID",
			corporationId: corporationID,
			squadId:       0,
			wants: testResult{
				response:    testSquadId,
				shouldError: false,
			},
		},
		{
			name:          "Invalid ID",
			corporationId: corporationID,
			squadId:       999,
			wants: testResult{
				response:    0,
				shouldError: true,
			},
		},
		{
			name:          "Negative ID",
			corporationId: corporationID,
			squadId:       -1,
			wants: testResult{
				response:    0,
				shouldError: true,
			},
		},
		{
			name:          "Invalid Corporation Id",
			corporationId: 999,
			squadId:       0,
			wants: testResult{
				response:    0,
				shouldError: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			resChan := make(chan gamecomm.ChanResponse)
			command := gamecomm.CorpCommand{
				CorporationId:   tt.corporationId,
				SquadIndex:      tt.squadId,
				ResponseChannel: resChan,
				Action:          gamecomm.GetSquad,
			}

			gameChannels.CorpChannel <- command

			res := <-resChan
			if !tt.wants.shouldError {
				assert.NilError(t, res.Err)
			} else {
				assert.Error(t, res.Err)
			}

			resSquad, ok := res.Val.(gamecomm.Squad)
			if !tt.wants.shouldError && !ok {
				t.Fatalf("type conversion failed - got: %v; expected: %v", reflect.TypeOf(res.Val), "gamecomm.Squad")
			}

			assert.Equal[uint64](t, resSquad.Id, tt.wants.response)
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
				command := gamecomm.CorpCommand{
					CorporationId:   corporationID,
					ResponseChannel: resChan,
					Action:          gamecomm.GetSquad,
					SquadIndex:      0,
				}

				gameChannels.CorpChannel <- command

				res := <-resChan
				assert.NilError(t, res.Err)

				resSquad, ok := res.Val.(gamecomm.Squad)
				if !ok {
					t.Errorf("type conversion failed - got: %v; expected: %v", reflect.TypeOf(res.Val), "gamecomm.Squad")
				}

				assert.Equal[uint64](t, resSquad.Id, testSquadId)
			}()
		}

		wg.Wait()
	})

}

func TestAddResourcesToBase(t *testing.T) {
	type testResult struct {
		response    int
		shouldError bool
	}

	// TODO: Nil CorpGroup

	tests := []struct {
		name          string
		corporationId uint64
		resource      string
		amount        int
		wants         testResult
		shouldError   bool
	}{

		{
			name:          "Valid Resource & Valid Amount",
			corporationId: corporationID,
			resource:      "iron",
			amount:        50,
			wants: testResult{
				response:    initialIronQuantity + 50,
				shouldError: false,
			},
		},
		{
			name:          "Valid Resource and Negative Amount",
			corporationId: corporationID,
			resource:      "iron",
			amount:        -50,
			wants: testResult{
				response:    0,
				shouldError: true,
			},
		},
		{
			name:          "Invalid Resource & Valid Amount",
			corporationId: corporationID,
			resource:      "sand",
			amount:        50,
			wants: testResult{
				response:    50,
				shouldError: false,
			},
		},
		{
			name:          "Valid Resource but not on map",
			corporationId: corporationID,
			resource:      "water",
			amount:        50,
			wants: testResult{
				response:    50,
				shouldError: false,
			},
		},
		{
			name:          "Invalid Corporation Id",
			corporationId: 999,
			resource:      "iron",
			wants: testResult{
				response:    0,
				shouldError: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			gameChannels := &gamecomm.GameChannels{
				CorpChannel: make(chan gamecomm.CorpCommand, 10),
			}

			cg := createTestCorpGroup(t, gameChannels)
			cg.Listen()

			resChan := make(chan gamecomm.ChanResponse)
			command := gamecomm.CorpCommand{
				CorporationId:   tt.corporationId,
				ResponseChannel: resChan,
				Action:          gamecomm.AddResourcesToBase,
				Amount:          tt.amount,
				Resource:        tt.resource,
			}

			gameChannels.CorpChannel <- command

			res := <-resChan
			if !tt.wants.shouldError {
				assert.NilError(t, res.Err)
			} else {
				assert.Error(t, res.Err)
			}

			resAmount, ok := res.Val.(int)
			if !tt.wants.shouldError && !ok {
				t.Fatalf("type conversion failed - got: %v; expected: %v", reflect.TypeOf(res.Val), "int")
			}

			assert.Equal[int](t, resAmount, tt.wants.response)
		})
	}

	// Subtest for Concurrent Access
	gameChannels := &gamecomm.GameChannels{
		CorpChannel: make(chan gamecomm.CorpCommand, 10),
	}

	cg := createTestCorpGroup(t, gameChannels)
	cg.Listen()

	t.Run("Concurrent Access", func(t *testing.T) {
		var wg sync.WaitGroup
		numGoroutines := 10
		amountPerGoroutine := 50
		expectedTotalIncrease := numGoroutines * amountPerGoroutine

		cg.RW.RLock()
		initialResourceAmount := cg.Corporations[corporationID].Bases[0].StoredResources["iron"]
		cg.RW.RUnlock()

		expectedFinalAmount := initialResourceAmount + expectedTotalIncrease

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				resChan := make(chan gamecomm.ChanResponse)
				command := gamecomm.CorpCommand{
					CorporationId:   corporationID,
					ResponseChannel: resChan,
					Action:          gamecomm.AddResourcesToBase,
					Amount:          amountPerGoroutine,
					Resource:        "iron",
				}

				gameChannels.CorpChannel <- command
				res := <-resChan
				assert.NilError(t, res.Err)
			}()
		}

		wg.Wait()

		cg.RW.RLock()
		finalResourceAmount := cg.Corporations[corporationID].Bases[0].StoredResources["iron"]
		cg.RW.RUnlock()

		assert.Equal(t, finalResourceAmount, expectedFinalAmount)
	})
}

func TestRemoveResourcesFromBase(t *testing.T) {
	type testResult struct {
		response    int
		shouldError bool
	}

	// TODO: Nil CorpGroup

	tests := []struct {
		name          string
		corporationId uint64
		resource      string
		amount        int
		wants         testResult
		shouldError   bool
	}{

		{
			name:          "Valid Resource & Valid Amount",
			corporationId: corporationID,
			resource:      "iron",
			amount:        50,
			wants: testResult{
				response:    initialIronQuantity - 50,
				shouldError: false,
			},
		},
		{
			name:          "Valid Resource and Negative Amount",
			corporationId: corporationID,
			resource:      "iron",
			amount:        -50,
			wants: testResult{
				response:    0,
				shouldError: true,
			},
		},
		{
			name:          "Invalid Resource & Valid Amount",
			corporationId: corporationID,
			resource:      "sand",
			amount:        50,
			wants: testResult{
				response:    0,
				shouldError: true,
			},
		},
		{
			name:          "Valid Resource but not on map",
			corporationId: corporationID,
			resource:      "water",
			amount:        50,
			wants: testResult{
				response:    0,
				shouldError: true,
			},
		},
		{
			name:          "Invalid Corporation Id",
			corporationId: 999,
			resource:      "iron",
			wants: testResult{
				response:    0,
				shouldError: true,
			},
		},
		{
			name:          "Remove More Than Available",
			corporationId: corporationID,
			resource:      "iron",
			amount:        initialIronQuantity + 1,
			wants: testResult{
				response:    0,
				shouldError: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			gameChannels := &gamecomm.GameChannels{
				CorpChannel: make(chan gamecomm.CorpCommand, 10),
			}

			cg := createTestCorpGroup(t, gameChannels)
			cg.Listen()

			resChan := make(chan gamecomm.ChanResponse)
			command := gamecomm.CorpCommand{
				CorporationId:   tt.corporationId,
				ResponseChannel: resChan,
				Action:          gamecomm.RemoveResourcesFromBase,
				Amount:          tt.amount,
				Resource:        tt.resource,
			}

			gameChannels.CorpChannel <- command

			res := <-resChan
			if !tt.wants.shouldError {
				assert.NilError(t, res.Err)
			} else {
				assert.Error(t, res.Err)
			}

			resAmount, ok := res.Val.(int)
			if !tt.wants.shouldError && !ok {
				t.Fatalf("type conversion failed - got: %v; expected: %v", reflect.TypeOf(res.Val), "int")
			}

			assert.Equal[int](t, resAmount, tt.wants.response)
		})
	}

	// Subtest for Concurrent Access
	gameChannels := &gamecomm.GameChannels{
		CorpChannel: make(chan gamecomm.CorpCommand, 10),
	}

	cg := createTestCorpGroup(t, gameChannels)
	cg.Listen()

	t.Run("Concurrent Access", func(t *testing.T) {
		var wg sync.WaitGroup
		numGoroutines := 10
		amountPerGoroutine := 50
		expectedTotalDecrease := numGoroutines * amountPerGoroutine

		cg.RW.RLock()
		initialResourceAmount := cg.Corporations[corporationID].Bases[0].StoredResources["iron"]
		cg.RW.RUnlock()

		expectedFinalAmount := initialResourceAmount - expectedTotalDecrease

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				resChan := make(chan gamecomm.ChanResponse)
				command := gamecomm.CorpCommand{
					CorporationId:   corporationID,
					ResponseChannel: resChan,
					Action:          gamecomm.RemoveResourcesFromBase,
					Amount:          amountPerGoroutine,
					Resource:        "iron",
				}

				gameChannels.CorpChannel <- command
				res := <-resChan
				assert.NilError(t, res.Err)
			}()
		}

		wg.Wait()

		cg.RW.RLock()
		finalResourceAmount := cg.Corporations[corporationID].Bases[0].StoredResources["iron"]
		cg.RW.RUnlock()

		assert.Equal(t, finalResourceAmount, expectedFinalAmount)
	})
}

func TestAddResourcesToSquad(t *testing.T) {
	type testResult struct {
		response    int
		shouldError bool
	}

	// TODO: Nil CorpGroup

	tests := []struct {
		name          string
		corporationId uint64
		squadId       int
		resource      string
		amount        int
		wants         testResult
		shouldError   bool
	}{

		{
			name:          "Valid Resource & Valid Amount",
			corporationId: corporationID,
			squadId:       0,
			resource:      "iron",
			amount:        50,
			wants: testResult{
				response:    initialIronQuantity + 50,
				shouldError: false,
			},
		},
		{
			name:          "Valid Resource and Negative Amount",
			corporationId: corporationID,
			squadId:       0,
			resource:      "iron",
			amount:        -50,
			wants: testResult{
				response:    0,
				shouldError: true,
			},
		},
		{
			name:          "Invalid Resource & Valid Amount",
			corporationId: corporationID,
			squadId:       0,
			resource:      "sand",
			amount:        50,
			wants: testResult{
				response:    50,
				shouldError: false,
			},
		},
		{
			name:          "Valid Resource but not on map",
			corporationId: corporationID,
			squadId:       0,
			resource:      "water",
			amount:        50,
			wants: testResult{
				response:    50,
				shouldError: false,
			},
		},
		{
			name:          "Invalid Squad Id",
			corporationId: corporationID,
			squadId:       999,
			resource:      "iron",
			wants: testResult{
				response:    0,
				shouldError: true,
			},
		},
		{
			name:          "Invalid Corporation Id",
			corporationId: 999,
			squadId:       0,
			resource:      "iron",
			wants: testResult{
				response:    0,
				shouldError: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			gameChannels := &gamecomm.GameChannels{
				CorpChannel: make(chan gamecomm.CorpCommand, 10),
			}

			cg := createTestCorpGroup(t, gameChannels)
			cg.Listen()

			resChan := make(chan gamecomm.ChanResponse)
			command := gamecomm.CorpCommand{
				CorporationId:   tt.corporationId,
				SquadIndex:      tt.squadId,
				ResponseChannel: resChan,
				Action:          gamecomm.AddResourcesToSquad,
				Amount:          tt.amount,
				Resource:        tt.resource,
			}

			gameChannels.CorpChannel <- command

			res := <-resChan
			if !tt.wants.shouldError {
				assert.NilError(t, res.Err)
			} else {
				assert.Error(t, res.Err)
			}

			resAmount, ok := res.Val.(int)
			if !tt.wants.shouldError && !ok {
				t.Fatalf("type conversion failed - got: %v; expected: %v", reflect.TypeOf(res.Val), "int")
			}

			assert.Equal[int](t, resAmount, tt.wants.response)
		})
	}

	// Subtest for Concurrent Access
	gameChannels := &gamecomm.GameChannels{
		CorpChannel: make(chan gamecomm.CorpCommand, 10),
	}

	cg := createTestCorpGroup(t, gameChannels)
	cg.Listen()

	t.Run("Concurrent Access", func(t *testing.T) {
		var wg sync.WaitGroup
		numGoroutines := 10
		amountPerGoroutine := 50
		expectedTotalDecrease := numGoroutines * amountPerGoroutine

		cg.RW.RLock()
		initialResourceAmount := cg.Corporations[corporationID].Squads[0].Cargo["iron"]
		cg.RW.RUnlock()

		expectedFinalAmount := initialResourceAmount - expectedTotalDecrease

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				resChan := make(chan gamecomm.ChanResponse)
				command := gamecomm.CorpCommand{
					CorporationId:   corporationID,
					SquadIndex:      0,
					ResponseChannel: resChan,
					Action:          gamecomm.RemoveResourcesFromSquad,
					Amount:          amountPerGoroutine,
					Resource:        "iron",
				}

				gameChannels.CorpChannel <- command
				res := <-resChan
				assert.NilError(t, res.Err)
			}()
		}

		wg.Wait()

		cg.RW.RLock()
		finalResourceAmount := cg.Corporations[corporationID].Squads[0].Cargo["iron"]
		cg.RW.RUnlock()

		assert.Equal(t, finalResourceAmount, expectedFinalAmount)
	})
}

func TestRemoveResourcesFromSquad(t *testing.T) {
	type testResult struct {
		response    int
		shouldError bool
	}

	// TODO: Nil CorpGroup

	tests := []struct {
		name          string
		corporationId uint64
		squadId       int
		resource      string
		amount        int
		wants         testResult
		shouldError   bool
	}{

		{
			name:          "Valid Resource & Valid Amount",
			corporationId: corporationID,
			squadId:       0,
			resource:      "iron",
			amount:        50,
			wants: testResult{
				response:    initialIronQuantity - 50,
				shouldError: false,
			},
		},
		{
			name:          "Valid Resource and Negative Amount",
			corporationId: corporationID,
			squadId:       0,
			resource:      "iron",
			amount:        -50,
			wants: testResult{
				response:    0,
				shouldError: true,
			},
		},
		{
			name:          "Invalid Resource & Valid Amount",
			corporationId: corporationID,
			squadId:       0,
			resource:      "sand",
			amount:        50,
			wants: testResult{
				response:    0,
				shouldError: true,
			},
		},
		{
			name:          "Valid Resource but not on map",
			corporationId: corporationID,
			squadId:       0,
			resource:      "water",
			amount:        50,
			wants: testResult{
				response:    0,
				shouldError: true,
			},
		},
		{
			name:          "Invalid Squad Id",
			corporationId: corporationID,
			squadId:       999,
			resource:      "iron",
			wants: testResult{
				response:    0,
				shouldError: true,
			},
		},
		{
			name:          "Invalid Corporation Id",
			corporationId: 999,
			squadId:       0,
			resource:      "iron",
			wants: testResult{
				response:    0,
				shouldError: true,
			},
		},
		{
			name:          "Remove More Than Available",
			corporationId: corporationID,
			squadId:       0,
			resource:      "iron",
			amount:        initialIronQuantity + 1,
			wants: testResult{
				response:    0,
				shouldError: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			gameChannels := &gamecomm.GameChannels{
				CorpChannel: make(chan gamecomm.CorpCommand, 10),
			}

			cg := createTestCorpGroup(t, gameChannels)
			cg.Listen()

			resChan := make(chan gamecomm.ChanResponse)
			command := gamecomm.CorpCommand{
				CorporationId:   tt.corporationId,
				SquadIndex:      tt.squadId,
				ResponseChannel: resChan,
				Action:          gamecomm.RemoveResourcesFromSquad,
				Amount:          tt.amount,
				Resource:        tt.resource,
			}

			gameChannels.CorpChannel <- command

			res := <-resChan
			if !tt.wants.shouldError {
				assert.NilError(t, res.Err)
			} else {
				assert.Error(t, res.Err)
			}

			resAmount, ok := res.Val.(int)
			if !tt.wants.shouldError && !ok {
				t.Fatalf("type conversion failed - got: %v; expected: %v", reflect.TypeOf(res.Val), "int")
			}

			assert.Equal[int](t, resAmount, tt.wants.response)
		})
	}

	// Subtest for Concurrent Access
	gameChannels := &gamecomm.GameChannels{
		CorpChannel: make(chan gamecomm.CorpCommand, 10),
	}

	cg := createTestCorpGroup(t, gameChannels)
	cg.Listen()

	t.Run("Concurrent Access", func(t *testing.T) {
		var wg sync.WaitGroup
		numGoroutines := 10
		amountPerGoroutine := 50
		expectedTotalDecrease := numGoroutines * amountPerGoroutine

		cg.RW.RLock()
		initialResourceAmount := cg.Corporations[corporationID].Squads[0].Cargo["iron"]
		cg.RW.RUnlock()

		expectedFinalAmount := initialResourceAmount - expectedTotalDecrease

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				resChan := make(chan gamecomm.ChanResponse)
				command := gamecomm.CorpCommand{
					CorporationId:   corporationID,
					SquadIndex:      0,
					ResponseChannel: resChan,
					Action:          gamecomm.RemoveResourcesFromSquad,
					Amount:          amountPerGoroutine,
					Resource:        "iron",
				}

				gameChannels.CorpChannel <- command
				res := <-resChan
				assert.NilError(t, res.Err)
			}()
		}

		wg.Wait()

		cg.RW.RLock()
		finalResourceAmount := cg.Corporations[corporationID].Squads[0].Cargo["iron"]
		cg.RW.RUnlock()

		assert.Equal(t, finalResourceAmount, expectedFinalAmount)
	})
}

func TestAddCredits(t *testing.T) {
	type testResult struct {
		response    float64
		shouldError bool
	}

	// TODO: Nil CorpGroup

	tests := []struct {
		name          string
		corporationId uint64
		amount        float64
		wants         testResult
		shouldError   bool
	}{

		{
			name:          "Valid Amount",
			corporationId: corporationID,
			amount:        50,
			wants: testResult{
				response:    initialCorporationCredits + 50,
				shouldError: false,
			},
		},
		{
			name:          "Valid Resource and Negative Amount",
			corporationId: corporationID,
			amount:        -50,
			wants: testResult{
				response:    0,
				shouldError: true,
			},
		},
		{
			name:          "Invalid Corporation Id",
			corporationId: 999,
			wants: testResult{
				response:    0,
				shouldError: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			gameChannels := &gamecomm.GameChannels{
				CorpChannel: make(chan gamecomm.CorpCommand, 10),
			}

			cg := createTestCorpGroup(t, gameChannels)
			cg.Listen()

			resChan := make(chan gamecomm.ChanResponse)
			command := gamecomm.CorpCommand{
				CorporationId:   tt.corporationId,
				ResponseChannel: resChan,
				Action:          gamecomm.AddCredits,
				AmountDecimal:   tt.amount,
			}

			gameChannels.CorpChannel <- command

			res := <-resChan
			if !tt.wants.shouldError {
				assert.NilError(t, res.Err)
			} else {
				assert.Error(t, res.Err)
			}

			resAmount, ok := res.Val.(float64)
			if !tt.wants.shouldError && !ok {
				t.Fatalf("type conversion failed - got: %v; expected: %v", reflect.TypeOf(res.Val), "int")
			}

			assert.Equal(t, resAmount, tt.wants.response)
		})
	}

	// Subtest for Concurrent Access
	gameChannels := &gamecomm.GameChannels{
		CorpChannel: make(chan gamecomm.CorpCommand, 10),
	}

	cg := createTestCorpGroup(t, gameChannels)
	cg.Listen()

	t.Run("Concurrent Access", func(t *testing.T) {
		var wg sync.WaitGroup
		numGoroutines := 10
		amountPerGoroutine := 50.0
		expectedTotalIncrease := float64(numGoroutines) * amountPerGoroutine

		cg.RW.RLock()
		initialCredits := cg.Corporations[corporationID].Credits
		cg.RW.RUnlock()

		expectedFinalAmount := initialCredits + expectedTotalIncrease

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				resChan := make(chan gamecomm.ChanResponse)
				command := gamecomm.CorpCommand{
					CorporationId:   corporationID,
					ResponseChannel: resChan,
					Action:          gamecomm.AddCredits,
					AmountDecimal:   amountPerGoroutine,
				}

				gameChannels.CorpChannel <- command
				res := <-resChan
				assert.NilError(t, res.Err)
			}()
		}

		wg.Wait()

		cg.RW.RLock()
		finalResourceAmount := cg.Corporations[corporationID].Credits
		cg.RW.RUnlock()

		assert.Equal(t, finalResourceAmount, expectedFinalAmount)
	})
}

func TestRemoveCredits(t *testing.T) {
	type testResult struct {
		response    float64
		shouldError bool
	}

	// TODO: Nil CorpGroup

	tests := []struct {
		name          string
		corporationId uint64
		amount        float64
		wants         testResult
		shouldError   bool
	}{

		{
			name:          "Valid Amount",
			corporationId: corporationID,
			amount:        50,
			wants: testResult{
				response:    initialCorporationCredits - 50,
				shouldError: false,
			},
		},
		{
			name:          "Valid Resource and Negative Amount",
			corporationId: corporationID,
			amount:        -50,
			wants: testResult{
				response:    0,
				shouldError: true,
			},
		},
		{
			name:          "Invalid Corporation Id",
			corporationId: 999,
			wants: testResult{
				response:    0,
				shouldError: true,
			},
		},
		{
			name:          "Remove More Credits Than Available",
			corporationId: corporationID,
			amount:        initialCorporationCredits + 1,
			wants: testResult{
				response:    0,
				shouldError: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			gameChannels := &gamecomm.GameChannels{
				CorpChannel: make(chan gamecomm.CorpCommand, 10),
			}

			cg := createTestCorpGroup(t, gameChannels)
			cg.Listen()

			resChan := make(chan gamecomm.ChanResponse)
			command := gamecomm.CorpCommand{
				CorporationId:   tt.corporationId,
				ResponseChannel: resChan,
				Action:          gamecomm.RemoveCredits,
				AmountDecimal:   tt.amount,
			}

			gameChannels.CorpChannel <- command

			res := <-resChan
			if !tt.wants.shouldError {
				assert.NilError(t, res.Err)
			} else {
				assert.Error(t, res.Err)
			}

			resAmount, ok := res.Val.(float64)
			if !tt.wants.shouldError && !ok {
				t.Fatalf("type conversion failed - got: %v; expected: %v", reflect.TypeOf(res.Val), "int")
			}

			assert.Equal(t, resAmount, tt.wants.response)
		})
	}

	// Subtest for Concurrent Access
	gameChannels := &gamecomm.GameChannels{
		CorpChannel: make(chan gamecomm.CorpCommand, 10),
	}

	cg := createTestCorpGroup(t, gameChannels)
	cg.Listen()

	t.Run("Concurrent Access", func(t *testing.T) {
		var wg sync.WaitGroup
		numGoroutines := 10
		amountPerGoroutine := 50.0
		expectedTotalDecrease := float64(numGoroutines) * amountPerGoroutine

		cg.RW.RLock()
		initialCredits := cg.Corporations[corporationID].Credits
		cg.RW.RUnlock()

		expectedFinalAmount := initialCredits - expectedTotalDecrease

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				resChan := make(chan gamecomm.ChanResponse)
				command := gamecomm.CorpCommand{
					CorporationId:   corporationID,
					ResponseChannel: resChan,
					Action:          gamecomm.RemoveCredits,
					AmountDecimal:   amountPerGoroutine,
				}

				gameChannels.CorpChannel <- command
				res := <-resChan
				assert.NilError(t, res.Err)
			}()
		}

		wg.Wait()

		cg.RW.RLock()
		finalResourceAmount := cg.Corporations[corporationID].Credits
		cg.RW.RUnlock()

		assert.Equal(t, finalResourceAmount, expectedFinalAmount)
	})
}
