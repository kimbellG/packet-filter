package controller

import (
	"context"

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

func (i *interactor) GetCount(ctx context.Context) (uint64, error) {
	count, err := i.sc.Count().Get(ctx)
	if err != nil {
		return 0, coderror.Errorf(err, "get count from scanner: %v", err)
	}

	return count, nil
}

func (i *interactor) RefreshCount(ctx context.Context) error {
	if err := i.sc.Count().Refresh(ctx); err != nil {
		return coderror.Errorf(err, "scanner: %v", err)
	}

	return nil
}
