package handler

import (
	"github.com/gorilla/mux"
	"github.com/kimbellG/packet-filter/service/internal/controller"
)

const (
	countPath = "/count"
)

func RegisterCount(router *mux.Router, cont controller.Controller) {
	handler := NewHandler(cont)

	router.HandleFunc(countPath, handler.GetCount).Methods("GET")
	router.HandleFunc(countPath, handler.RefreshCount).Methods("DELETE")
}
