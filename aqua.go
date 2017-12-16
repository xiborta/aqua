
package main

import (
	"log"
	"net"
	"strings"

	"github.com/dustin/go-coap"
)


func handleAqua(l *net.UDPConn, a *net.UDPAddr, m *coap.Message) *coap.Message {

	log.Printf("Reading: %s %v", strings.Join(m.Path(), "/"), string(m.Payload))	
	
	return nil
}

func main() {
	
  mux := coap.NewServeMux()
	mux.Handle("/aqua/", coap.FuncHandler(handleAqua))

  log.Printf("installing /aqua/ handler ...")

	log.Fatal(coap.ListenAndServe("udp", ":5683", mux))
}
