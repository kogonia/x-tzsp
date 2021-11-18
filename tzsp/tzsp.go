package tzsp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	packetKeepAlive = 10 * time.Second

	tzspHeaderLength    = 5
	ethHeaderLength     = 14
	ipv4HeaderMinLength = 20

	ipv4EtherType uint16 = 0x0800

	tcpProto = 6
	udpProto = 17
)

var (
	errHeaderNotTZSP = errors.New("not TZSP")
	// errHeaderTooShort     = errors.New("header too short")
	errNotIPv4            = errors.New("not IPv4")
	errNotMAC             = errors.New("failed to parse MAC from ethernet header")
	errIPv4HeaderTooShort = errors.New("ip header less 20 byte")
)

func Parse(b []byte, routerAddr net.IP) error {
	b, err := checkTZSPHeader(b)
	if err != nil {
		return err
	}

	p := packet{routerAddr: routerAddr}

	b, err = parseEthHeader(b, &p)
	if err != nil {
		return err
	}

	err = parseIPv4Header(b, &p)
	if err != nil {
		return err
	}
	p.exp = time.Now().Add(packetKeepAlive).Unix()
	go p.process()

	return nil
}

func checkTZSPHeader(b []byte) ([]byte, error) {
	if bytes.Compare(b[:tzspHeaderLength], []byte{1, 0, 0, 1, 1}) != 0 {
		return nil, errHeaderNotTZSP
	}
	return b[tzspHeaderLength:], nil
}

func parseEthHeader(bOld []byte, p *packet) (b []byte, err error) {
	ethHeader := bOld[:ethHeaderLength]
	b = bOld[ethHeaderLength:]

	// dstMac, err = toMAC(ethHeader[:6])
	// if err != nil {
	// 	return b, errNotMAC
	// }
	p.mac, err = toMAC(ethHeader[6:12])
	if err != nil {
		return b, errNotMAC
	}

	etherType := binary.BigEndian.Uint16(ethHeader[12:14])
	if etherType != ipv4EtherType {
		return b, errNotIPv4
	}

	return
}

func parseIPv4Header(b []byte, p *packet) error {
	ipv4HeaderSize := b[0] & 0x0F * 4 // header size in byte
	if ipv4HeaderSize < ipv4HeaderMinLength {
		return errIPv4HeaderTooShort
	}

	ipHeader := b[:ipv4HeaderSize]
	b = b[ipv4HeaderSize:]

	protocol := ipHeader[9]

	p.srcAddr = toIP(ipHeader[12:16])

	p.dstAddr = toIP(ipHeader[16:20])
	if protocol == tcpProto || protocol == udpProto {
		p.srcPort = binary.BigEndian.Uint16(b[:2])
		p.dstPort = binary.BigEndian.Uint16(b[2:4])
	}

	if protocol == tcpProto {
		p.sn = binary.BigEndian.Uint32(b[4:8])
		p.ackSN = binary.BigEndian.Uint32(b[8:12])
	}
	return nil
}

func toIP(octets []byte) net.IP {
	var addr string
	for i := range octets {
		addr += strconv.Itoa(int(octets[i])) + "."
	}
	return net.ParseIP(addr[:len(addr)-1])
}

func toMAC(octets []byte) (net.HardwareAddr, error) {
	return net.ParseMAC(strings.Replace(fmt.Sprintf("% x", octets), " ", ":", -1))
}
