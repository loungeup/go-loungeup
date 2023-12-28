package cache

type Reader interface {
	Read(key string) any
}

type Writer interface {
	Write(key string, value any)
}

type ReadWriter interface {
	Reader
	Writer
}
