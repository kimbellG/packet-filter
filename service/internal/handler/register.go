package handler

import (
	"github.com/gorilla/mux"
	"github.com/kimbellG/packet-filter/service/internal/controller"
)

const (
	countPath  = "/count"
	filterPath = "/filter/{proto}"
)

func RegisterCount(router *mux.Router, cont controller.Controller) {
	handler := NewHandler(cont)

	router.HandleFunc(countPath, handler.GetCount).Methods("GET")
	router.HandleFunc(countPath, handler.RefreshCount).Methods("DELETE")

	router.HandleFunc(filterPath, handler.BlockL2Proto).Methods("GET")
	router.HandleFunc(filterPath, handler.UnBlockL2Proto).Methods("DELETE")
}
