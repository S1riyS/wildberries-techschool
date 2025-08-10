package slogext

import (
	"log/slog"
)

// Err returns log/slog an Attr that represents an error.
func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
