package tms

import "golang.org/x/exp/slog"

const (
	errExtract = "could not extract object by key"
)

func getFromCtx(key string) interface{} {
	value, ok := ctxMgr.GetValue(key)
	if !ok {
		logger.Error(errExtract, slog.String("key", key))
	}

	return value
}

func manipulateOnObjectFromCtx(key string, action func(object interface{})) {
	if object, ok := ctxMgr.GetValue(key); ok {
		action(object)
	} else {
		logger.Error(errExtract, slog.String("key", key))
	}
}
