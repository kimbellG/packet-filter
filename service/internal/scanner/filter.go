package scanner

import (
	"context"
	"encoding/binary"
	"fmt"
	"net"

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

	BGID := make([]byte, 2)
	bcc.GetHostByteOrder().PutUint16(BGID, id)

	return BGID, nil
}

type XDPFilter struct {
	protoBlacklist  *bcc.Table
	subnetBlacklist *bcc.Table
}

func NewXDPFilter(blackProto, blackSubnet *bcc.Table) *XDPFilter {
	return &XDPFilter{
		protoBlacklist:  blackProto,
		subnetBlacklist: blackSubnet,
	}
}

func (xf *XDPFilter) Block(protocol controller.IPProto) error {
	id, err := getIDOfProtocol(protocol)
	if err != nil {
		return coderror.Errorf(err, "get ip id of protocol: %v", err)
	}

	value := make([]byte, 22)
	bcc.GetHostByteOrder().PutUint16(value, 1)

	if err := xf.protoBlacklist.Set(id, value); err != nil {
		return coderror.Newf(codes.BCCSetToTableError, "set new protocol to blacklist: %v", err)
	}

	return nil
}

func (xf *XDPFilter) UnBlock(protocol controller.IPProto) error {
	id, err := getIDOfProtocol(protocol)
	if err != nil {
		return coderror.Errorf(err, "get id of protocol: %v", err)
	}

	if err := xf.protoBlacklist.Delete(id); err != nil {
		return coderror.Newf(codes.BCCDeleteFromTableError, "delete protocol from talbe: %v", err)
	}

	return nil
}

func (xf *XDPFilter) BlockSubnet(ctx context.Context, addr net.IP) (uint64, error) {
	select {
	case <-ctx.Done():
		// TODO: Set a valid error code
		return 0, coderror.Newf(codes.Unknown, "stopped by context")
	default:
		return xf.blockSubnet(addr)
	}
}

func (xf *XDPFilter) blockSubnet(addr net.IP) (uint64, error) {
	addrBytes := xf.addrInBytes(addr)

	count, err := xf.subnetBlacklist.Get(addrBytes)
	if err != nil {
		return xf.addSubnetToBlacklist(addrBytes)
	}

	return binary.BigEndian.Uint64(count), nil
}

func (xf *XDPFilter) addrInBytes(addr net.IP) []byte {
	addrBytes := make([]byte, 4)

	binary.BigEndian.PutUint32(addrBytes, binary.BigEndian.Uint32(addr[12:16]))

	return addrBytes
}

func (xf *XDPFilter) addSubnetToBlacklist(addr []byte) (uint64, error) {
	count := make([]byte, 8)

	binary.BigEndian.PutUint32(count, 0)

	if err := xf.subnetBlacklist.Set(addr, count); err != nil {
		return 0, coderror.Newf(codes.BCCSetToTableError, "set new subnet to the blacklist: %v", err)
	}

	return 0, nil
}

func (xf *XDPFilter) UnblockSubnet(ctx context.Context, addr net.IP) error {
	select {
	case <-ctx.Done():
		// TODO: Set a valid error code
		return coderror.Newf(codes.Unknown, "stopped by context")
	default:
		return xf.unblockSubnet(addr)
	}
}

func (xf *XDPFilter) unblockSubnet(addr net.IP) error {
	addrBytes := xf.addrInBytes(addr)

	if err := xf.subnetBlacklist.Delete(addrBytes); err != nil {
		return coderror.Newf(codes.BCCDeleteFromTableError, "delete from bcc table: %v", err)
	}

	return nil
}
