package resresultsets

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/loungeup/go-loungeup/errors"
	"github.com/nats-io/nats.go/jetstream"
)

type jetStreamKeyValueStore struct{ store jetstream.KeyValue }

func NewJetStreamKeyValueStore(store jetstream.KeyValue) *jetStreamKeyValueStore {
	return &jetStreamKeyValueStore{store}
}

var _ (Store) = (*jetStreamKeyValueStore)(nil)

func (store *jetStreamKeyValueStore) ReadByID(id uuid.UUID) (*ResultSet, error) {
	entry, err := store.store.Get(context.Background(), id.String())
	if err != nil {
		if errors.Is(err, jetstream.ErrKeyNotFound) {
			return nil, &errors.Error{Code: errors.CodeNotFound}
		} else {
			return nil, err
		}
	}

	model := &jetStreamResultSetModel{}
	if err := json.Unmarshal(entry.Value(), model); err != nil {
		return nil, fmt.Errorf("could not decode JetStream result set model: %w", err)
	}

	return mapJetStreamModelToResultSet(model), nil
}

func (store *jetStreamKeyValueStore) Write(set *ResultSet) error {
	encodedModel, err := json.Marshal(mapResultSetToJetStreamModel(set))
	if err != nil {
		return fmt.Errorf("could not encode JetStream result set model: %w", err)
	}

	if _, err := store.store.Put(context.Background(), set.ID.String(), encodedModel); err != nil {
		return fmt.Errorf("could not write result set model to JetStream: %w", err)
	}

	return nil
}

type jetStreamResultSetModel struct {
	ID         uuid.UUID       `json:"id"`
	Collection json.RawMessage `json:"collection"`
}

func mapJetStreamModelToResultSet(model *jetStreamResultSetModel) *ResultSet {
	return &ResultSet{
		ID: model.ID,
		Collection: func() any {
			var result any
			_ = json.Unmarshal(model.Collection, &result)

			return result
		}(),
	}
}

func mapResultSetToJetStreamModel(set *ResultSet) *jetStreamResultSetModel {
	return &jetStreamResultSetModel{
		ID: set.ID,
		Collection: func() json.RawMessage {
			result, _ := json.Marshal(set.Collection)

			return result
		}(),
	}
}
