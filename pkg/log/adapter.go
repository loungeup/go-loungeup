package log

import (
	"fmt"

	"github.com/dgraph-io/badger/v3"
	resLogger "github.com/jirenius/go-res/logger"
)

// Adapter for external libraries.
// It is used to adapt the logger by implementing interfaces required by external libraries.
// It is not intended to be used directly by the application.
type Adapter struct{ underlyingLogger *Logger }

// Interfaces implemented by the adapter.
var (
	_ resLogger.Logger = (*Adapter)(nil)
	_ badger.Logger    = (*Adapter)(nil)
)

func (a *Adapter) Debugf(message string, attributes ...any) {
	a.underlyingLogger.Debug(fmt.Sprintf(message, attributes...))
}

func (a *Adapter) Errorf(message string, attributes ...any) {
	a.underlyingLogger.Error(fmt.Sprintf(message, attributes...))
}

func (a *Adapter) Infof(message string, attributes ...any) {
	a.underlyingLogger.Debug(fmt.Sprintf(message, attributes...))
}

func (a *Adapter) Tracef(message string, attributes ...any) {
	a.underlyingLogger.Debug(fmt.Sprintf(message, attributes...))
}

func (a *Adapter) Warningf(message string, attributes ...any) {
	a.underlyingLogger.Debug(fmt.Sprintf(message, attributes...))
}
