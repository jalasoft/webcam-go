package webcam

import (
	"fmt"
	"log"
	"os"
	"syscall"
	"v4l2"
	"v4l2/ioctl"
)

//-----------------------------------------------------
//SNAPSHOT
//-----------------------------------------------------

type snapshot struct {
	framesize *DiscreteFrameSize
	data []byte
	length uint32
}

func (s *snapshot) FrameSize() *DiscreteFrameSize {
	return s.framesize
}

func (s *snapshot) Data() []byte {
	return s.data
}

func (s *snapshot) Length() uint32 {
	return s.length
}

//-----------------------------------------------------
//STILL CAMERA
//-----------------------------------------------------

type stillcamera struct {
	file *os.File
}

func (s *stillcamera) TakeSnapshot(frameSize *DiscreteFrameSize) (Snapshot, error) {

	log.Printf("Setting up frame size %dx%d", frameSize.Width, frameSize.Height)
	if err := s.setFrameSize(frameSize, v4l2.V4L2_PIX_FMT_MJPEG); err != nil {
		return nil, err
	}

	log.Printf("Frame size set up")
	log.Printf("Requesting buffer")
	if err := s.requestMmapBuffer(); err != nil {
		return nil, err
	}
	log.Printf("Buffer requested successfully")
	log.Printf("Querying mmap buffer")
	offset, length, err := s.queryMmapBuffer()

	if err != nil {
		return nil, err
	}

	log.Printf("Mmap buffer obtained. Offset=%v, length=%v\n", offset, length)
	log.Printf("Retrieving mapped memory block, offset=%d, length=%d", offset, length)

	data, err := s.mapBuffer(offset, length)
	if err != nil {
		return nil, err
	}
	
	var dataCopy []byte = make([]byte, length)

	log.Println("Activating streaming")
	if err := s.activateStreaming(); err != nil {
		return nil, err
	}

	log.Println("Queueing buffer")
	var buffer v4l2.V4l2Buffer
	buffer.Index = uint32(0)
	buffer.Type = v4l2.V4L2_BUF_TYPE_VIDEO_CAPTURE
	buffer.Memory = v4l2.V4L2_MEMORY_MMAP

	if err := s.queueBuffer(&buffer); err != nil {
		return nil, err
	}
	log.Println(fmt.Sprintf("Buffer filled with %d bytes", buffer.Length))

	log.Println("Dequeuing the buffer")
	if err := s.dequeueBuffer(&buffer); err != nil {
		return nil, err
	}

	//----
	file, _ := os.Create("/home/honzales/dalsi.jpg")
	file.Write(data)
	file.Close()
	//----

	copy(dataCopy, data)
	
	log.Printf("Releasing mapped memory block")
	if err := s.munmapBuffer(data); err != nil {
		return nil, err
	}

	log.Println("Deactivating streaming")
	if err := s.deactivateStreaming(); err != nil {
		return nil, err
	}



	return &snapshot{frameSize, dataCopy, length}, nil
}

func (s *stillcamera) setFrameSize(frameSize *DiscreteFrameSize, pixelFormat uint32) error {
	var format v4l2.V4l2Format

	var pixFormat v4l2.V4l2PixFormat
	pixFormat.Width = frameSize.Width
	pixFormat.Height = frameSize.Height
	pixFormat.Pixelformat = pixelFormat
	pixFormat.Field = v4l2.V4L2_FIELD_NONE

	format.SetPixFormat(&pixFormat)

	return ioctl.SetFrameSize(s.file.Fd(), &format)
}

func (s *stillcamera) requestMmapBuffer() error {

	var request v4l2.V4l2RequestBuffers
	request.Count = 1
	request.Type = v4l2.V4L2_BUF_TYPE_VIDEO_CAPTURE
	request.Memory = v4l2.V4L2_MEMORY_MMAP

	return ioctl.RequestBuffer(s.file.Fd(), &request)
}

func (s *stillcamera) queryMmapBuffer() (uint32, uint32, error) {

	buffer := &v4l2.V4l2Buffer{}
	buffer.Index = uint32(0)
	buffer.Type = v4l2.V4L2_BUF_TYPE_VIDEO_CAPTURE
	buffer.Memory = v4l2.V4L2_MEMORY_MMAP

	ioctl.QueryBuffer(s.file.Fd(), buffer)

	//fmt.Printf("%v\n", buffer)

	return buffer.Offset(), buffer.Length, nil
}

func (s *stillcamera) mapBuffer(offset uint32, length uint32) ([]byte, error) {
	return syscall.Mmap(int(s.file.Fd()), int64(offset), int(length), syscall.PROT_READ | syscall.PROT_WRITE, syscall.MAP_SHARED)
}

func (s *stillcamera) munmapBuffer(data []byte) error {
	return syscall.Munmap(data)
}

func (s *stillcamera) activateStreaming() error {
	return ioctl.ActivateStreaming(s.file.Fd(), v4l2.V4L2_BUF_TYPE_VIDEO_CAPTURE)
}

func (s *stillcamera) deactivateStreaming() error {
	return ioctl.DeactivateStreaming(s.file.Fd(), v4l2.V4L2_BUF_TYPE_VIDEO_CAPTURE)
}

func (s *stillcamera) queueBuffer(buffer* v4l2.V4l2Buffer) error {
	
	if err := ioctl.QueueBuffer(s.file.Fd(), buffer); err != nil {
		return err
	}

	return nil
}

func (s *stillcamera) dequeueBuffer(buffer *v4l2.V4l2Buffer) error {
	
	err := ioctl.DequeueBuffer(s.file.Fd(), buffer)

	if err != nil {
		return err
	}

	return nil
}
