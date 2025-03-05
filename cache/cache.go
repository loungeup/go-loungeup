package cache

import "time"

type Reader interface {
	Read(key string) any
}

type Writer interface {
	// Write the value to the cache with the default duration.
	Write(key string, value any)

	// Write the value to the cache with the given duration.
	WriteWithDuration(key string, value any, duration time.Duration)
}

type ReadWriter interface {
	Reader
	Writer
}
