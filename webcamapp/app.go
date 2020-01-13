package main

import (
	"fmt"
	"log"
	"v4l2"
	"webcam"
)

func main() {
	//printCapability("/dev/video0")
	printFormatSupport("/dev/video0")
}

func printFormatSupport(path string) {
	device, err := webcam.OpenVideoDevice(path)

	if err != nil {
		log.Fatalf("%v\n", err)
	}

	defer device.Close()

	supports, err := device.Formats().Supports(v4l2.V4L2_BUF_TYPE_VIDEO_CAPTURE, v4l2.V4L2_PIX_FMT_MJPEG)

	if err != nil {
		log.Fatalf("%v\n", err)
	}

	fmt.Printf("Device %s supports format %s: %t\n", device.Name(), "V4L2_PIX_FMT_MJPEG", supports)
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
