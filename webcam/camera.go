package webcam

import (
	"fmt"
	"log"
	"os"
	"v4l2"
)

//-----------------------------------------------------
//SNAPSHOT
//-----------------------------------------------------

type snapshot struct {
	framesize *DiscreteFrameSize
	data      []byte
	length    uint32
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

type camera struct {
	file *os.File
}

func (s *camera) takeSnapshotChan(frameSize *DiscreteFrameSize, ch chan Snapshot) {

	stream := &streaming{file: s.file}

	if err := stream.Open(frameSize); err != nil {
		log.Fatalf("%v\n", err)
	}

	defer stream.Close()

	snap, err := stream.Snapshot()

	if err != nil {
		panic(fmt.Sprintf("%v\n", err))
	}

	ch <- snap
	close(ch)
}

func (s *camera) takeSnapshot(frameSize *DiscreteFrameSize) (Snapshot, error) {

	var sn *snapshot

	err := s.takeSnapshotAsync(frameSize, func(snap Snapshot) {
		var dataCopy []byte = make([]byte, snap.Length())
		copy(dataCopy, snap.Data())
		sn = &snapshot{snap.FrameSize(), dataCopy, snap.Length()}
	})

	if err != nil {
		return nil, err
	}

	return sn, nil
}

func (s *camera) takeSnapshotAsync(frameSize *DiscreteFrameSize, handler SnapshotHandler) error {

	log.Printf("Setting up frame size %dx%d", frameSize.Width, frameSize.Height)
	if err := setFrameSize(s.file.Fd(), frameSize, v4l2.V4L2_PIX_FMT_MJPEG); err != nil {
		return err
	}

	log.Printf("Frame size set up")
	log.Printf("Requesting buffer")
	if err := requestMmapBuffer(s.file.Fd()); err != nil {
		return err
	}
	log.Printf("Buffer requested successfully")
	log.Printf("Querying mmap buffer")
	offset, length, err := queryMmapBuffer(s.file.Fd())

	if err != nil {
		return err
	}

	log.Printf("Mmap buffer obtained. Offset=%v, length=%v\n", offset, length)
	log.Printf("Retrieving mapped memory block, offset=%d, length=%d", offset, length)

	data, err := mapBuffer(s.file.Fd(), offset, length)
	if err != nil {
		return err
	}

	log.Println("Activating streaming")
	if err := activateStreaming(s.file.Fd()); err != nil {
		return err
	}

	log.Println("Queueing buffer")
	var buffer v4l2.V4l2Buffer
	buffer.Index = uint32(0)
	buffer.Type = v4l2.V4L2_BUF_TYPE_VIDEO_CAPTURE
	buffer.Memory = v4l2.V4L2_MEMORY_MMAP

	if err := queueBuffer(s.file.Fd(), &buffer); err != nil {
		return err
	}
	log.Println(fmt.Sprintf("Buffer filled with %d bytes", buffer.Length))

	log.Println("Dequeuing the buffer")
	if err := dequeueBuffer(s.file.Fd(), &buffer); err != nil {
		return err
	}

	snapshot := &snapshot{frameSize, data, length}
	handler(snapshot)

	log.Printf("Releasing mapped memory block")
	if err := munmapBuffer(data); err != nil {
		return err
	}

	log.Println("Deactivating streaming")
	if err := deactivateStreaming(s.file.Fd()); err != nil {
		return err
	}

	return nil
}

//--------------------------------------------------------------------------------------------------
//STREAMING
//--------------------------------------------------------------------------------------------------

type streaming struct {
	file      *os.File
	frameSize *DiscreteFrameSize
	length    uint32
	data      []byte
}

func (s *streaming) Open(frameSize *DiscreteFrameSize) error {
	log.Printf("Setting up frame size %dx%d", frameSize.Width, frameSize.Height)
	if err := setFrameSize(s.file.Fd(), frameSize, v4l2.V4L2_PIX_FMT_MJPEG); err != nil {
		return err
	}

	s.frameSize = frameSize

	log.Printf("Frame size set up")
	log.Printf("Requesting buffer")
	if err := requestMmapBuffer(s.file.Fd()); err != nil {
		return err
	}
	log.Printf("Buffer requested successfully")
	log.Printf("Querying mmap buffer")
	offset, length, err := queryMmapBuffer(s.file.Fd())

	if err != nil {
		return err
	}

	s.length = length

	log.Printf("Mmap buffer obtained. Offset=%v, length=%v\n", offset, length)
	log.Printf("Retrieving mapped memory block, offset=%d, length=%d", offset, length)

	data, err := mapBuffer(s.file.Fd(), offset, length)
	if err != nil {
		return err
	}

	s.data = data

	log.Println("Activating streaming")
	if err := activateStreaming(s.file.Fd()); err != nil {
		return err
	}

	return nil
}

/*
func (s *streaming) Stream(tickChannel <-chan bool, frameChannel chan<- Snapshot) {

	for range tickChannel {
		snap, err := s.makeSnapshot()

		if err != nil {
			close(frameChannel)
			panic(fmt.Sprintf("%v\n", err))
		}

		frameChannel <- snap
	}

	close(frameChannel)
}*/

func (s *streaming) Snapshot() (Snapshot, error) {

	log.Println("Queueing buffer")
	var buffer v4l2.V4l2Buffer
	buffer.Index = uint32(0)
	buffer.Type = v4l2.V4L2_BUF_TYPE_VIDEO_CAPTURE
	buffer.Memory = v4l2.V4L2_MEMORY_MMAP

	if err := queueBuffer(s.file.Fd(), &buffer); err != nil {
		return nil, err
	}
	log.Println(fmt.Sprintf("Buffer filled with %d bytes", buffer.Length))

	log.Println("Dequeuing the buffer")
	if err := dequeueBuffer(s.file.Fd(), &buffer); err != nil {
		return nil, err
	}

	snapshot := &snapshot{s.frameSize, s.data, s.length}
	return snapshot, nil
}

func (s *streaming) Close() error {
	log.Printf("Releasing mapped memory block")
	if err := munmapBuffer(s.data); err != nil {
		return err
	}

	log.Println("Deactivating streaming")
	if err := deactivateStreaming(s.file.Fd()); err != nil {
		return err
	}

	return nil
}
