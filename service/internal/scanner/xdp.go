package scanner

import (
	"encoding/binary"
	"errors"

	"github.com/iovisor/gobpf/bcc"
	"github.com/kimbellG/packet-filter/service/internal/controller"
	"github.com/kimbellG/packet-filter/service/pkg/coderror"
	"github.com/kimbellG/packet-filter/service/pkg/codes"
)

type XDPScanner struct {
	count  *XDPCounter
	filter *XDPFilter
}

func NewXDPScanner(module *bcc.Module) *XDPScanner {
	return &XDPScanner{
		count:  NewXDPCounter(bcc.NewTable(module.TableId(CountMapName), module)),
		filter: NewXDPFilter(bcc.NewTable(module.TableId(BlackListMapName), module)),
	}
}

func (xs *XDPScanner) Count() controller.Count {
	return xs.count
}

func (xs *XDPScanner) Filter() controller.Filter {
	return xs.filter

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
