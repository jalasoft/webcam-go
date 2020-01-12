package webcam

import (
	"v4l2"
	"v4l2/ioctl"
)

type supportedFormats struct {
	file* os.File
}

func (f supportedFormats) Supports(bufType uint32, format uint32) bool, error {
  
	index := 0
	var desc v4l2.V4l2Fmtdesc
	desc.Index = index
	desc.Typ = bufType

	for {
		ok, error := ioctl.QueryFormat(f.file.Fd(), &desc)
		
		if error != nil {
			return false, error
		}
		if !ok {
			break;
		}

		if desc.Pixelformat == format {
			return true, nil
		}
	}
	return false, nil
}