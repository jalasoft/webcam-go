package webcam

import "v4l2"

func (c v4l2.V4l2Capability) bus() string {
	return string(c.BusInfo[:])
}
