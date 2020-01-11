package webcam

import (
	"fmt"
	"v4l2/ioctl"
)

func OpenVideoDevice(path string) VideoDevice, error {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	cap, err := ioctl.ReadCapability(file.Fd())

	if (cap.Capabilities & v4l2.V4L2_CAP_VIDEO_CAPTURE) == 0 {
		return nil, errors.New(fmt.Sprintf("Device %s is not a video capturing device.", path))
	}

}

type VideoDevice interface {

} 
