package camserver

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"webcam"

	"github.com/gorilla/mux"
)

type snapshot struct {
	Width  uint32 `json:"width"`
	Height uint32 `json:"height"`
	Data   string `json:"data"`
}

func snapshotHandler(writer http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)
	name := vars["name"]

	ok, file := parameters.GetVideoFile(name)

	if !ok {
		log.Printf("There is no device '%s'", name)
		writer.Write([]byte(fmt.Sprintf("No device '%s'\n", name)))
		return
	}

	device, err := webcam.OpenVideoDevice(file.Path)

	if err != nil {
		message := fmt.Sprintf("Cannot read device '%s'\n", name)
		log.Printf(message)
		writer.Write([]byte(message))
	}

	defer func() {
		if err := device.Close(); err != nil {
			log.Printf("Cannot close device %s: %v\n", file.Path, err)
		}
	}()

	framesize := webcam.DiscreteFrameSize{640, 480}

	snap, err := device.TakeSnapshot(&framesize)

	if err != nil {
		message := fmt.Sprintf("Cannot take snapshot: %v\n", err)
		log.Printf(message)
		writer.Write([]byte(message))
		return
	}

	payload := snapshot{}
	payload.Width = 640
	payload.Height = 480
	payload.Data = base64.StdEncoding.EncodeToString(snap.Data())

	b, err := json.MarshalIndent(payload, "", "  ")

	if err != nil {
		message := "Cannot marshall response to json"
		log.Printf(message)
		panic(message)
	}

	writer.Write(b)
}
