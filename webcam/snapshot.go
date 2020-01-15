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

	var err error

	log.Printf("Setting up frame size %dx%d", frameSize.Width, frameSize.Height)
	err = s.setFrameSize(&frameSize, v4l2.V4L2_PIX_FMT_MJPEG)
	if err != nil {
		return err
	}

	log.Printf("Frame size set up")
	log.Printf("Requesting buffer")
	err = s.requestBuffer(v4l2.V4L2_MEMORY_MMAP)
	if err != nil {
		return err
	}
	log.Printf("Buffer requested successfully")

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

	return ioctl.SetFrameSize(s.file.Fd(), &format)
}

func (s *snapshot) requestBuffer(memory uint32) error {

	var request v4l2.V4l2RequestBuffers
	request.Count = 1
	request.Type = v4l2.V4L2_BUF_TYPE_VIDEO_CAPTURE
	request.Memory = memory

	return ioctl.RequestBuffer(s.file.Fd(), &request)
}
