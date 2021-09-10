package scanner

import (
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

const (
	MainFnName = "filter_main"

	countKey     = "1"
	CountMapName = "pacinfo"

	BlackListMapName = "protocol_blacklist"
)

func InitXDP(netInterface string, source string) (*bcc.Module, error) {
	module := bcc.NewModule(source, []string{
		"-DFUNCNAME=" + MainFnName,
		"-DMAPNAME=" + CountMapName,
		"-DCOUNTKEY=" + countKey,
	})

	fn, err := module.Load(MainFnName, C.BPF_PROG_TYPE_XDP, 1, 65536)
	if err != nil {
		return nil, fmt.Errorf("Failed to load xdp func in ebpf module: %v", err)
	}

	if err := module.AttachXDP(netInterface, fn); err != nil {
		return nil, fmt.Errorf("Failed to attach module to xdp subsystem: %v", err)
	}

	return module, nil
}
