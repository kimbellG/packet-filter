package controller

import "context"

//Scanner is interface for communicative beetween controller and package scanner
type Scanner interface {
	Count() Count
}

//Count is interface for count controll
type Count interface {
	Get(ctx context.Context) uint64
	Refresh(ctx context.Context) error
}
