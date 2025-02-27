package resresultsets

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/jirenius/go-res"
	"github.com/jirenius/go-res/restest"
	"github.com/stretchr/testify/require"
)

func TestServer(t *testing.T) {
	server := NewServer(res.NewService("test"), &MockStore{
		ReadByIDFunc: func(id uuid.UUID) (*ResultSet, error) {
			return &ResultSet{
				ID:         id,
				Collection: []res.Ref{res.Ref("foo"), res.Ref("bar")},
			}, nil
		},
		WriteFunc: func(set *ResultSet) error { return nil },
	})

	session := restest.NewSession(t, server.service)
	defer session.Close()

	setRID, err := server.CreateResultSet([]res.Ref{res.Ref("foo"), res.Ref("bar")})
	require.NoError(t, err)

	session.Get(setRID).Response().AssertCollection(json.RawMessage(`[
		{"rid": "foo"},
		{"rid": "bar"}
	]`))
}
