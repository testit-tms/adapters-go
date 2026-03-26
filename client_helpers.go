package tms

import (
	"context"
	"fmt"
	"io"
	"net/http"

	tmsclient "github.com/testit-tms/api-client-golang/v3"
	"golang.org/x/exp/slog"
)

func (c *tmsClient) authContext() context.Context {
	return context.WithValue(context.Background(), tmsclient.ContextAPIKeys, map[string]tmsclient.APIKey{
		"Bearer or PrivateToken": {
			Key:    c.cfg.Token,
			Prefix: "PrivateToken",
		},
	})
}

func responseBodyString(resp *http.Response) string {
	if resp == nil || resp.Body == nil {
		return ""
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	return string(b)
}

func logAndWrapAPIError(l *slog.Logger, op, msg string, err error, resp *http.Response) error {
	if body := responseBodyString(resp); body != "" {
		l.Error(msg, "error", err, slog.String("response", body), slog.String("op", op))
	} else {
		l.Error(msg, "error", err, slog.String("op", op))
	}
	return fmt.Errorf("%s: %s: %w", op, msg, err)
}
