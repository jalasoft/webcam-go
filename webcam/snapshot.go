package webcam

import (
	"fmt"
	"os"
)

type snapshot struct {
	file *os.File
}

func (s *snapshot) Take(frameSize DiscreteFrameSize) {

	fmt.Println("Cvak...")

	//ioctl.SetFrameSize(s.file.Fd())
}
