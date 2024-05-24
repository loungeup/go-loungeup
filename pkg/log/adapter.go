package log

import (
	"fmt"

	"github.com/dgraph-io/badger/v3"
	resLogger "github.com/jirenius/go-res/logger"
)

// adapter for external libraries.
// It is used to adapt the logger by implementing interfaces required by external libraries.
// It is not intended to be used directly by the application.
type adapter struct{ underlyingLogger *logger }

// Interfaces implemented by the adapter.
var (
	_ resLogger.Logger = (*adapter)(nil)
	_ badger.Logger    = (*adapter)(nil)
)

func (a *adapter) Debugf(message string, attributes ...any) {
	a.underlyingLogger.Debug(fmt.Sprintf(message, attributes...))
}

func (a *adapter) Errorf(message string, attributes ...any) {
	a.underlyingLogger.Error(fmt.Sprintf(message, attributes...))
}

func (a *adapter) Infof(message string, attributes ...any) {
	a.underlyingLogger.Debug(fmt.Sprintf(message, attributes...))
}

func (a *adapter) Tracef(message string, attributes ...any) {
	a.underlyingLogger.Debug(fmt.Sprintf(message, attributes...))
}

func (a *adapter) Warningf(message string, attributes ...any) {
	a.underlyingLogger.Debug(fmt.Sprintf(message, attributes...))
}
