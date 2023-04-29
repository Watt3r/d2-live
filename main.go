package main

import (
	"log"
	"net/http"

	statsd "github.com/DataDog/datadog-go/v5/statsd"
	"github.com/husobee/vestigo"
	"github.com/watt3r/d2-live/internal/handlers"
)

var Version string

func rootHandler(rw http.ResponseWriter, req *http.Request) {
	http.Redirect(rw, req, "https://github.com/Watt3r/d2-live", 301)
}

func main() {
	metricsClient, err := statsd.New("172.17.33.150:8125",
		statsd.WithTags([]string{"env:prod", "service:myservice"}),
	)
	if err != nil {
		log.Fatal(err)
	}

	c := handlers.Controller{
		Metrics: metricsClient,
		Version: Version,
	}

	router := vestigo.NewRouter()

	router.Get("/", rootHandler, c.StatsdMiddleware)

	router.Get("/info", c.GetInfoHandler, c.StatsdMiddleware)

	router.Get("/svg/:encodedD2", c.GetD2SVGHandler, c.StatsdMiddleware)

	log.Fatal(http.ListenAndServe(":8090", router))
}
