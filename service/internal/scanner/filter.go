package scanner

import (
	"encoding/binary"
	"fmt"

	"github.com/iovisor/gobpf/bcc"
	"github.com/kimbellG/packet-filter/service/internal/controller"
	"github.com/kimbellG/packet-filter/service/pkg/coderror"
	"github.com/kimbellG/packet-filter/service/pkg/codes"
)

type IPProto int

var idsOfProtos = map[controller.IPProto]uint16{
	controller.IP:       0,
	controller.ICMP:     1,
	controller.IGMP:     2,
	controller.IPIP:     4,
	controller.TCP:      6,
	controller.EGP:      8,
	controller.PUP:      12,
	controller.UDP:      17,
	controller.IDP:      22,
	controller.TP:       29,
	controller.DCCP:     33,
	controller.IPV6:     41,
	controller.RSVP:     46,
	controller.GRE:      47,
	controller.ESP:      50,
	controller.AH:       51,
	controller.MTP:      92,
	controller.BEETPH:   94,
	controller.ENCAP:    98,
	controller.PIM:      103,
	controller.COMP:     108,
	controller.SCTP:     132,
	controller.UDPLITE:  136,
	controller.MPLS:     137,
	controller.ETHERNET: 143,
}

func getIDOfProtocol(proto controller.IPProto) ([]byte, error) {
	id, ok := idsOfProtos[proto]
	if !ok {
		return nil, coderror.New(codes.InvalidProtocolInIP, fmt.Errorf("protocol with id(%v) don't encapsulated in ip protocol", proto))
	}

	BGID := make([]byte, 16)
	binary.BigEndian.PutUint16(BGID, id)

	return BGID, nil
}

type XDPFilter struct {
	blacklist *bcc.Table
}

func NewXDPFilter(blacklist *bcc.Table) *XDPFilter {
	return &XDPFilter{
		blacklist: blacklist,
	}
}

func (xf *XDPFilter) Block(protocol controller.IPProto) error {
	id, err := getIDOfProtocol(protocol)
	if err != nil {
		return coderror.Errorf(err, "get ip id of protocol: %v", err)
	}

	value := make([]byte, 16)
	binary.BigEndian.PutUint16(value, 1)

	if err := xf.blacklist.Set(id, value); err != nil {
		return coderror.Newf(codes.BCCSetToTableError, "set new protocol to blacklist: %v", err)
	}

	return nil
}

func (xf *XDPFilter) UnBlock(protocol controller.IPProto) error {
	id, err := getIDOfProtocol(protocol)
	if err != nil {
		return coderror.Errorf(err, "get id of protocol: %v", err)
	}

	if err := xf.blacklist.Delete(id); err != nil {
		return coderror.Newf(codes.BCCDeleteFromTableError, "delete protocol from talbe: %v", err)
	}

	return nil
}
