package controller

import (
	"context"
	"net"

	"github.com/kimbellG/packet-filter/service/pkg/coderror"
)

func (i *interactor) BlockSubnet(ctx context.Context, addr net.IP) (uint64, error) {
	count, err := i.sc.Filter().BlockSubnet(ctx, addr)
	if err != nil {
		return count, coderror.Errorf(err, "xdp filter: %v", err)
	}

	return count, nil
}

func (i *interactor) UnblockSubnet(ctx context.Context, addr net.IP) error {
	if err := i.sc.Filter().UnblockSubnet(ctx, addr); err != nil {
		return coderror.Errorf(err, "xdp filter: %v", err)
	}

	return nil
}
