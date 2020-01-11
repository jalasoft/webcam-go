package webcam

import (
	"log"
	"os"
	"v4l2"
)

type device struct {
	file       *os.File
	capability v4l2.V4l2Capability
}

func (d device) Close() {
	log.Printf("Closing video device.\n")
	d.file.Close()
}
