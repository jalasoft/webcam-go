package camserver

import (
	"camserver/params"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type cameraInfo struct {
	Driver  string `json:"driver"`
	Card    string `json:"card"`
	Businfo string `json:"bus_info"`
}

var parameters params.Params

func StartServer() {

	par, err := params.ParseParams()

	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	parameters = par

	log.Printf("starting server on port %d", parameters.Port)

	router := mux.NewRouter()

	router.HandleFunc("/camera/", allCamerasHandler).Methods("GET")
	router.HandleFunc("/camera/{name}", cameraHandler).Methods("GET")
	router.HandleFunc("/camera/{name}/snapshot", snapshotHandler).Methods("GET")

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", parameters.Port),
		Handler: router,
	}

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
