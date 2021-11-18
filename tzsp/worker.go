package tzsp

import (
	"runtime"

	"github.com/kogonia/x-tzsp/storage"
)

type worker struct {
	ch chan packet
}

var w = worker{ch: make(chan packet, 10_000)}

var st = storage.Init()

func Init() {
	for i := 0; i < runtime.NumCPU(); i++ {
		go workerInit()
	}
	go st.ExpirationCheck(packetKeepAlive)
}

func workerInit() {
	for p := range w.ch {
		if p.ackSN == 0 {
			st.SavePacket(newKey(p), newValue(p))
		} else {
			st.FindPacket(newKey(p), newValue(p))
		}
		st.Done()
	}
}

func (p packet) process() {
	st.Add(1)
	w.ch <- p
	st.Wait()
}
