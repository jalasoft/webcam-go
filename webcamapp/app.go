package main

import (
	"camserver"
	"fmt"
	"log"
	"os"
	"time"
	"v4l2"
	"webcam"
)

func main() {

	camserver.StartServer()

	//writeBinaryFile()
	//printConstants()
	//streamVideo("/dev/video0")
	//printAllFrameSizes("/dev/video0")
	//printCapability("/dev/video0")
	//printFormatSupport("/dev/video0")
}

func streamVideo(file string) {
	device, err := webcam.OpenVideoDevice(file)

	if err != nil {
		log.Fatalf("%v\n", err)
	}

	defer func() {
		if err2 := device.Close(); err2 != nil {
			log.Fatalf("%v\n")
		}
	}()

	ticks := make(chan bool, 1)
	snaps := make(chan webcam.Snapshot)

	go device.Stream(&webcam.DiscreteFrameSize{1280, 960}, ticks, snaps)
	go tickDriving(40, ticks)

	index := uint(0)
	for s := range snaps {
		file, err := os.Create(fmt.Sprintf("/home/honzales/stream/%d.jpg", index))
		if err != nil {
			log.Fatalf("%v\n", err)
		}

		file.Write(s.Data())

		index++
	}

	fmt.Printf("Konec streamu")
}

func tickDriving(count int, ticks chan bool) {

	for i := 0; i < count; i++ {
		time.Sleep(300 * time.Millisecond)
		ticks <- true

		ok := <-ticks

		if !ok {
			break
		}
	}

	close(ticks)
}

func takeSnapshot(file string) {
	device, err := webcam.OpenVideoDevice(file)

	if err != nil {
		log.Fatalf("%v\n", err)
	}

	defer device.Close()

	ch := make(chan webcam.Snapshot)
	go device.TakeSnapshotChan(&webcam.DiscreteFrameSize{1280, 960}, ch)

	for s := range ch {
		fmt.Printf("Mam obrazek o velikosti %dB\n", s.Length())

		outfile, err := os.Create("/home/honzales/snapshot.jpg")
		defer outfile.Close()

		if err != nil {
			log.Fatalf("%v\n", err)
		}

		outfile.Write(s.Data())
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
