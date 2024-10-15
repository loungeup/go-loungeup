package resresultsets

import "github.com/jirenius/go-res"

type KeysetPaginationModel struct {
	LastKey   any     `json:"lastKey,omitempty"`
	ResultSet res.Ref `json:"resultSet"`
}
