package controller

type ResResponse struct {
	ResponseCode    int         `json:"responseCode"`
	ResponseMessage string      `json:"responseMessage"`
	Data            interface{} `json:"data,omitempty"`
}
