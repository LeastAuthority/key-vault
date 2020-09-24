package sentry

import (
	gosentry "github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"
)

// Handler is responsible for sending errors to Sentry.
type Handler struct {
	sendSynchronously bool
	logger            *logrus.Logger
}

// NewHandler creates a new handler.
func NewHandler(logger *logrus.Logger, sendSynchronously bool) *Handler {
	return &Handler{
		sendSynchronously: sendSynchronously,
		logger:            logger,
	}
}

// Handle sends the error to Rollbar.
func (h *Handler) Handle(err error) {
	h.logger.WithError(err).Error("got the error in platform errors handler")
	eventID := gosentry.CaptureException(err)

	if h.sendSynchronously && eventID != nil {
		gosentry.Flush(0)
	}
}
