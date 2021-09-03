package controller

import (
	"github.com/kimbellG/packet-filter/service/pkg/coderror"
)

type interactor struct {
	sc Scanner
}

func NewController(sc Scanner) Controller {
	return &interactor{
		sc: sc,
	}
}

func (i *interactor) GetCount() uint64 {
	return i.sc.Count().Get()
}

func (i *interactor) RefreshCount() error {
	if err := i.sc.Count().Refresh(); err != nil {
		return coderror.Errorf(err, "scanner: %v", err)
	}

	return nil
}
