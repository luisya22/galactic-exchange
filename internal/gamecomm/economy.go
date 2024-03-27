package gamecomm

// World Channels
type EconomyCommand struct {
	Action          EconomyCommandType
	MarketListingId string
	ZoneId          string
	Amount          int
	Resource        string
	Price           float64
	CorporationId   uint64
	BuyerPlanetId   string
	ResponseChannel chan ChanResponse
}

type EconomyCommandType int

const (
	AddMarketListing EconomyCommandType = iota
	BuyMarketListing
	GetMarketListings
	GetMarketListingsByResource
	EditMarketListingPrice
	GetMarketPrice
)

type MarketListing struct {
	Id            string
	ResourceName  string
	Amount        int
	Price         float64
	CorporationId uint64
}
