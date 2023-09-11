package logger

type Logger interface {
	Infof(format string, args ...any)
	Errorf(format string, args ...any)
	Tracef(format string, args ...any)
}
