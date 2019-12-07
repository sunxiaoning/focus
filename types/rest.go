package types

import (
	"encoding/json"
	"net/http"
)

type restResponse struct {
	ResponseCode    int         `json:"responseCode"`
	ResponseMessage string      `json:"responseMessage"`
	Data            interface{} `json:"data,omitempty"`
}

func NewRestRestResponse(rw http.ResponseWriter, data interface{}) error {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(rw)
	encoder.SetEscapeHTML(false)
	res := &restResponse{Success, "", data}
	return encoder.Encode(res)
}
