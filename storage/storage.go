package storage

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

type Storage struct {
	sync.Mutex
	sync.WaitGroup
	m map[SKey]SValue
}

type SKey struct {
	RouterAddr string
	Dst        string
	SN         uint32
	Exp        int64
}

type SValue struct {
	MAC     net.HardwareAddr
	SrcAddr net.IP
	SrcPort uint16
	DstAddr net.IP
	DstPort uint16
}

func Init() Storage {
	return Storage{
		Mutex:     sync.Mutex{},
		WaitGroup: sync.WaitGroup{},
		m:         make(map[SKey]SValue, 10_000),
	}
}
func (st *Storage) SavePacket(k SKey, v SValue) {
	st.Lock()
	defer st.Unlock()
	st.m[k] = v
}

func (st *Storage) FindPacket(k SKey, v SValue) {
	if len(st.m) == 0 {
		return
	}

	for i, sm := range st.m {
		if i.RouterAddr == k.RouterAddr && i.SN+1 == k.SN && i.Dst == k.Dst {
			fmt.Printf("[%s] %v %s:%d %s:%d %s:%d\n",
				k.RouterAddr,
				strings.Replace(sm.MAC.String(), ":", "-", -1),
				sm.SrcAddr, sm.SrcPort,
				v.DstAddr, v.DstPort,
				v.SrcAddr, v.SrcPort)
			st.Lock()
			delete(st.m, i)
			st.Unlock()
			return
		}
	}
}

func (st *Storage) ExpirationCheck(checkInterval time.Duration) {
	for range time.Tick(checkInterval) {
		if len(st.m) > 0 {
			for k := range st.m {
				if k.Exp >= time.Now().Unix() {
					st.Lock()
					delete(st.m, k)
					st.Unlock()
				}
			}
		}
	}
}
