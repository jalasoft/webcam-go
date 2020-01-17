package webcam

import (
	"log"
	"os"
	"syscall"
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
	err = s.requestMmapBuffer()
	if err != nil {
		return err
	}
	log.Printf("Buffer requested successfully")
	log.Printf("Querying mmap buffer")
	offset, length, err := s.queryMmapBuffer()

	if err != nil {
		return err
	}

	log.Printf("Mmap buffer obtained. Offset=%v, length=%v\n", offset, length)
	log.Printf("Retrieving mapped memory block, offset=%d, length=%d", offset, length)

	data, err := s.mapBuffer(offset, length)

	if err != nil {
		return err
	}

	log.Println("Activating streaming")
	s.activateStreaming()

	log.Println("Deactivating streaming")
	s.deactivateStreaming()

	log.Printf("Releasing mapped memory block")
	err2 := s.munmapBuffer(data)

	if err2 != nil {
		return err2
	}

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

func (s *snapshot) requestMmapBuffer() error {

	var request v4l2.V4l2RequestBuffers
	request.Count = 1
	request.Type = v4l2.V4L2_BUF_TYPE_VIDEO_CAPTURE
	request.Memory = v4l2.V4L2_MEMORY_MMAP

	return ioctl.RequestBuffer(s.file.Fd(), &request)
}

func (s *snapshot) queryMmapBuffer() (uint32, uint32, error) {

	var buffer v4l2.V4l2Buffer
	buffer.Index = uint32(0)
	buffer.Type = v4l2.V4L2_BUF_TYPE_VIDEO_CAPTURE
	buffer.Memory = v4l2.V4L2_MEMORY_MMAP

	ioctl.QueryBuffer(s.file.Fd(), &buffer)

	return buffer.Offset(), buffer.Length, nil
}

func (s *snapshot) mapBuffer(offset uint32, length uint32) ([]byte, error) {
	return syscall.Mmap(int(s.file.Fd()), int64(offset), int(length), syscall.PROT_READ, syscall.MAP_SHARED)
}

func (s *snapshot) munmapBuffer(data []byte) error {
	return syscall.Munmap(data)
}

func (s *snapshot) activateStreaming() {
	ioctl.ActivateStreaming(s.file.Fd(), v4l2.V4L2_BUF_TYPE_VIDEO_CAPTURE)
}

func (s *snapshot) deactivateStreaming() {
	ioctl.DeactivateStreaming(s.file.Fd(), v4l2.V4L2_BUF_TYPE_VIDEO_CAPTURE)
}
