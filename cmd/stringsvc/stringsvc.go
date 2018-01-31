package main

import (
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"github.com/groob/bostongo-meetup-jan-2018/stringsvc"
)

func main() {
	logger := log.NewLogfmtLogger(os.Stderr)

	var svc stringsvc.Service
	svc = stringsvc.New()
	svc = stringsvc.LoggingMiddleware(logger)(svc)

	endpoints := stringsvc.MakeServerEndpoints(svc)

	r := mux.NewRouter()
	stringsvc.MakeHTTPHandler(r, endpoints, logger)

	http.ListenAndServe("8080", r)
}
