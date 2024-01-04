package mission

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/luisya22/galactic-exchange/gameclock"
	"github.com/luisya22/galactic-exchange/gamecomm"
)

type MissionScheduler struct {
	missions       map[string]*Mission
	eventScheduler *EventScheduler
	missionChannel chan gamecomm.MissionCommand
	RW             sync.RWMutex
	gameClock      *gameclock.GameClock
}

type Mission struct {
	Id               string
	CorporationId    uint64
	Squads           []int
	PlanetId         string
	DestinationTime  time.Time
	ReturnalTime     time.Time
	Status           string
	Type             gamecomm.MissionType
	Resources        []string
	NotificationChan chan string
	Amount           int // TODO: You should have and object for transfers {resource, amount}
}

// startMission
// missionStatus -> createEvents
// events -> sendMissionStatus, after each event send notification
// if all events received, end mission and send notification

func NewMissionScheduler(gameChannels *gamecomm.GameChannels, gc *gameclock.GameClock) *MissionScheduler {

	missions := make(map[string]*Mission, 0)
	eventScheduler := NewEventScheduler(gameChannels, missions, gc)

	return &MissionScheduler{
		missions:       missions,
		eventScheduler: eventScheduler,
		missionChannel: gameChannels.MissionChannel,
		gameClock:      gc,
	}
}

func (ms *MissionScheduler) Run() {

	go ms.eventScheduler.Run()

	for m := range ms.missionChannel {

		mission, err := CreateMission(m)
		if err != nil {
			continue
		}

		ms.StartMission(mission)
	}
}

func CreateMission(mc gamecomm.MissionCommand) (Mission, error) {

	uuid, err := uuid.NewUUID()
	if err != nil {
		return Mission{}, fmt.Errorf("Error: %v", err)
	}

	missionId := uuid.String()

	mission := Mission{
		Id:               missionId,
		CorporationId:    mc.CorporationId,
		Squads:           mc.Squads,
		PlanetId:         mc.PlanetId,
		DestinationTime:  mc.DestinationTime,
		ReturnalTime:     mc.ReturnalTime,
		Status:           "In Progress",
		Type:             mc.Type,
		Resources:        mc.Resources,
		NotificationChan: mc.NotificationChan,
		Amount:           mc.Amount,
	}

	return mission, nil
}

func (ms *MissionScheduler) StartMission(m Mission) {
	ms.RW.Lock()
	ms.missions[m.Id] = &m
	ms.RW.Unlock()

	switch m.Type {
	case gamecomm.SquadMission:
		ms.CreateSquadMission(m)
		break
	case gamecomm.TransferMission:
		ms.CreateTransferMission(m)
		break
	default:
		ms.RW.Lock()
		delete(ms.missions, m.Id)
		ms.RW.Unlock()
	}
}
