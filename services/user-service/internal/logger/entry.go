package logger

import (
	"context"
	"user-service/internal/common/requestcontext"

	"github.com/sirupsen/logrus"
)

func WithContext(ctx context.Context) *logrus.Entry {
	fields := logrus.Fields{}

	if requestID := requestcontext.GetRequestID(ctx); requestID != "" {
		fields["request_id"] = requestID
	}
	return logrus.WithFields(fields)
}
