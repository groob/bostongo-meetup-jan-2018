package stringsvc

import "github.com/go-kit/kit/log"

type Middleware func(next Service) Service

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return logmw{logger, next}
	}
}

type logmw struct {
	logger log.Logger
	next   Service
}
