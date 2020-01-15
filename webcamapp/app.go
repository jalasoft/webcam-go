package main

import (
	"fmt"
	"log"
	"v4l2"
	"webcam"
)

func main() {
	takeSnapshot("/dev/video0")
	//printAllFrameSizes("/dev/video0")
	//printCapability("/dev/video0")
	//printFormatSupport("/dev/video0")
}

func takeSnapshot(file string) {
	device, err := webcam.OpenVideoDevice(file)

	if err != nil {
		log.Fatalf("%v\n", err)
	}

	defer device.Close()

	snapshot := device.Snapshot()

	err2 := snapshot.Take(webcam.DiscreteFrameSize{1280, 960})

	if err2 != nil {
		log.Fatalf("%v\n", err)
	}
}

func printAllFrameSizes(file string) {

	device, err := webcam.OpenVideoDevice(file)

	if err != nil {
		log.Fatalf("%v\n", err)
	}

	defer device.Close()

	sizes := device.FrameSizes()

	discretes, err := sizes.AllDiscrete(v4l2.V4L2_PIX_FMT_MJPEG)

	if err != nil {
		log.Fatalf("%v\n", err)
	}

	for _, d := range discretes {
		fmt.Printf("%v\n", d)
	}

	fmt.Println("--------------------------")

	supports, err := sizes.SupportsDiscrete(v4l2.V4L2_PIX_FMT_MJPEG, 1184, 656)

	if err != nil {
		log.Fatalf("%v\n", err)
	}

	fmt.Printf("1184x656 je podporovano: %t\n", supports)
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
