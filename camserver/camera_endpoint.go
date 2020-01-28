package camserver

import (
	"camserver/params"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"v4l2"
	"webcam"

	"github.com/gorilla/mux"
)

type supported_resolution struct {
	Width  uint32 `json:"width"`
	Height uint32 `json:"height"`
}

type camera_full_info struct {
	Info        camera_info            `json:"info"`
	Resolutions []supported_resolution `json:"resolutions"`
}

func cameraHandler(writer http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)
	name := vars["name"]

	file, ok := parameters.GetVideoFile(name)

	if !ok {
		log.Printf("There is no device '%s'", name)
		writer.Write([]byte(fmt.Sprintf("No device '%s'\n", name)))
		return
	}

	err, info := readCameraFullInfo(file)

	if err != nil {
		log.Printf("Cannot get full info for device %v: %v\n", file, err)
		writer.Write([]byte(fmt.Sprintf("Cannot get full info from device %v", file.Name)))
		return
	}
	b, err := json.MarshalIndent(info, "", "  ")

	if err != nil {
		log.Printf("Cannot marshal response: %v", err)
		writer.Write([]byte("Cannot marshall response"))
		return
	}

	writer.Write(b)
}

func readCameraFullInfo(file params.VideoFile) (error, camera_full_info) {

	fullInfo := camera_full_info{}
	fullInfo.Info = camera_info{}

	device, err := webcam.OpenVideoDevice(file.Path)

	if err != nil {
		fullInfo.Info.Driver = fmt.Sprintf("cannot load: %v", err)
		return err, fullInfo
	}

	defer func() {
		if err := device.Close(); err != nil {
			log.Printf("Cannot close device %s: %v\n", file.Path, err)
		}
	}()

	cap := device.Capability()
	fullInfo.Info.Driver = trim(cap.Driver())
	fullInfo.Info.Card = trim(cap.Card())
	fullInfo.Info.Businfo = trim(cap.BusInfo())
	fullInfo.Info.Version = cap.Version()

	frames, err := device.FrameSizes().AllDiscrete(v4l2.V4L2_PIX_FMT_MJPEG)

	if err != nil {
		log.Printf("Cannot load frame sizes: %v", err)
		return err, fullInfo
	}

	resolutions := make([]supported_resolution, 0, 10)

	for _, frame := range frames {
		resolutions = append(resolutions, supported_resolution{frame.Width, frame.Height})
	}

	fullInfo.Resolutions = resolutions

	return nil, fullInfo
}
