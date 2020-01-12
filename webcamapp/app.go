package main

import (
	"fmt"
	"log"
	"webcam"
)

func main() {

	printCapability("/dev/video0")
}

func printCapability(file string) {
	device, err := webcam.OpenVideoDevice(file)

	if err != nil {
		log.Fatalf("%v\n", err)
	}

	defer device.Close()

	fmt.Printf("Video Device %s\n", device.Name())
	fmt.Printf("%v\n", device.Capability())
}
