package sentry

import (
	"time"

	gosentry "github.com/getsentry/sentry-go"
)

// Init initializes Sentry package.
func Init(dsn string) error {
	return gosentry.Init(gosentry.ClientOptions{
		Debug:            true,
		Dsn:              dsn,
		AttachStacktrace: true,
		Integrations: func(integrations []gosentry.Integration) []gosentry.Integration {
			return append(integrations, new(ExtractExtra), new(EventFormatter))
		},
	})
}

// Flush flushes sentry logs.
func Flush() {
	gosentry.Flush(2 * time.Second)
}
