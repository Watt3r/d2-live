package handlers

import (
	"encoding/json"
	"net/http"
)

type InfoResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

func (c *Controller) GetInfoHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(InfoResponse{Status: "Success", Version: c.Version})
}
