package client

// internalClient used to interact with NATS services using the RES protocol. This client is meant to be used by
// internal services for service-to-service communication.
type internalClient struct {
	Bookings      *bookingsClient
	ComputedAttrs *computedAttrsClient
	Currency      *currencyClient
	Entities      *entitiesClient
	Guests        *guestsClient
	Integrations  *integrationsClient
	Products      *productsClient
	ProxyDB       *proxyDBClient
	RoomTypes     *roomTypesClient
	Segments      *segmentsClient
}
