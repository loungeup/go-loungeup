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
		"status": "started",
		"progress": 0
	}`))

	assert.Equal(t, errors.CodeInvalid, errors.ErrorCode(server.SetTaskProgress(taskRID, -1)))

	assert.NoError(t, server.SetTaskProgress(taskRID, 50))
	session.GetMsg().AssertChangeEvent(taskRID, json.RawMessage(`{
		"status": "started",
		"progress": 50
	}`))

	assert.NoError(t, server.CompleteTask(taskRID, true))
	session.GetMsg().AssertChangeEvent(taskRID, json.RawMessage(`{
		"status": "completed",
		"progress": 100,
		"result": {
			"data": true
		}
	}`))
	session.Get(taskRID).Response().AssertModel(json.RawMessage(`{
		"status": "completed",
		"progress": 100,
		"result": {
			"data": true
		}
	}`))

	assert.NoError(t, server.FailTask(taskRID, &errors.Error{Message: "Unknown error"}))
	session.GetMsg().AssertChangeEvent(taskRID, json.RawMessage(`{
		"status": "failed",
		"progress": 0,
		"error": "Unknown error"
	}`))
}
