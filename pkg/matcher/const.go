package matcher

type wellKnownMatchers struct {
	Products        string
	RoomTypes       string
	TicketLocations string
	TicketTypes     string
	TicketTargets   string
}

var WellKnownMatchers = wellKnownMatchers{
	Products:        "productsMatcher",
	RoomTypes:       "roomTypesMatcher",
	TicketLocations: "ticketLocationsMatcher",
	TicketTypes:     "ticketTypesMatcher",
	TicketTargets:   "ticketTargetsMatcher",
}

type wellKnownMatcherKeys struct {
	UpgradeProductID Matchable
	Default          Matchable
}

var WellKnownMatcherKeys = wellKnownMatcherKeys{
	UpgradeProductID: "upgradeProductId",
	Default:          "default",
}
