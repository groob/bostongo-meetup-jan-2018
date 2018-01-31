package stringsvc

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

type Endpoints struct {
	UppercaseEndpoint endpoint.Endpoint
	CountEndpoint     endpoint.Endpoint
}

func MakeServerEndpoints(svc Service) Endpoints {
	return Endpoints{
		UppercaseEndpoint: makeUppercaseEndpoint(svc),
		CountEndpoint:     makeCountEndpoint(svc),
	}
}

func MakeHTTPHandler(r *mux.Router, e Endpoints, logger log.Logger) {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(ErrorEncoder),
		httptransport.ServerErrorLogger(logger),
	}

	// POST     /v1/uppercase			uppercase a string
	// POST     /v1/count				count the length of a string

	r.Methods("POST").Path("/v1/uppercase").Handler(httptransport.NewServer(
		e.UppercaseEndpoint,
		decodeUppercaseRequest,
		EncodeJSONResponse,
		options...,
	))

	r.Methods("POST").Path("/v1/count").Handler(httptransport.NewServer(
		e.CountEndpoint,
		decodeCountRequest,
		EncodeJSONResponse,
		options...,
	))
}

// failer is an interface that should be implemented by response types.
// Response encoders can check if responses are Failer, and if so if they've
// failed, and if so encode them using a separate write path based on the error.
type failer interface {
	Failed() error
}

func EncodeJSONResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if f, ok := response.(failer); ok && f.Failed() != nil {
		ErrorEncoder(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if headerer, ok := response.(httptransport.Headerer); ok {
		for k := range headerer.Headers() {
			w.Header().Set(k, headerer.Headers().Get(k))
		}
	}
	code := http.StatusOK
	if sc, ok := response.(httptransport.StatusCoder); ok {
		code = sc.StatusCode()
	}
	w.WriteHeader(code)

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(response)
}

func ErrorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	errMap := map[string]interface{}{"error": err.Error()}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")

	if headerer, ok := err.(httptransport.Headerer); ok {
		for k := range headerer.Headers() {
			w.Header().Set(k, headerer.Headers().Get(k))
		}
	}

	code := http.StatusInternalServerError
	if sc, ok := err.(httptransport.StatusCoder); ok {
		code = sc.StatusCode()
	}
	w.WriteHeader(code)

	enc.Encode(errMap)
}
