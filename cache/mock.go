package cache

import "time"

type Mock struct {
	ReadFunc              func(key string) any
	WriteFunc             func(key string, value any)
	WriteWithDurationFunc func(key string, value any, duration time.Duration)
}

var _ ReadWriter = (*Mock)(nil)

func (m *Mock) Read(key string) any {
	return m.ReadFunc(key)
}

func (m *Mock) Write(key string, value any) {
	m.WriteFunc(key, value)
}

func (m *Mock) WriteWithDuration(key string, value any, duration time.Duration) {
	m.WriteWithDurationFunc(key, value, duration)
}
