package slogutil

import (
	"log/slog"

	"github.com/google/uuid"
)

// ErrorAttrKey is the key used for error attributes.
const ErrorAttrKey = "error"

// NewErrorAttr returns an Attr for an error value.
func NewErrorAttr(err error) slog.Attr { return slog.Any(ErrorAttrKey, err) }

// NewTraceIDAttr returns an Attr with a unique trace ID.
func NewTraceIDAttr() slog.Attr { return slog.String("traceId", uuid.NewString()) }
