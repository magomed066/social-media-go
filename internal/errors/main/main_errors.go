package mainErrors

import (
	"errors"
	"net/http"
	"social-media-app/internal/json"

	"charm.land/log/v2"
)

var ErrNotFound = errors.New("record not found")

func InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Errorf("internal server error %s, path: %s, error: %s", r.Method, r.URL.Path, err)

	json.WriteJSONError(w, http.StatusInternalServerError, "internal server error")
}

func BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Errorf("bad request error %s, path: %s, error: %s", r.Method, r.URL.Path, err)

	json.WriteJSONError(w, http.StatusBadRequest, err.Error())
}

func NotFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Errorf("not found error %s, path: %s, error: %s", r.Method, r.URL.Path, err)

	json.WriteJSONError(w, http.StatusNotFound, "record not found")
}
