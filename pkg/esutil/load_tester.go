package esutil

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/loungeup/go-loungeup/pkg/log"
	"github.com/tidwall/gjson"
)

const defaultLoadThreshold = 70

type LoadTester struct {
	client        *elasticsearch.Client
	typedClient   *elasticsearch.TypedClient
	loadThreshold int
	logger        *log.Logger
}

type LoadTesterOption func(*LoadTester)

func NewLoadTesterWithClient(client *elasticsearch.Client, options ...LoadTesterOption) *LoadTester {
	result := &LoadTester{
		client:        client,
		loadThreshold: defaultLoadThreshold,
		logger:        log.Default(),
	}
	for _, option := range options {
		option(result)
	}

	return result
}

func NewLoadTesterWithTypedClient(client *elasticsearch.TypedClient, options ...LoadTesterOption) *LoadTester {
	result := &LoadTester{
		typedClient:   client,
		loadThreshold: defaultLoadThreshold,
		logger:        log.Default(),
	}
	for _, option := range options {
		option(result)
	}

	return result
}

func WithLoadTesterThreshold(threshold int) LoadTesterOption {
	return func(t *LoadTester) { t.loadThreshold = threshold }
}

func WithLoadTesterLogger(logger *log.Logger) LoadTesterOption {
	return func(t *LoadTester) { t.logger = logger }
}

// Test returns true if the load of the server is acceptable, false otherwise. It returns an error if the load can not
// be tested.
func (t *LoadTester) Test() (bool, error) {
	metrics, err := t.fetchClusterMetrics()
	if err != nil {
		return false, fmt.Errorf("could not get server load: %w", err)
	}

	cpuUsage := metrics.cpuUsagePercent()
	memoryUsage := metrics.memoryUsagePercent()

	l1 := t.logger.With(
		slog.Int("cpuUsage", cpuUsage),
		slog.Int("loadThreshold", t.loadThreshold),
		slog.Int("memoryUsage", memoryUsage),
		slog.Int("nodesCount", metrics.nodesCount),
		slog.Int("usedCpuPercent", metrics.usedCPUPercent),
		slog.Int64("availableMemoryBytes", metrics.availableMemoryBytes),
		slog.Int64("usedMemoryBytes", metrics.usedMemoryBytes),
	)

	result := cpuUsage <= t.loadThreshold && memoryUsage <= t.loadThreshold
	if result {
		l1.Debug("Load is acceptable")
	} else {
		l1.Error("Load is too high")
	}

	return result, nil
}

func (t *LoadTester) fetchClusterMetrics() (*clusterMetrics, error) {
	if t.client != nil {
		return t.fetchClusterMetricsWithClient()
	} else if t.typedClient != nil {
		return t.fetchClusterMetricsWithTypedClient()
	}

	return nil, fmt.Errorf("could not fetch server metrics without client")
}

func (t *LoadTester) fetchClusterMetricsWithClient() (*clusterMetrics, error) {
	response, err := t.client.Cluster.Stats(
		t.client.Cluster.Stats.WithFilterPath(strings.Join([]string{
			"nodes.count.total",
			"nodes.jvm.mem.heap_max_in_bytes",
			"nodes.jvm.mem.heap_used_in_bytes",
			"nodes.process.cpu.percent",
		}, ",")),
	)
	if err != nil {
		return nil, fmt.Errorf("could not execute request: %w", err)
	}
	defer response.Body.Close()

	if response.IsError() {
		return nil, fmt.Errorf("invalid status code: %d", response.StatusCode)
	}

	bodyBuilder := &strings.Builder{}
	if _, err := io.Copy(bodyBuilder, response.Body); err != nil {
		return nil, fmt.Errorf("could not read response body: %w", err)
	}

	bodyAsString := bodyBuilder.String()

	return &clusterMetrics{
		availableMemoryBytes: gjson.Get(bodyAsString, "nodes.jvm.mem.heap_max_in_bytes").Int(),
		nodesCount:           int(gjson.Get(bodyAsString, "nodes.count.total").Int()),
		usedCPUPercent:       int(gjson.Get(bodyAsString, "nodes.process.cpu.percent").Int()),
		usedMemoryBytes:      gjson.Get(bodyAsString, "nodes.jvm.mem.heap_used_in_bytes").Int(),
	}, nil
}

func (t *LoadTester) fetchClusterMetricsWithTypedClient() (*clusterMetrics, error) {
	response, err := t.typedClient.Cluster.Stats().FilterPath(
		"nodes.count.total",
		"nodes.jvm.mem.heap_max_in_bytes",
		"nodes.jvm.mem.heap_used_in_bytes",
		"nodes.process.cpu.percent",
	).Do(context.Background())
	if err != nil {
		return nil, fmt.Errorf("could not execute request: %w", err)
	}

	return &clusterMetrics{
		availableMemoryBytes: response.Nodes.Jvm.Mem.HeapMaxInBytes,
		nodesCount:           response.Nodes.Count.Total,
		usedCPUPercent:       response.Nodes.Process.Cpu.Percent,
		usedMemoryBytes:      response.Nodes.Jvm.Mem.HeapUsedInBytes,
	}, nil
}

type clusterMetrics struct {
	availableMemoryBytes int64
	nodesCount           int
	usedCPUPercent       int
	usedMemoryBytes      int64
}

func (m *clusterMetrics) cpuUsagePercent() int {
	return m.usedCPUPercent / m.nodesCount
}

func (m *clusterMetrics) memoryUsagePercent() int {
	return int(m.usedMemoryBytes * 100 / m.availableMemoryBytes)
}
