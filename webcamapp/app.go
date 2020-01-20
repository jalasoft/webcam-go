package main

import (
	"fmt"
	"log"
	"os"
	"v4l2"
	"v4l2/ioctl"
	"webcam"
)

func main() {
	//writeBinaryFile()
	//printConstants()
	takeSnapshot("/dev/video0")
	//printAllFrameSizes("/dev/video0")
	//printCapability("/dev/video0")
	//printFormatSupport("/dev/video0")
}

func writeBinaryFile() {
	file, err := os.Create("/home/honzales/binarka")

	if err != nil {
		log.Fatalf("%v\n", err)
	}

	defer file.Close()

	var bindata []byte = []byte{0xaa, 0x12, 0x56}

	file.Write(bindata)

}

func printConstants() {
	fmt.Printf("VIDIOC_QUERYCAP: %v\n", ioctl.VIDIOC_QUERYCAP)
	fmt.Printf("VIDIOC_ENUM_FMT: %v\n", ioctl.VIDIOC_ENUM_FMT)
	fmt.Printf("VIDIOC_ENUM_FRAMESIZES: %v\n", ioctl.VIDIOC_ENUM_FRAMESIZES)
	fmt.Printf("VIDIOC_S_FMT: %v\n", ioctl.VIDIOC_S_FMT)
	fmt.Printf("VIDIOC_REQBUFS: %v\n", ioctl.VIDIOC_REQBUFS)
	fmt.Printf("VIDIOC_QUERYBUF: %v\n", ioctl.VIDIOC_QUERYBUF)
	fmt.Printf("VIDIOC_STREAMON: %v\n", ioctl.VIDIOC_STREAMON)
	fmt.Printf("VIDIOC_STREAMOFF: %v\n", ioctl.VIDIOC_STREAMOFF)
	fmt.Printf("VIDIOC_DQBUF: %v\n", ioctl.VIDIOC_DQBUF)
	fmt.Printf("VIDIOC_QBUF: %v\n", ioctl.VIDIOC_QBUF)

	fmt.Printf("V4L2_PIX_FMT_MJPEG: %v\n", v4l2.V4L2_PIX_FMT_MJPEG)
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
		fmt.Println("Mam snapshot....")

		file, err := os.Create("/home/honzales/cam1.jpg")
		if err != nil {
			log.Fatalf("%v\n", err)
		}

		defer file.Close()

		file.Write(s.Data())
	}

	/*
		snapshot, err := device.TakeSnapshot(&webcam.DiscreteFrameSize{1280, 960})
		if err != nil {
			log.Fatalf("%v\n", err)
		}

		fmt.Printf("Mam obrazek o velikosti %dB\n", snapshot.Length())

		outfile, err := os.Create("/home/honzales/snapshot.jpg")
		if err != nil {
			log.Fatalf("%v\n", err)
		}

		defer outfile.Close()

		//fmt.Printf("Zapisuju data do souboru\n")

		//fmt.Printf("%x", snapshot.Data())

		binary.Write(outfile, binary.LittleEndian, snapshot.Data())
		//fmt.Printf("Hotovo\n")*/

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
