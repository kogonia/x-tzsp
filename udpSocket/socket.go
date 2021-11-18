package udpSocket

import (
	"log"
	"net"

	"github.com/kogonia/x-tzsp/tzsp"
)

const (
	tzspHeaderLength    = 5
	ethHeaderLength     = 14
	ipv4HeaderMinLength = 20
)

var buf = make([]byte, 4096)

func Start(addr, port string) error {
	conn, err := listen(addr + ":" + port)
	if err != nil {
		return err
	}

	go tzsp.Init()

	handleData(conn)

	return nil
}

func listen(addr string) (conn *net.UDPConn, err error) {
	udpAddr, err := net.ResolveUDPAddr("udp4", addr)
	if err != nil {
		return
	}

	conn, err = net.ListenUDP("udp", udpAddr)
	if err != nil {
		return
	}
	// log.Printf("listening on %s", udpAddr)

	return
}

func handleData(conn *net.UDPConn) {
	for {
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Println("failed to read from UDP")
		}
		if n < tzspHeaderLength+ethHeaderLength+ipv4HeaderMinLength {
			log.Printf("packet too short: %v", buf[:n])
			continue
		}
		go func(b []byte, routerAddr net.IP) {
			if err := tzsp.Parse(b, routerAddr); err != nil {
				log.Printf("failed to parse data from '%v': %v ", routerAddr, err)
			}
		}(buf[0:n], addr.IP)
	}
}
