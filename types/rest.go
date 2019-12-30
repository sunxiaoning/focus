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

type PageRequest struct {
	PageIndex int `json:"pageIndex"`
	PageSize  int `json:"pageSize"`
}

type PageResponse struct {
	Total      int         `json:"total"`
	ResultList interface{} `json:"resultList"`
}

func NewPageResponse(total int, resultList interface{}) *PageResponse {
	return &PageResponse{
		Total:      total,
		ResultList: resultList,
	}
}

func NewRestRestResponse(rw http.ResponseWriter, data interface{}) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(rw)
	encoder.SetEscapeHTML(false)
	res := &restResponse{Success, "", data}
	panic(encoder.Encode(res))
}
