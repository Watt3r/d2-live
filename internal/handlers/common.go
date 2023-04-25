package handlers

import (
	"net/http"
	"strings"
	"time"

	statsd "github.com/DataDog/datadog-go/v5/statsd"
)

type Controller struct {
	Metrics *statsd.Client
}

func (c *Controller) StatsdMiddleware(f http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		start := time.Now()
		f(rw, req)
		path := strings.Split(req.URL.Path, "/")[1]
		c.Metrics.Histogram("d2-live."+path, time.Now().Sub(start).Seconds(), []string{}, 1)
	}
}
