package corporation

import "github.com/luisya22/galactic-exchange/ship"

type Squad struct {
	Ships       *ship.Ship
	CrewMembers []*CrewMember
	// Officers []Officers   coming soon...
}
