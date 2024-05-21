package matcher

type wellKnownMatchers struct {
	Products        string
	RoomTypes       string
	TicketLocations string
	TicketTypes     string
	TicketTargets   string
}

const (
	ProductsMatcher        = "productsMatcher"
	RoomTypesMatcher       = "roomTypesMatcher"
	ticketLocationsMatcher = "ticketLocationsMatcher"
	ticketTypesMatcher     = "ticketTypesMatcher"
	ticketTargetsMatcher   = "ticketTargetsMatcher"
)

var WellKnownMatchers = wellKnownMatchers{
	Products:        ProductsMatcher,
	RoomTypes:       RoomTypesMatcher,
	TicketLocations: ticketLocationsMatcher,
	TicketTypes:     ticketTypesMatcher,
	TicketTargets:   ticketTargetsMatcher,
}

type wellKnownMatcherKeys struct {
	UpgradeProductID Matchable
	Default          Matchable
}

var WellKnownMatcherKeys = wellKnownMatcherKeys{
	UpgradeProductID: "upgradeProductId",
	Default:          "default",
}
