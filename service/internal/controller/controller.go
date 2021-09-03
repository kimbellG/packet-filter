package controller

import "context"

type Controller interface {
	GetCount(ctx context.Context) uint64
	RefreshCount(ctx context.Context) error
}
