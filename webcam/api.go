package webcam

import (
	"errors"
	"fmt"
	"log"
	"os"
	"v4l2"
	"v4l2/ioctl"
)

func OpenVideoDevice(path string) (VideoDevice, error) {
	file, err := os.Open(path)

	log.Printf("Opening device %s\n", path)

	if err != nil {
		return nil, err
	}

	log.Println("Reading capability")
	cap, err := ioctl.ReadCapability(file.Fd())

	if err != nil {
		return nil, err
	}

	var dev device = device{file, v4l2Capability{cap}}

	if !dev.Capability().HasCapability(v4l2.V4L2_CAP_VIDEO_CAPTURE) {
		return nil, errors.New(fmt.Sprintf("Device %s is not a video capturing device.", dev.Name()))
	}

	if !dev.Capability().HasCapability(v4l2.V4L2_CAP_STREAMING) {
		return nil, errors.New(fmt.Sprintf("Device %s is not able to stream frames.", dev.Name()))
	}

	log.Println("Device is a video device")
	return dev, nil
}

type VideoDevice interface {
	Name() string
	Capability() Capability
	Close()
}

type Capability interface {
	Driver() string
	Card() string
	BusInfo() string
	Version() uint32
	HasCapability(cap uint32) bool
}
