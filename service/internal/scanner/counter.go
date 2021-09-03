package scanner

import (
	"context"
	"encoding/binary"

	"github.com/iovisor/gobpf/bcc"
	"github.com/kimbellG/packet-filter/service/internal/scanner/tablekey"
	"github.com/kimbellG/packet-filter/service/pkg/coderror"
	"github.com/kimbellG/packet-filter/service/pkg/codes"
)

type XDPCounter struct {
	countKey         uint8
	packageInfoTable *bcc.Table
}

func NewXDPCounter(packageInfoTable *bcc.Table) *XDPCounter {
	return &XDPCounter{
		countKey:         tablekey.Count,
		packageInfoTable: packageInfoTable,
	}
}

func (xs *XDPCounter) Get(ctx context.Context) (uint64, error) {
	countBytes, err := xs.getValueFromTable(xs.countKey)
	if err != nil {
		return 0, coderror.Errorf(err, "get count from table: %v", err)
	}

	count, err := decodeUIntValue(countBytes)
	if err != nil {
		return 0, coderror.Errorf(err, "decoding count from table: %v", err)
	}

	return count, nil
}

func (xs *XDPCounter) getValueFromTable(key uint8) ([]byte, error) {
	count, err := xs.packageInfoTable.Get(encodeKey(key))
	if err != nil {
		return nil, coderror.New(codes.BCCGetFromTableError, err)
	}

	return count, nil
}

func (xs *XDPCounter) Refresh(ctx context.Context) error {
	var newValue uint64

	if err := xs.setValueToTable(xs.countKey, newValue); err != nil {
		return coderror.Errorf(err, "set value to bpf table: %v", err)
	}

	return nil
}

func (xs *XDPCounter) setValueToTable(key uint8, value uint64) error {
	if err := xs.packageInfoTable.Set(encodeKey(key), encodeValueToBytes(value)); err != nil {
		return coderror.New(codes.BCCSetToTableError, err)
	}

	return nil
}

func encodeValueToBytes(value uint64) []byte {
	buf := make([]byte, 8)
	binary.PutUvarint(buf, value)

	return buf
}
