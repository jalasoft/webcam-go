package webcam

import (
	"log"
	"os"
)

type device struct {
	file       *os.File
	capability v4l2Capability
}

func (d device) Name() string {
	return d.file.Name()
}

func (d device) Capability() Capability {
	return d.capability
}

func (d device) Close() {
	log.Printf("Closing video device.\n")
	d.file.Close()
}
