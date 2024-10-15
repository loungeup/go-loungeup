package resresultsets

import (
	"github.com/google/uuid"
	"github.com/jirenius/go-res"
	"github.com/loungeup/go-loungeup/pkg/cache"
	"github.com/loungeup/go-loungeup/pkg/errors"
	"github.com/loungeup/go-loungeup/pkg/log"
)

type Server struct {
	cache   cache.ReadWriter
	service *res.Service
}

// NewServer creates a new server.
func NewServer(cache cache.ReadWriter, service *res.Service) *Server {
	result := &Server{cache, service}
	result.addRESHandlers()

	return result
}

// CreateResultSet with the given collection and returns its RID.
func (s *Server) CreateResultSet(collection any) string {
	set := &resultSet{
		serviceName: s.service.FullPath(),
		id:          uuid.New(),
		collection:  collection,
	}
	cacheResultSet(s.cache, set)

	return set.rid()
}

// addRESHandlers to the server.
func (s *Server) addRESHandlers() {
	s.service.Handle("result-sets.$resultSetID", res.GetCollection(func(request res.CollectionRequest) {
		set, err := readCachedResultSet(s.cache, request.ResourceName())
		if err != nil {
			errors.LogAndWriteRESError(log.Default(), request, err)
			return
		}

		request.Collection(set.collection)
	}))
}

func cacheResultSet(cache cache.ReadWriter, set *resultSet) {
	cache.Write(set.rid(), set)
}

func readCachedResultSet(cache cache.ReadWriter, rid string) (*resultSet, error) {
	if set, ok := cache.Read(rid).(*resultSet); ok {
		return set, nil
	}

	return nil, &errors.Error{Code: errors.CodeNotFound, Message: "Result set not found"}
}
