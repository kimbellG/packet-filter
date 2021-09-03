package handler

import (
	"github.com/kimbellG/packet-filter/service/internal/controller"
)

type Handler struct {
	cont controller.Controller
}

func NewHandler(cont controller.Controller) *Handler {
	return &Handler{
		cont: cont,
	}
}
