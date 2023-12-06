// Package transport acts as a thin layer for the transport layer (HTTP, RES, etc.) It allows us to interact with our
// services. It is used by the client package. The separation between the client and transport packages allows us to
// easily mock the transport layer in our tests.
package transport
