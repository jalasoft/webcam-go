package camserver

import (
	"camserver/params"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"webcam"
)

type camera_info struct {
	Name    string `json:"name"`
	File    string `json:"file"`
	Driver  string `json:"driver"`
	Card    string `json:"card"`
	Businfo string `json:"bus_info"`
	Version uint32 `json:"version"`
}

func allCamerasHandler(writer http.ResponseWriter, request *http.Request) {

	payload := make([]camera_info, 0)
	channel := make(chan camera_info)

	go processCameraInfo(channel)

	for info := range channel {
		payload = append(payload, info)
	}

	b, err := json.MarshalIndent(payload, "", "  ")

	if err != nil {
		writer.Write([]byte(fmt.Sprintf("%v", err)))
	}

	writer.Write(b)
}

func processCameraInfo(channel chan<- camera_info) {
	var group sync.WaitGroup

	deviceCount := len(parameters.Files)
	group.Add(deviceCount)

	for _, file := range parameters.Files {
		go readCameraInfo(file, &group, channel)
	}

	group.Wait()
	close(channel)
}

func readCameraInfo(file params.VideoFile, group *sync.WaitGroup, channel chan<- camera_info) {

	device, err := webcam.OpenVideoDevice(file.Path)

	info := camera_info{}
	info.Name = file.Name
	info.File = file.Path

	if err != nil {
		info.Driver = fmt.Sprintf("cannot load: %v", err)

		channel <- info
		group.Done()
		return
	}

	defer func() {
		if err := device.Close(); err != nil {
			log.Printf("Cannot close device %s: %v\n", file.Path, err)
		}
	}()

	cap := device.Capability()

	info.Driver = trim(cap.Driver())
	info.Card = trim(cap.Card())
	info.Businfo = trim(cap.BusInfo())
	info.Version = cap.Version()

	fmt.Printf("'%v'\n", info.Driver)

	channel <- info

	group.Done()
	return
}

func trim(value string) string {
	return strings.Trim(value, string('\u0000'))
}
