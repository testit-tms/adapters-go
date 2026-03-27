package client_helpers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	tmsclient "github.com/testit-tms/api-client-golang/v3"
	"golang.org/x/exp/slog"
)

func AuthContext(token string) context.Context {
	return context.WithValue(context.Background(), tmsclient.ContextAPIKeys, map[string]tmsclient.APIKey{
		"Bearer or PrivateToken": {
			Key:    token,
			Prefix: "PrivateToken",
		},
	})
}

func LogAndWrapAPIError(l *slog.Logger, op, msg string, err error, resp *http.Response) error {
	if resp != nil && resp.Body != nil {
		if b, readErr := io.ReadAll(resp.Body); readErr == nil && len(b) != 0 {
			l.Error(msg, "error", err, slog.String("response", string(b)), slog.String("op", op))
			return fmt.Errorf("%s: %s: %w", op, msg, err)
		}
	}
	l.Error(msg, "error", err, slog.String("op", op))
	return fmt.Errorf("%s: %s: %w", op, msg, err)
}

func Retry(tries int, delay time.Duration, fn func() error) error {
	var lastErr error
	for i := 0; i < tries; i++ {
		if err := fn(); err == nil {
			return nil
		} else {
			lastErr = err
			time.Sleep(delay)
		}
	}
	return lastErr
}
