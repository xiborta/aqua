package main

import (
	"fmt"
	"time"

	"github.com/kidoman/embd"
	"github.com/kidoman/embd/convertors/mcp3008"
	_ "github.com/kidoman/embd/host/all"
)

const (
	slot  = 0
	speed = 1000000
	bpw   = 8
	delay = 0
)

func main() {

	if err := embd.InitSPI(); err != nil {
		panic(err)
	}
	defer embd.CloseSPI()

	spiBus := embd.NewSPIBus(embd.SPIMode0, slot, speed, bpw, delay)
	defer spiBus.Close()

	adc := mcp3008.New(mcp3008.SingleMode, spiBus)

	channels := []int{0, 1}

	for {
		time.Sleep(1 * time.Second)

		for _, channel := range channels {
			val, err := adc.AnalogValueAt(channel)
			if err != nil {
				panic(err)
			}
			fmt.Printf("Channel %d value is: %v\n", channel, val)

		}
	}
}
