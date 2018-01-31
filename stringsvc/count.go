package stringsvc

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-kit/kit/endpoint"
)

func (StringService) Count(_ context.Context, s string) int {
	return len(s)
}

type countRequest struct {
	S string `json:"s"`
}

type countResponse struct {
	V int `json:"v"`
}

func makeCountEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(countRequest)
		v := svc.Count(ctx, req.S)
		return countResponse{v}, nil
	}
}

func decodeCountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request countRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func (mw logmw) Count(ctx context.Context, s string) (n int) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "count",
			"input", s,
			"n", n,
			"took", time.Since(begin),
		)
	}(time.Now())

	n = mw.next.Count(ctx, s)
	return
}
