package controller

//Scanner is interface for communicative beetween controller and package scanner
type Scanner interface {
	Count() Count
}

//Count is interface for count controll
type Count interface {
	Get() uint64
	Refresh() error
}
