package controller

import (
	"context"
	"net"
)

//Scanner is interface for communicative beetween controller and package scanner
type Scanner interface {
	Count() Count
	Filter() Filter
}

type IPProto int

const (
	IP IPProto = iota
	ICMP
	IGMP
	IPIP
	TCP
	EGP
	PUP
	UDP
	IDP
	TP
	DCCP
	IPV6
	RSVP
	GRE
	ESP
	AH
	MTP
	BEETPH
	ENCAP
	PIM
	COMP
	SCTP
	UDPLITE
	MPLS
	ETHERNET
)

type Filter interface {
	Block(protocol IPProto) error
	UnBlock(protocol IPProto) error

	BlockSubnet(ctx context.Context, addr net.IP) (uint64, error)
	UnblockSubnet(ctx context.Context, addr net.IP) error
}

//Count is interface for count controll
type Count interface {
	Get(ctx context.Context) (uint64, error)
	Refresh(ctx context.Context) error
}
