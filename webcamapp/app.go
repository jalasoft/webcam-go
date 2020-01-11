package main

import (
	"fmt"
	"log"
	"os"
	"v4l2/ioctl"
)

func main() {
	fmt.Println("ahoooj")

	file, err := os.Open("/dev/video0")

	if err != nil {
		log.Fatalf("%v\n", err)
	}

	defer file.Close()

	cap, err := ioctl.ReadCapability(file.Fd())

	if err != nil {
		log.Fatalf("%v\n", err)
	}

	fmt.Printf("%v\n", cap)
}
