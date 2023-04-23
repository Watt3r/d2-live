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
	"strings"
	"time"

	statsd "github.com/DataDog/datadog-go/v5/statsd"
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

func (c *Controller) rootHandler(rw http.ResponseWriter, req *http.Request) {
	http.Redirect(rw, req, "https://github.com/Watt3r/d2-live", 301)
}

type InfoResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

func (c *Controller) infoHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(InfoResponse{Status: "Success", Version: Version})
}

func (c *Controller) GetD2SVGHandler(rw http.ResponseWriter, req *http.Request) {
	ctx := context.Background()

	svg, err := c.handleGetD2SVG(ctx, req)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	rw.Header().Set("Content-Type", "image/svg+xml")
	rw.Write(svg)
}

func (c *Controller) handleGetD2SVG(ctx context.Context, req *http.Request) ([]byte, error) {
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

type Controller struct {
	Metrics *statsd.Client
}

func (c *Controller) statsdMiddleware(f http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		start := time.Now()
		f(rw, req)
		path := strings.Split(req.URL.Path, "/")[1]
		c.Metrics.Histogram("d2-live."+path, time.Now().Sub(start).Seconds(), []string{"version:" + Version}, 1)
	}
}

func main() {
	metricsClient, err := statsd.New("172.17.33.150:8125",
		statsd.WithTags([]string{"env:prod", "service:myservice"}),
	)
	if err != nil {
		log.Fatal(err)
	}

	c := Controller{
		Metrics: metricsClient,
	}

	router := vestigo.NewRouter()

	router.Get("/", c.rootHandler)

	router.Get("/info", c.infoHandler, c.statsdMiddleware)

	router.Get("/svg/:encodedD2", c.GetD2SVGHandler, c.statsdMiddleware)

	log.Fatal(http.ListenAndServe(":8090", router))
}
