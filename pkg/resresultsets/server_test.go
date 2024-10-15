package resresultsets

import (
	"encoding/json"
	"testing"

	"github.com/jirenius/go-res"
	"github.com/jirenius/go-res/restest"
	"github.com/loungeup/go-loungeup/pkg/cache"
	"github.com/stretchr/testify/require"
)

func TestServer(t *testing.T) {
	cache, err := cache.NewRistretto(cache.MediumRistrettoCache)
	require.NoError(t, err)

	server := NewServer(cache, res.NewService("test"))

	session := restest.NewSession(t, server.service)
	defer session.Close()

	setRID := server.CreateResultSet([]res.Ref{res.Ref("foo"), res.Ref("bar")})

	session.Get(setRID).Response().AssertCollection(json.RawMessage(`[
		{"rid": "foo"},
		{"rid": "bar"}
	]`))
}
