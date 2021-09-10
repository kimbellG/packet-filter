package handler

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kimbellG/packet-filter/service/internal/controller"
	"github.com/kimbellG/packet-filter/service/pkg/coderror"
	"github.com/kimbellG/packet-filter/service/pkg/codes"
)

var protoToController = map[string]controller.IPProto{
	"IP":       controller.IP,
	"ICMP":     controller.ICMP,
	"IPIP":     controller.IPIP,
	"TCP":      controller.TCP,
	"EGP":      controller.EGP,
	"PUP":      controller.PUP,
	"UPD":      controller.UDP,
	"TP":       controller.TP,
	"DCCP":     controller.DCCP,
	"IPV6":     controller.IPV6,
	"RSVP":     controller.RSVP,
	"GRE":      controller.GRE,
	"ESP":      controller.ESP,
	"AH":       controller.AH,
	"MTP":      controller.MTP,
	"BEETPH":   controller.BEETPH,
	"ENCAP":    controller.ENCAP,
	"PIM":      controller.PIM,
	"COMP":     controller.COMP,
	"SCTP":     controller.SCTP,
	"UDPLITE":  controller.UDPLITE,
	"MPLS":     controller.MPLS,
	"EHTERNET": controller.ETHERNET,
}

func parseProto(proto string) (controller.IPProto, error) {
	protoID, ok := protoToController[proto]
	if !ok {
		return protoID, coderror.New(codes.InvalidProtocolInIP, fmt.Errorf("invalid protocol name: %s", proto))
	}

	return protoID, nil
}

func (h *Handler) blockL2Proto(w http.ResponseWriter, r *http.Request) error {
	proto := mux.Vars(r)["proto"]

	protoID, err := parseProto(proto)
	if err != nil {
		return coderror.Errorf(err, "parsing protocol path: %v", err)
	}

	if err := h.cont.BlockL2Protocol(protoID); err != nil {
		return coderror.Errorf(err, "controller: %v", err)
	}

	return nil
}

func (h *Handler) unBlockL2Proto(w http.ResponseWriter, r *http.Request) error {
	proto := mux.Vars(r)["proto"]

	protoID, err := parseProto(proto)
	if err != nil {
		return coderror.Errorf(err, "parsing protocol path: %v", err)
	}

	if err := h.cont.UnBlockL2Protocol(protoID); err != nil {
		return coderror.Errorf(err, "controller: %v", err)
	}

	return nil
}

func (h *Handler) BlockL2Proto(w http.ResponseWriter, r *http.Request) {
	HTTPHandlerWithError(h.blockL2Proto).Serve(w, r)
}

func (h *Handler) UnBlockL2Proto(w http.ResponseWriter, r *http.Request) {
	HTTPHandlerWithError(h.unBlockL2Proto).Serve(w, r)
}
