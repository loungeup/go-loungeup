package client

import (
	"encoding/json"
	"fmt"

	"github.com/jirenius/go-res"
	"github.com/jirenius/go-res/resprot"
	"github.com/loungeup/go-loungeup/resmodels"
	"github.com/loungeup/go-loungeup/transport"
)

//go:generate mockgen -source entity.go -destination=./mocks/mock_entity.go -package=mocks

type EntitiesManager interface {
	ReadEntity(selector *resmodels.EntitySelector) (*resmodels.Entity, error)
	ReadEntityAccounts(selector *resmodels.EntityAccountsSelector) ([]*resmodels.Entity, error)
	ReadAccountParents(selector *resmodels.EntitySelector) ([]*resmodels.Entity, error)
	ReadEntityCustomFields(selector *resmodels.EntityCustomFieldsSelector) (*resmodels.EntityCustomFields, error)
	PatchEntity(selector *resmodels.EntitySelector, updates *resmodels.EntityUpdates) error
	ReadEntityFeatures(selector *resmodels.EntitySelector) (*resmodels.EntityFeatures, error)
	BuildESQueryEntity(selector *resmodels.EntitySelector, params *resmodels.BuildEntityESQueryParams) (json.RawMessage, error)
}

// EntitiesClient provides methods to interact with entities.
type EntitiesClient struct {
	base *BaseClient
}

func NewEntitiesClient(base *BaseClient) *EntitiesClient {
	return &EntitiesClient{
		base: base,
	}
}

func (c *EntitiesClient) BuildESQueryEntity(
	selector *resmodels.EntitySelector,
	params *resmodels.BuildEntityESQueryParams,
) (json.RawMessage, error) {
	return transport.CallRESResult[json.RawMessage](
		c.base.resClient,
		"guestprofile.entities."+selector.EntityID.String()+".build-elasticsearch-query",
		resprot.Request{
			Params: params,
			Token:  json.RawMessage(`{"agentRoles": ["service"]}`),
		},
	)
}

func (c *EntitiesClient) ReadEntity(selector *resmodels.EntitySelector) (*resmodels.Entity, error) {
	return c.readEntityByRID(selector.RID())
}

func (c *EntitiesClient) ReadEntityAccounts(selector *resmodels.EntityAccountsSelector) ([]*resmodels.Entity, error) {
	references, err := transport.GetRESCollection[res.Ref](
		c.base.resClient,
		selector.RID(),
		resprot.Request{
			Query: selector.EncodedQuery(),
		},
	)
	if err != nil {
		return nil, err
	}

	result := []*resmodels.Entity{}

	for _, reference := range references {
		account, err := c.readEntityByRID(string(reference))
		if err != nil {
			return nil, err
		}

		result = append(result, account)
	}

	return result, nil
}

func (c *EntitiesClient) ReadAccountParents(selector *resmodels.EntitySelector) ([]*resmodels.Entity, error) {
	entity, err := c.readEntityByRID(selector.RID())
	if err != nil {
		return nil, err
	}

	if entity.Type != resmodels.EntityTypeAccount {
		return nil, fmt.Errorf("entity is not an account")
	}

	result := []*resmodels.Entity{}

	if entity.Chain != "" {
		chain, err := c.readEntityByRID(string(entity.Chain))
		if err != nil {
			return nil, err
		}

		result = append(result, chain)
	}

	if entity.Group != "" {
		group, err := c.readEntityByRID(string(entity.Group))
		if err != nil {
			return nil, err
		}

		result = append(result, group)
	}

	return result, nil
}

func (c *EntitiesClient) ReadEntityCustomFields(
	selector *resmodels.EntityCustomFieldsSelector,
) (*resmodels.EntityCustomFields, error) {
	if cachedResult, ok := c.base.ReadCache(selector.RID()).(*resmodels.EntityCustomFields); ok {
		return cachedResult, nil
	}

	result, err := transport.GetRESModel[*resmodels.EntityCustomFields](
		c.base.resClient,
		selector.RID(),
		resprot.Request{},
	)
	if err != nil {
		return nil, err
	}

	c.base.WriteCache(selector.RID(), result)

	return result, nil
}

func (c *EntitiesClient) ReadEntityFeatures(selector *resmodels.EntitySelector) (*resmodels.EntityFeatures, error) {
	cacheKey := selector.RID() + ".features"

	if cachedResult, ok := c.base.cache.Read(cacheKey).(*resmodels.EntityFeatures); ok {
		return cachedResult, nil
	}

	rids, err := transport.GetRESCollection[res.Ref](
		c.base.resClient,
		selector.RID()+".features",
		resprot.Request{},
	)
	if err != nil {
		return nil, err
	}

	rawEntityFeatures := []*resmodels.RawEntityFeature{}

	for _, rid := range rids {
		rawEntityFeature, err := c.readEntityFeatureByRid(string(rid))
		if err != nil {
			return nil, err
		}

		rawEntityFeatures = append(rawEntityFeatures, rawEntityFeature)
	}

	return resmodels.MapRawEntityFeaturesToEntityFeatures(rawEntityFeatures), nil
}

func (c *EntitiesClient) readEntityFeatureByRid(rid string) (
	*resmodels.RawEntityFeature, error,
) {
	result, err := transport.GetRESModel[*resmodels.RawEntityFeature](
		c.base.resClient,
		rid,
		resprot.Request{})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *EntitiesClient) PatchEntity(selector *resmodels.EntitySelector, updates *resmodels.EntityUpdates) error {
	encodedUpdates, err := json.Marshal(updates)
	if err != nil {
		return fmt.Errorf("could not encode updates: %w", err)
	}

	if response := c.base.resClient.Request(
		"call."+selector.RID()+".patch",
		resprot.Request{Params: json.RawMessage(encodedUpdates)},
	); response.HasError() {
		return response.Error
	}

	return nil
}

func (c *EntitiesClient) readEntityByRID(resourceID string) (*resmodels.Entity, error) {
	if cachedResult, ok := c.base.ReadCache(resourceID).(*resmodels.Entity); ok {
		return cachedResult, nil
	}

	result, err := transport.GetRESModel[*resmodels.Entity](c.base.resClient, resourceID, resprot.Request{})
	if err != nil {
		return nil, err
	}

	c.base.WriteCache(resourceID, result)

	return result, nil
}
