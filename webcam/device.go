package webcam

import (
	"log"
	"os"
)

type device struct {
	file       *os.File
	capability v4l2Capability
	formats    supportedFormats
	framesizes *framesizes
	stillcamera   *stillcamera
}

func (d *device) Name() string {
	return d.file.Name()
}

func (d *device) Capability() Capability {
	return d.capability
}

func (d *device) Formats() SupportedFormats {
	return d.formats
}

func (d *device) FrameSizes() FrameSizes {
	return d.framesizes
}

func (d *device) TakeSnapshot(frameSize *DiscreteFrameSize) (Snapshot, error) {
	return d.stillcamera.TakeSnapshot(frameSize)
}

func (d *device) Close() {
	log.Printf("Closing video device.\n")
	d.file.Close()
}
