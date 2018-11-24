package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

const devicesPath = "/sys/bus/w1/devices/"

var (
	certFile = flag.String("cert", "someCertFile", "A PEM eoncoded certificate file.")
	keyFile  = flag.String("key", "someKeyFile", "A PEM encoded private key file.")
	topic    = flag.String("topic", "the aws iot publishing endpoint", "the aws iot publishing endpoint URL.")
)

func main() {

	flag.Parse()

	for {

		var deviceIDs []string

		files, err := ioutil.ReadDir(devicesPath)
		check(err)
		for _, file := range files {
			if strings.HasPrefix(file.Name(), "28-") {
				deviceIDs = append(deviceIDs, file.Name())
			}
		}

		for _, deviceID := range deviceIDs {

			data, err := ioutil.ReadFile(devicesPath + "/" + deviceID + "/w1_slave")
			check(err)

			sample := string(data)

			if strings.Contains(sample, "NO") {
				log.Print(sample)
			} else {
				pos := strings.Index(sample, "t=")
				if pos >= 0 {
					value := sample[pos+2 : pos+5]
					f, err := strconv.ParseFloat(value, 64)
					check(err)
					celsius := f * 0.1

					temp := strconv.FormatFloat(celsius, 'f', 1, 64)

					sendMeasure(deviceID, "temperature", temp)

					log.Println(deviceID + ": " + temp + " C")
				} else {
					log.Print("Cannot read data: " + sample + " from device: " + deviceID)
				}
			}
		}

		time.Sleep(time.Minute)
	}
}

func sendMeasure(deviceID string, measure string, value string) {

	cert, err := tls.LoadX509KeyPair(*certFile, *keyFile)
	if err != nil {
		log.Fatal(err)
	}

	// Setup HTTPS client
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: transport}

	values := map[string]string{"ts": "now", measure: value}

	jsonValue, _ := json.Marshal(values)

	// POST to aws topic
	resp, err := client.Post(*topic, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Print(err)
	}
	defer resp.Body.Close()
}
