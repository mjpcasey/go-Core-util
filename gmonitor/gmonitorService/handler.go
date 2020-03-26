package gmonitorService

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func handlerHealthCheck(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Success"))
}