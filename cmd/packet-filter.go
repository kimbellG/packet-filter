package main

import (
	_ "embed"

	"github.com/kimbellG/packet-filter/service"
)

//go:embed xdp.c
var source string

func main() {
	service.Run(source)
}
