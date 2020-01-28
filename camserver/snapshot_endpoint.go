package camserver

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"webcam"

	"github.com/gorilla/mux"
)

type snapshot struct {
	Width  uint32 `json:"width"`
	Height uint32 `json:"height"`
	Data   string `json:"data"`
}

const (
	DEFAULT_WIDTH  uint32 = 640
	DEFAULT_HEIGHT uint32 = 480
)

func snapshotHandler(writer http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)
	name := vars["name"]

	format, ok := resolveOutputFormat(request)

	if !ok {
		logAndWriteResponse("Bad value of param 'format'", nil, writer)
		return
	}

	file, ok := parameters.GetVideoFile(name)

	if !ok {
		logAndWriteResponse(fmt.Sprintf("There is no device '%s'", name), nil, writer)
		return
	}

	device, err := webcam.OpenVideoDevice(file.Path)

	if err != nil {
		logAndWriteResponse(fmt.Sprintf("Cannot read device '%s'\n", name), err, writer)
		return
	}

	defer func() {
		if err := device.Close(); err != nil {
			log.Printf("Cannot close device %s: %v\n", file.Path, err)
		}
	}()

	framesize, err := resolveFrameSize(request, device)

	if err != nil {
		logAndWriteResponse("No frame size resolved", err, writer)
		return
	}

	snap, err := device.TakeSnapshot(&framesize)

	if err != nil {
		logAndWriteResponse("Cannot take snapshot", err, writer)
		return
	}

	contentType := resolveContentType(format)
	b := formatPayload(snap, format)

	writer.Header().Set("Content-Type", contentType)
	writer.Write(b)
}

func logAndWriteResponse(m string, err error, writer http.ResponseWriter) {
	var message string
	if err != nil {
		message = fmt.Sprintf("%v: %v\n", m, err)
	} else {
		message = m
	}

	log.Printf(message)
	writer.Write([]byte(message))
}

//-------------------------------------------------------------------------------
//RESOLVING FRAME SIZE
//-------------------------------------------------------------------------------

type eval func(frameSize webcam.DiscreteFrameSize) uint32

func findNearestFrameSizeByWidth(frameSizes []webcam.DiscreteFrameSize, width uint32) webcam.DiscreteFrameSize {
	return findFrameSize(frameSizes, func(frameSize webcam.DiscreteFrameSize) uint32 {
		return uint32(math.Abs(float64(frameSize.Width - width)))
	})
}

func findNearestFrameSizeByHeight(frameSizes []webcam.DiscreteFrameSize, height uint32) webcam.DiscreteFrameSize {
	return findFrameSize(frameSizes, func(frameSize webcam.DiscreteFrameSize) uint32 {
		return uint32(math.Abs(float64(frameSize.Height - height)))
	})
}

func findFrameSize(sizes []webcam.DiscreteFrameSize, evaluation eval) webcam.DiscreteFrameSize {
	var evalTotal uint32 = math.MaxUint32
	var minIndex int = 0

	for i, s := range sizes {
		eval := evaluation(s)
		if eval < evalTotal {
			minIndex = i
			evalTotal = eval
		}
	}

	return sizes[minIndex]
}

func resolveFrameSize(request *http.Request, device webcam.VideoDevice) (webcam.DiscreteFrameSize, error) {
	queries := request.URL.Query()

	widthStr, wok := queries["width"]
	heightStr, hok := queries["height"]

	result := webcam.DiscreteFrameSize{}

	if !wok && !hok {
		log.Println(fmt.Sprintf("No resolution setup. Setting default %dx%d", DEFAULT_WIDTH, DEFAULT_HEIGHT))
		result.Width = DEFAULT_WIDTH
		result.Height = DEFAULT_HEIGHT
		return result, nil
	}

	sizes, err := device.FrameSizes().AllDiscreteMJPEG()
	if err != nil {
		return result, err
	}

	if !wok {
		log.Println("Width missing, looking for appropriate one")
		height, err := strconv.Atoi(heightStr[0])

		if err != nil {
			return result, err
		}

		framesize := findNearestFrameSizeByHeight(sizes, uint32(height))
		return framesize, nil
	}

	if !hok || hok && wok {
		log.Println("Width missing, looking for appropriate one")
		width, err := strconv.Atoi(widthStr[0])

		if err != nil {
			return result, err
		}

		framesize := findNearestFrameSizeByWidth(sizes, uint32(width))
		return framesize, nil
	}

	panic("Unreachable")
}

//----------------------------------------------------------------------------
//RESOLVING OUTPUT FORMAT
//----------------------------------------------------------------------------

func resolveOutputFormat(request *http.Request) (string, bool) {

	queries := request.URL.Query()

	formats, ok := queries["format"]

	if !ok {
		return "json", true
	}

	var format string = formats[0]

	switch format {
	case "json":
		return "json", true

	case "raw":
		return "raw", true

	default:
		return "", false
	}
}

func formatPayload(snap webcam.Snapshot, format string) []byte {
	if format == "json" {
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

		return b
	}

	if format == "raw" {
		return snap.Data()
	}

	panic("No format")
}

func resolveContentType(format string) string {

	switch format {
	case "json":
		return "application/json"

	case "raw":
		return "image/jpeg"

	default:
		panic("no format")
	}
}
