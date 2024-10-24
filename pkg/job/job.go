package job

import (
	"log/slog"
	"time"

	"github.com/loungeup/go-loungeup/pkg/log"
)

type Controller[JobSlice ~[]Job, Job any] struct {
	logger       *log.Logger
	reader       Reader[JobSlice, Job]
	readInterval time.Duration
	runners      []Runner[Job]
}

func NewController[JobSlice ~[]Job, Job any](
	reader Reader[JobSlice, Job],
	runner Runner[Job],
	options ...controllerOption,
) *Controller[JobSlice, Job] {
	if reader == nil {
		panic("reader is required")
	}

	if runner == nil {
		panic("runner is required")
	}

	const (
		defaultReadInterval = 5 * time.Second
		defaultRunnersCount = 1
	)

	config := &controllerConfig{
		logger:       log.Default(),
		readInterval: defaultReadInterval,
		runnersCount: defaultRunnersCount,
	}
	for _, option := range options {
		option(config)
	}

	return &Controller[JobSlice, Job]{
		logger:       config.logger,
		reader:       reader,
		readInterval: config.readInterval,
		runners: func() []Runner[Job] {
			result := []Runner[Job]{}
			for i := 0; i < config.runnersCount; i++ {
				result = append(result, runner)
			}
			return result
		}(),
	}
}

func WithControllerLogger(logger *log.Logger) controllerOption {
	return func(c *controllerConfig) { c.logger = logger }
}

func WithControllerReadInterval(interval time.Duration) controllerOption {
	return func(c *controllerConfig) { c.readInterval = interval }
}

func WithControllerRunnersCount(count int) controllerOption {
	return func(c *controllerConfig) { c.runnersCount = count }
}

func (m *Controller[JobSlice, Job]) Run() error {
	jobChannel := make(chan Job)

	go func() {
		ticker := time.NewTicker(m.readInterval)
		defer ticker.Stop()

		var lastJob Job

		for {
			select {
			case <-ticker.C:
				m.logger.Debug("Reading jobs")

				jobs, err := m.reader.Read(lastJob)
				if err != nil {
					m.logger.Error("Could not read jobs", slog.Any("error", err))
					continue
				}

				for _, job := range jobs {
					jobChannel <- job
				}

				m.logger.Debug("Jobs read", slog.Int("count", len(jobs)))

				if len(jobs) > 0 {
					lastJob = jobs[len(jobs)-1]
				}
			}
		}
	}()

	for index, runner := range m.runners {
		go func(runnerID int, runner Runner[Job]) {
			l1 := m.logger.With(slog.Int("runnerId", runnerID))

			for job := range jobChannel {
				l1.Debug("Running job")

				if err := runner.Run(job); err != nil {
					l1.Error("Could not run job", slog.Any("error", err))
				}
			}
		}(index+1, runner)
	}

	return nil
}

type Reader[JobSlice ~[]Job, Job any] interface {
	Read(lastJob Job) (JobSlice, error)
}

type Runner[Job any] interface {
	Run(job Job) error
}

type controllerConfig struct {
	logger       *log.Logger
	readInterval time.Duration
	runnersCount int
}

type controllerOption func(*controllerConfig)
