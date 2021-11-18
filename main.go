package main

import (
	"log"

	"github.com/kogonia/x-tzsp/udpSocket"
)

func main() {
	if err := udpSocket.Start("0.0.0.0", "5514"); err != nil {
		log.Fatal(err)
	}
}
