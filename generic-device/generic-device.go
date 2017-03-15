package main

import (
	"log"
	"os"
	"time"

	// DBus helper library from DeviceHive
	// Not an absolute requirement, it just makes the use of DBus from Go easier
	"github.com/devicehive/IoT-framework/godbus-helpers/cloud"

	// Go library providing system information
	"github.com/shirou/gopsutil/cpu"
)

func main() {
	cloud, err := cloud.NewDbusForComDevicehiveCloud()
	if err != nil {
		log.Panic(err)
	}

	h, _ := os.Hostname()

	for {
		c, err := cpu.CPUPercent(time.Second, false)
		if err != nil {
			log.Panic(err)
		}

		// Sends a notification with the following parameters:
		// Name: stats
		// Payload: map[string]interface{}{"cpu usage": c[0], "device name": h,}
		// Priority: 1
		cloud.SendNotification("stats", map[string]interface{}{
			"cpu usage":    c[0],
			"device name":  h,
		}, 1)

		log.Printf("Sent a notification with cpu usage %f and device name %s", c[0], h)

		time.Sleep(time.Second)
	}
}
