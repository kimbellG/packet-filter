package controller

type Controller interface {
	GetCount() uint64
	RefreshCount() error
}
