package gameclock

import (
	"fmt"
	"sync"
	"time"
)

const (
	hoursPerDay   = 24
	daysPerMonth  = 30
	monthsPerYear = 12
)

const (
	Day   = hoursPerDay
	Month = hoursPerDay * daysPerMonth
	Year  = hoursPerDay * daysPerMonth * monthsPerYear
)

type GameTime uint64
type GameTimeDuration uint64

type GameClock struct {
	currentTime         GameTime
	tickerInterval      time.Duration
	rw                  sync.RWMutex
	newTickerMultiplier chan float64
}

func NewGameClock(initialTime GameTime, gameSpeedMultiplier float64) *GameClock {
	gc := &GameClock{
		currentTime:         initialTime,
		newTickerMultiplier: make(chan float64),
	}

	gc.setTickerInterval(gameSpeedMultiplier)

	return gc
}

func (gc *GameClock) Update() {

	gc.rw.Lock()
	defer gc.rw.Unlock()

	gc.currentTime++
}

func (gc *GameClock) GetCurrentTime() GameTime {
	gc.rw.RLock()
	defer gc.rw.RUnlock()
	return gc.currentTime
}

func (gc *GameClock) GetCurrentDate() string {
	gc.rw.RLock()
	defer gc.rw.RUnlock()

	totalHours := gc.currentTime
	hours := totalHours % hoursPerDay
	days := (totalHours / hoursPerDay) % daysPerMonth
	months := (totalHours / (hoursPerDay * daysPerMonth)) % monthsPerYear
	years := totalHours / (hoursPerDay * daysPerMonth * monthsPerYear)

	return fmt.Sprintf("Year: %d, Month: %d, Day: %d, Hour: %d", years+1, months+1, days+1, hours)
}

func ConvertDateToGameTime(year, month, day, hour uint64) uint64 {
	yearHours := (year - 1) * monthsPerYear * daysPerMonth * hoursPerDay
	monthHours := (month - 1) * daysPerMonth * hoursPerDay
	daysHours := (day - 1) * hoursPerDay

	return yearHours + monthHours + daysHours + hour
}

func (gt GameTime) After(gametime GameTime) bool {
	return gt > gametime
}

func (gt GameTime) Before(gametime GameTime) bool {
	return gt < gametime
}

func (gt GameTime) Add(gtd GameTimeDuration) GameTime {
	return gt + GameTime(gtd)
}

func (gc *GameClock) setTickerInterval(multiplier float64) {
	hoursInYear := 365 * 24

	gc.rw.Lock()
	defer gc.rw.Unlock()

	gc.tickerInterval = time.Duration(float64(time.Hour) / multiplier / float64(hoursInYear))
}

func (gc *GameClock) UpdateTicker(multiplier float64) {
	gc.newTickerMultiplier <- multiplier
}

func (gc *GameClock) StartTime() {

	gc.rw.RLock()
	tickerInterval := gc.tickerInterval
	gc.rw.RUnlock()

	ticker := time.NewTicker(tickerInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			gc.Update()
		case newMultiplier := <-gc.newTickerMultiplier:
			ticker.Stop()
			gc.setTickerInterval(newMultiplier)

			gc.rw.RLock()
			tickerInterval := gc.tickerInterval
			gc.rw.RUnlock()
			ticker = time.NewTicker(tickerInterval)
		}
	}
}

// TODO: Get Date from Current Time
