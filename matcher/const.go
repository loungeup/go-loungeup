package matcher

type wellKnownMatchers struct {
	Products        string
	RoomTypes       string
	TicketLocations string
	TicketTypes     string
	TicketTargets   string
}

const (
	Products        = "productsMatcher"
	RoomTypes       = "roomTypesMatcher"
	TicketLocations = "ticketLocationsMatcher"
	TicketTypes     = "ticketTypesMatcher"
	TicketTargets   = "ticketTargetsMatcher"
)

var WellKnownMatchers = wellKnownMatchers{
	Products:        Products,
	RoomTypes:       RoomTypes,
	TicketLocations: TicketLocations,
	TicketTypes:     TicketTypes,
	TicketTargets:   TicketTargets,
}

type wellKnownMatcherKeys struct {
	UpgradeProductID Matchable
	Default          Matchable
}

var WellKnownMatcherKeys = wellKnownMatcherKeys{
	UpgradeProductID: "upgradeProductId",
	Default:          "default",
}
