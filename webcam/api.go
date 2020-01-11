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

	if (cap.Capabilities & v4l2.V4L2_CAP_VIDEO_CAPTURE) == 0 {
		return nil, errors.New(fmt.Sprintf("Device %s is not a video capturing device.", path))
	}

	log.Println("Device is a video device")
	return device{file, cap}, nil
}

type VideoDevice interface {
	Close()
}
