package params

import (
	"errors"
	"flag"
	"fmt"
	"strings"
)

type Port uint

type VideoFile struct {
	Name string
	Path string
}

type videofiles_parser struct {
	files []VideoFile
}

func (d *videofiles_parser) Set(str string) error {

	index := strings.Index(str, "=")

	name := str[:index]
	file := str[index+1:]

	videofile := VideoFile{name, file}
	d.files = append(d.files, videofile)
	return nil
}

func (d *videofiles_parser) String() string {
	return fmt.Sprintf("%v", d.files)
}

//---------------------------------------------------------------------------
//PARAMS
//---------------------------------------------------------------------------

type Params struct {
	Port  Port
	Files []VideoFile
}

func (p Params) GetVideoFile(name string) (VideoFile, bool) {

	for _, f := range p.Files {
		if f.Name == name {
			return f, true
		}
	}

	return VideoFile{}, false
}

//------------------------------------------------------------------------------
//------------------------------------------------------------------------------

func ParseParams() (Params, error) {

	var videofiles videofiles_parser
	flag.Var(&videofiles, "device", "proste devajsi")

	var port uint
	flag.UintVar(&port, "port", 8989, "proste port")

	flag.Parse()

	if len(videofiles.files) == 0 {
		return Params{}, errors.New("No video device entered. Use --device parameters")
	}

	return Params{Port(port), videofiles.files}, nil
}
