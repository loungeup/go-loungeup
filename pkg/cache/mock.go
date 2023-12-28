package cache

type Mock struct {
	ReadFunc  func(key string) any
	WriteFunc func(key string, value any)
}

var _ ReadWriter = (*Mock)(nil)

func (m *Mock) Read(key string) any         { return m.ReadFunc(key) }
func (m *Mock) Write(key string, value any) { m.WriteFunc(key, value) }
