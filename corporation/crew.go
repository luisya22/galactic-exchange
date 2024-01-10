package corporation

import "github.com/luisya22/galactic-exchange/gamecomm"

type CrewMember struct {
	ID         uint64
	Name       string
	Species    string
	Skills     map[string]int
	AssignedTo uint64
}

func (cm *CrewMember) Copy() gamecomm.CrewMember {
	return gamecomm.CrewMember{
		ID:         cm.ID,
		Name:       cm.Name,
		Species:    cm.Species,
		Skills:     cm.Skills,
		AssignedTo: cm.AssignedTo,
	}
}
