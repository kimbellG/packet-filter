package handler

import (
	"encoding/json"
	"net/http"

	"github.com/kimbellG/packet-filter/service/pkg/coderror"
	"github.com/kimbellG/packet-filter/service/pkg/codes"
)

func (h *Handler) GetCount(w http.ResponseWriter, r *http.Request) {
	HTTPHandlerWithError(h.getCount).Serve(w, r)
}

func (h *Handler) getCount(w http.ResponseWriter, r *http.Request) error {
	type Response struct {
		Count uint64 `json:"count"`
	}

	count, err := h.cont.GetCount(r.Context())
	if err != nil {
		return coderror.Errorf(err, "controller: %v", err)
	}

	if err := json.NewEncoder(w).Encode(&Response{count}); err != nil {
		return coderror.Newf(codes.JSONEncodeError, "encode json: %v", err)
	}

	return nil
}

func (h *Handler) RefreshCount(w http.ResponseWriter, r *http.Request) {
	HTTPHandlerWithError(h.refreshCount).Serve(w, r)
}

func (h *Handler) refreshCount(w http.ResponseWriter, r *http.Request) error {
	if err := h.cont.RefreshCount(r.Context()); err != nil {
		return coderror.Errorf(err, "refresh count: %v", err)
	}

	return nil
}
