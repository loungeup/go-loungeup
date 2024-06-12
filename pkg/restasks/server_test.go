package restasks

import (
	"encoding/json"
	"testing"

	"github.com/jirenius/go-res"
	"github.com/jirenius/go-res/restest"
	"github.com/loungeup/go-loungeup/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryServer(t *testing.T) {
	server := NewServer(res.NewService("test"), NewInMemoryStore(10))

	session := restest.NewSession(t, server.service)
	defer session.Close()

	taskRID, err := server.CreateTask()
	assert.NoError(t, err)
	session.Get(taskRID).Response().AssertModel(json.RawMessage(`{
		"status": "started"
	}`))

	assert.NoError(t, server.CompleteTask(taskRID, true))
	session.GetMsg().AssertChangeEvent(taskRID, json.RawMessage(`{
		"status": "completed",
		"result": {
			"data": true
		}
	}`))
	session.Get(taskRID).Response().AssertModel(json.RawMessage(`{
		"status": "completed",
		"result": {
			"data": true
		}
	}`))

	assert.NoError(t, server.FailTask(taskRID, &errors.Error{Message: "Unknown error"}))
	session.GetMsg().AssertChangeEvent(taskRID, json.RawMessage(`{
		"status": "failed",
		"error": "Unknown error"
	}`))
}
