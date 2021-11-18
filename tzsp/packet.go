package tzsp

import (
	"net"
	"strconv"

	"github.com/kogonia/x-tzsp/storage"
)

type packet struct {
	routerAddr net.IP
	mac        net.HardwareAddr
	srcAddr    net.IP
	srcPort    uint16
	natSrcAddr net.IP
	natSrcPort uint16
	dstAddr    net.IP
	dstPort    uint16
	exp        int64
	sn, ackSN  uint32
}

type StorageKey struct {
	RouterAddr string
	Dst        string
	SN         uint32
	Exp        int64
}

type StorageValue struct {
	MAC     net.HardwareAddr
	SrcAddr net.IP
	SrcPort uint16
	DstAddr net.IP
	DstPort uint16
}

func newKey(p packet) (k storage.SKey) {
	if p.ackSN == 0 {
		// TCP Syn
		k = storage.SKey{
			RouterAddr: p.routerAddr.String(),
			Dst:        p.dstAddr.String() + ":" + strconv.Itoa(int(p.dstPort)),
			SN:         p.sn,
			Exp:        p.exp,
		}
	} else {
		// TCP Ack
		k = storage.SKey{
			RouterAddr: p.routerAddr.String(),
			Dst:        p.srcAddr.String() + ":" + strconv.Itoa(int(p.srcPort)),
			SN:         p.ackSN,
		}
	}
	return
}

func newValue(p packet) storage.SValue {
	return storage.SValue{
		MAC:     p.mac,
		SrcAddr: p.srcAddr,
		SrcPort: p.srcPort,
		DstAddr: p.dstAddr,
		DstPort: p.dstPort,
	}
}
