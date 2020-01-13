package webcam

import (
	"log"
	"os"
)

type device struct {
	file       *os.File
	capability v4l2Capability
	formats    supportedFormats
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

func (d *device) Close() {
	log.Printf("Closing video device.\n")
	d.file.Close()
}
