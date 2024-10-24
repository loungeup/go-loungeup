package job_test

import (
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/loungeup/go-loungeup/pkg/job"
	"github.com/stretchr/testify/require"
)

func TestOrchestrator(t *testing.T) {
	currentGuestIDIndex := 1

	wg := &sync.WaitGroup{}
	wg.Add(6)

	manager := &guestIDsMock{
		readFunc: func(lastGuestID uuid.UUID) (uuid.UUIDs, error) {
			if lastGuestID == uuid.Nil {
				return uuid.UUIDs{
					uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					uuid.MustParse("00000000-0000-0000-0000-000000000003"),
				}, nil
			}

			require.Equal(t, "00000000-0000-0000-0000-000000000003", lastGuestID.String())

			return uuid.UUIDs{
				uuid.MustParse("00000000-0000-0000-0000-000000000004"),
				uuid.MustParse("00000000-0000-0000-0000-000000000005"),
				uuid.MustParse("00000000-0000-0000-0000-000000000006"),
			}, nil
		},
		runFunc: func(guestID uuid.UUID) error {
			require.Equal(t, "00000000-0000-0000-0000-00000000000"+strconv.Itoa(currentGuestIDIndex), guestID.String())
			currentGuestIDIndex++
			wg.Done()
			return nil
		},
	}

	require.NoError(t, job.NewController(manager, manager, job.WithControllerReadInterval(time.Second)).Run())
	wg.Wait()
}

type guestIDsMock struct {
	readFunc func(lastGuestID uuid.UUID) (uuid.UUIDs, error)
	runFunc  func(guestID uuid.UUID) error
}

func (m *guestIDsMock) Read(lastGuestID uuid.UUID) (uuid.UUIDs, error) {
	return m.readFunc(lastGuestID)
}

func (m *guestIDsMock) Run(guestID uuid.UUID) error {
	return m.runFunc(guestID)
}
