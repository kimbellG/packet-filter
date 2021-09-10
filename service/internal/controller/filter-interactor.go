package controller

import "github.com/kimbellG/packet-filter/service/pkg/coderror"

func (i *interactor) BlockL2Protocol(proto IPProto) error {
	if err := i.sc.Filter().Block(proto); err != nil {
		return coderror.Errorf(err, "scanner: %v", err)
	}

	return nil
}

func (i *interactor) UnBlockL2Protocol(proto IPProto) error {
	if err := i.sc.Filter().UnBlock(proto); err != nil {
		return coderror.Errorf(err, "scanner: %v", err)
	}

	return nil
}
