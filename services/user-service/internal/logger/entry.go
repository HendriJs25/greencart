package logger

import (
	"context"
	"user-service/internal/common/helper"

	"github.com/sirupsen/logrus"
)

func WithContext(ctx context.Context) *logrus.Entry {
	fields := logrus.Fields{}

	if requestID := helper.GetRequestID(ctx); requestID != "" {
		fields["request_id"] = requestID
	}
	return logrus.WithFields(fields)
}
