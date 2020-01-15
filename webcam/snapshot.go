package webcam

import (
	"log"
	"os"
	"v4l2"
	"v4l2/ioctl"
)

type snapshot struct {
	file *os.File
}

func (s *snapshot) Take(frameSize DiscreteFrameSize) error {

	s.setFrameSize(&frameSize, v4l2.V4L2_PIX_FMT_MJPEG)

	return nil
}

func (s *snapshot) setFrameSize(frameSize *DiscreteFrameSize, pixelFormat uint32) error {
	var format v4l2.V4l2Format

	var pixFormat v4l2.V4l2PixFormat
	pixFormat.Width = frameSize.Width
	pixFormat.Height = frameSize.Height
	pixFormat.Pixelformat = pixelFormat
	pixFormat.Field = v4l2.V4L2_FIELD_NONE

	format.SetPixFormat(&pixFormat)

	log.Printf("Setting up frame size %dx%d", frameSize.Width, frameSize.Height)

	err := ioctl.SetFrameSize(s.file.Fd(), &format)

	if err != nil {
		return err
	}

	log.Printf("Frame size set up")
	return nil
}
