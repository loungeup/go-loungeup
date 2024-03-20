package matcher

type wellKnownMatchers struct {
	Products  string
	RoomTypes string
}

var WellKnownMatchers = wellKnownMatchers{
	Products:  "productsMatcher",
	RoomTypes: "roomTypesMatcher",
}

type wellKnownMatcherKeys struct {
	UpgradeProductID Matchable
}

var WellKnownMatcherKeys = wellKnownMatcherKeys{
	UpgradeProductID: "upgradeProductId",
}
