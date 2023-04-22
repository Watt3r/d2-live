package main

import (
	"bytes"
	"compress/flate"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"sort"

	"github.com/husobee/vestigo"
	"oss.terrastruct.com/d2/d2graph"
	"oss.terrastruct.com/d2/d2layouts/d2dagrelayout"
	"oss.terrastruct.com/d2/d2lib"
	"oss.terrastruct.com/d2/d2renderers/d2svg"
	"oss.terrastruct.com/d2/lib/textmeasure"
)

var Version string

var compressionDict = "->" +
	"<-" +
	"--" +
	"<->"

// Decode decodes a compressed base64 D2 string.
func Decode(encoded string) (_ string, err error) {
	b64Decoded, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}

	zr := flate.NewReaderDict(bytes.NewReader(b64Decoded), []byte(compressionDict))
	var b bytes.Buffer
	if _, err := io.Copy(&b, zr); err != nil {
		return "", err
	}
	if err := zr.Close(); err != nil {
		return "", nil
	}
	return b.String(), nil
}

func init() {
	var common []string
	for k := range d2graph.StyleKeywords {
		common = append(common, k)
	}
	for k := range d2graph.ReservedKeywords {
		common = append(common, k)
	}
	for k := range d2graph.ReservedKeywordHolders {
		common = append(common, k)
	}
	sort.Strings(common)
	for _, k := range common {
		compressionDict += k
	}
}

type InfoResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

func infoHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(InfoResponse{Status: "Success", Version: Version})
}

func GetD2SVGHandler(rw http.ResponseWriter, req *http.Request) {
	ctx := context.Background()

	svg, err := handleGetD2SVG(ctx, req)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	rw.Header().Set("Content-Type", "image/svg+xml")
	rw.Write(svg)
}

func handleGetD2SVG(ctx context.Context, req *http.Request) ([]byte, error) {
	urlencoded := vestigo.Param(req, "encodedD2")
	decoded, err := Decode(urlencoded)
	if err != nil {
		return nil, errors.New("Invalid Base64 data.")
	}

	ruler, _ := textmeasure.NewRuler()

	diagram, _, _ := d2lib.Compile(ctx, decoded, &d2lib.CompileOptions{
		Layout: d2dagrelayout.DefaultLayout,
		Ruler:  ruler,
	})

	// Render to SVG
	out, err := d2svg.Render(diagram, &d2svg.RenderOpts{
		Pad: d2svg.DEFAULT_PADDING,
	})
	if err != nil {
		return nil, errors.New("Invalid D2 data.")
	}
	return out, nil
}

func main() {
	router := vestigo.NewRouter()

	router.Get("/info", infoHandler)

	router.Get("/svg/:encodedD2", GetD2SVGHandler)

	log.Fatal(http.ListenAndServe(":8090", router))
}
