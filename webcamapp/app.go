package main

import (
	"log"
	"webcam"
)

func main() {
	device, err := webcam.OpenVideoDevice("/dev/video0")

	if err != nil {
		log.Fatalf("%v\n", err)
	}

	defer device.Close()

}
