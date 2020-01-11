package webcam

import (
	"os"
	"v4l2"
)

type device struct {
	file       os.File
	capability v4l2.V4l2Capability
}

func (d device) Close() {
	d.file.Close()
}
