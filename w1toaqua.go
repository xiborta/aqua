package main

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dustin/go-coap"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

const devicesPath = "/sys/bus/w1/devices/"

func main() {

	var publishingURL string

	if len(os.Args) > 1 {
		publishingURL = os.Args[1]
		log.Println("Will send data to: " + publishingURL)
	}

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
					if publishingURL != "" {
						sendMeasure(publishingURL, deviceID, "temperature", temp)
					} else {
						log.Println(deviceID + ": " + temp + " C")
					}
				} else {
					log.Print("Cannot read data: " + sample + " from device: " + deviceID)
				}
			}
		}

		time.Sleep(time.Second)
	}
}

func sendMeasure(url string, deviceID string, measure string, value string) {

	req := coap.Message{
		Type:    coap.NonConfirmable,
		Code:    coap.POST,
		Payload: []byte(value),
	}

	path := "/aqua/" + deviceID + "/" + measure + "/"

	req.SetPathString(path)

	c, err := coap.Dial("udp", url)
	if err != nil {
		log.Fatal("Error dialing: %v", err)
	}

	rv, err := c.Send(req)
	if err != nil {
		log.Fatal("Error sending request: %v", err)
	}

	if rv != nil {
		log.Print("Response payload: %s", rv.Payload)
	}
}
