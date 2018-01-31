package stringsvc

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/go-kit/kit/endpoint"
)

func (StringService) Uppercase(_ context.Context, s string) (string, error) {
	if s == "" {
		return "", ErrEmpty
	}
	return strings.ToUpper(s), nil
}

// For each method, we define request and response structs
type uppercaseRequest struct {
	S string `json:"s"`
}

type uppercaseResponse struct {
	V   string `json:"v"`
	Err error  `json:"err,omitempty"`
}

func (r uppercaseResponse) Failed() error { return r.Err }

func makeUppercaseEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(uppercaseRequest)
		v, err := svc.Uppercase(ctx, req.S)
		if err != nil {
			return uppercaseResponse{v, err}, nil
		}
		return uppercaseResponse{v, nil}, nil
	}
}

func decodeUppercaseRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request uppercaseRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func (mw logmw) Uppercase(ctx context.Context, s string) (output string, err error) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "uppercase",
			"input", s,
			"output", output,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	output, err = mw.next.Uppercase(ctx, s)
	return
}
