package scanner

import (
	"encoding/binary"
	"errors"

	"github.com/iovisor/gobpf/bcc"
	"github.com/kimbellG/packet-filter/service/pkg/coderror"
	"github.com/kimbellG/packet-filter/service/pkg/codes"
)

type XDPScanner struct {
	count *XDPCounter
}

func NewXDPScanner(module *bcc.Module) *XDPScanner {
	return &XDPScanner{
		count: NewXDPCounter(bcc.NewTable(module.TableId(MapName), module)),
	}
}

func encodeKey(key uint8) []byte {
	buf := make([]byte, 8)
	binary.PutUvarint(buf, uint64(key))

	return buf[:1]
}

func decodeUIntValue(value []byte) (uint64, error) {
	if value == nil {
		return 0, coderror.New(codes.BCCNilValueFromTableError, errors.New("nil value"))
	}

	return bcc.GetHostByteOrder().Uint64(value), nil
}
