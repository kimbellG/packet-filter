package handler

import (
	"fmt"
	"net/http"

	"github.com/kimbellG/packet-filter/service/pkg/coderror"
	"github.com/kimbellG/packet-filter/service/pkg/codes"
)

type HTTPHandlerWithError func(w http.ResponseWriter, r *http.Request) error

func (fn HTTPHandlerWithError) Serve(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		msg, code := decodeError(err)
		http.Error(w, msg, code)
	}
}

var httpCode = map[int]int{
	codes.Unknown.Int():                   http.StatusInternalServerError,
	codes.JSONEncodeError.Int():           http.StatusInternalServerError,
	codes.InvalidProtocolInIP.Int():       http.StatusNotFound,
	codes.BCCGetFromTableError.Int():      http.StatusInternalServerError,
	codes.BCCSetToTableError.Int():        http.StatusInternalServerError,
	codes.BCCDeleteFromTableError.Int():   http.StatusInternalServerError,
	codes.BCCNilValueFromTableError.Int(): http.StatusInternalServerError,
}

func decodeError(err error) (string, int) {
	cerr, ok := coderror.As(err)
	if !ok {
		return fmt.Sprintf("Failed http request: %v", err.Error()), http.StatusInternalServerError
	}

	hcode, ok := httpCode[cerr.Code().Int()]
	if !ok {
		return fmt.Sprintf("Failed http request: status %s: %v", cerr.Code().String(), cerr.Error()), http.StatusInternalServerError
	}

	return fmt.Sprintf("Failed http request: status %s: %v", cerr.Code().String(), cerr.Error()), hcode
}
