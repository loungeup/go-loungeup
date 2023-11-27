package loungeup

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestEntityRelatedIDs(t *testing.T) {
	var (
		accountID  = uuid.New()
		chainID    = uuid.New()
		groupID    = uuid.New()
		resellerID = uuid.New()
	)

	tests := map[string]struct {
		in   Entity
		want []uuid.UUID
	}{
		"simple": {
			in:   Entity{ID: accountID, ChainID: chainID, GroupID: groupID, ResellerID: resellerID},
			want: []uuid.UUID{accountID, chainID, groupID, resellerID},
		},
		"without chain ID": {
			in:   Entity{ID: accountID, GroupID: groupID, ResellerID: resellerID},
			want: []uuid.UUID{accountID, groupID, resellerID},
		},
		"without group ID": {
			in:   Entity{ID: accountID, ChainID: chainID, ResellerID: resellerID},
			want: []uuid.UUID{accountID, chainID, resellerID},
		},
		"without reseller ID": {
			in:   Entity{ID: accountID, ChainID: chainID, GroupID: groupID},
			want: []uuid.UUID{accountID, chainID, groupID},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.in.RelatedIDs())
		})
	}
}
