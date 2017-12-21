package main

import (
	"github.com/dustin/go-coap"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	aquaTemp = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "temperature_celsius",
		Help: "Current temperature of the aquarium in degree Celsius.",
	})
)

func init() {
	// Metrics have to be registered to be exposed:
	prometheus.MustRegister(aquaTemp)
}

func handleAqua(l *net.UDPConn, a *net.UDPAddr, m *coap.Message) *coap.Message {

	log.Printf("Reading: %s %v", strings.Join(m.Path(), "/"), string(m.Payload))

	f, _ := strconv.ParseFloat(string(m.Payload), 64)

	aquaTemp.Set(f)

	return nil
}

func startCoap() {

	mux := coap.NewServeMux()
	mux.Handle("/aqua/", coap.FuncHandler(handleAqua))

	log.Printf("installing /aqua/ handler ...")

	log.Fatal(coap.ListenAndServe("udp", ":5683", mux))
}

func main() {

	go startCoap()

	// The Handler function provides a default handler to expose metrics
	// via an HTTP server. "/metrics" is the usual endpoint for that.
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))

}
