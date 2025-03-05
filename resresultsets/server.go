package resresultsets

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jirenius/go-res"
	"github.com/loungeup/go-loungeup/errors"
)

type Store interface {
	ReadByID(id uuid.UUID) (*ResultSet, error)
	Write(set *ResultSet) error
}

type Server struct {
	service *res.Service
	store   Store
}

// NewServer creates a new server.
func NewServer(service *res.Service, store Store) *Server {
	result := &Server{service, store}
	result.addRESHandlers()

	return result
}

// CreateResultSet with the given collection and returns its RID.
func (server *Server) CreateResultSet(collection any) (string, error) {
	set := &ResultSet{
		ID:         uuid.New(),
		Collection: collection,
	}
	if err := server.store.Write(set); err != nil {
		return "", fmt.Errorf("could not write result set: %w", err)
	}

	return server.makeResultSetRID(set), nil
}

// addRESHandlers to the server.
func (server *Server) addRESHandlers() {
	server.service.Handle("result-sets.$resultSetID", res.GetCollection(func(request res.CollectionRequest) {
		id, err := uuid.Parse(request.PathParam("resultSetID"))
		if err != nil {
			request.Error(&res.Error{
				Code:    res.CodeInvalidParams,
				Message: "Invalid result set ID",
				Data:    err.Error(),
			})

			return
		}

		set, err := server.store.ReadByID(id)
		if err != nil {
			if errors.ErrorCode(err) == errors.CodeNotFound {
				request.NotFound()
			} else {
				request.Error(&res.Error{
					Code:    res.CodeInternalError,
					Message: "Could not read result set",
					Data: map[string]string{
						"errorMessage": err.Error(),
						"id":           id.String(),
					},
				})
			}

			return
		}

		request.Collection(set.Collection)
	}))
}

func (server *Server) makeResultSetRID(set *ResultSet) string {
	return server.service.FullPath() + ".result-sets." + set.ID.String()
}
