package economy

import "github.com/luisya22/galactic-exchange/internal/gameclock"

type transaction struct {
	planetId      string
	corporationId uint64
	resource      string
	credits       float64
	time          gameclock.GameTime
}

// Save Transactions per zone and globally
func (e *Economy) addTransaction(zoneId string, planetId string, corporationId uint64, resource string, credits float64, t gameclock.GameTime) error {
	e.rw.Lock()
	defer e.rw.Unlock()

	tran := transaction{
		planetId:      planetId,
		corporationId: corporationId,
		resource:      resource,
		credits:       credits,
		time:          t,
	}

	index := len(e.transactions)

	// Transaction
	e.transactions = append(e.transactions, tran)
	if len(e.transactions) > e.limit {
		e.transactions = e.transactions[len(e.transactions)-e.limit:]
	}

	// Zone Transaction
	_, ok := e.zoneTransactions[zoneId]
	if !ok {
		e.zoneTransactions[zoneId] = []int{}
	}

	e.zoneTransactions[zoneId] = append(e.zoneTransactions[zoneId], index)

	// Planet Transaction
	_, ok = e.planetTransactions[planetId]
	if !ok {
		e.planetTransactions[planetId] = []int{}
	}

	e.planetTransactions[planetId] = append(e.planetTransactions[planetId], index)

	// TODO: Save Corporation-Planet Trade Relations level

	return nil
}
