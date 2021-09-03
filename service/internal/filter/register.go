package filter

import (
	_ "embed"
	"fmt"

	"github.com/iovisor/gobpf/bcc"
)

/*
#cgo CFLAGS: -I/usr/include/bcc/compat
#cgo LDFLAGS: -lbcc
#include <bcc/bcc_common.h>
#include <bcc/libbpf.h>
*/
import "C"

//go:embed xdp.c
var source string

const (
	FnName  = "counter"
	MapName = "pacinfo"

	countKey = "1"
)

func InitXDP(netInterface string) (*bcc.Module, error) {
	module := bcc.NewModule(source, []string{
		"-DFUNCNAME=" + FnName,
		"-DMAPNAME=" + MapName,
		"-DCOUNTKEY=" + countKey,
	})

	fn, err := module.Load(FnName, C.BPF_PROG_TYPE_XDP, 1, 65536)
	if err != nil {
		return nil, fmt.Errorf("Failed to load xdp func in ebpf module: %v", err)
	}

	if err := module.AttachXDP(netInterface, fn); err != nil {
		return nil, fmt.Errorf("Failed to attach module to xdp subsystem: %v", err)
	}

	return module, nil
}
