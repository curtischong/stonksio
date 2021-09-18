package request

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

func (handler *RequestHandler) sendStandardHeaders(
	w http.ResponseWriter,
) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
}
func (handler *RequestHandler) sendStatusOK(
	w http.ResponseWriter,
) {
	handler.sendStandardHeaders(w)
	w.WriteHeader(http.StatusOK)
}

func (handler *RequestHandler) sendInternalServerError(
	w http.ResponseWriter,
	err error,
) {
	handler.sendStandardHeaders(w)
	log.Error(err)
	w.WriteHeader(http.StatusInternalServerError)
	return
}

func (handler *RequestHandler) sendStatusBadRequest(
	w http.ResponseWriter,
	err error,
) {
	handler.sendStandardHeaders(w)
	log.Error(err)
	w.WriteHeader(http.StatusBadRequest)
	return
}
