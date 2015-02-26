package main

import (
    "github.com/hybridgroup/gobot/platforms/gpio"
    "github.com/hybridgroup/gobot/platforms/raspi"
)

var physicalPumps []*gpio.DirectPinDriver

func initializePumps() {
    adapter := raspi.NewRaspiAdaptor("raspi")
	physicalPumps = append(physicalPumps, gpio.NewDirectPinDriver(adapter, "pump_1", "8"))
	physicalPumps = append(physicalPumps, gpio.NewDirectPinDriver(adapter, "pump_2", "10"))
}