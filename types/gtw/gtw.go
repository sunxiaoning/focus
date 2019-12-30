package gtwtype

// 网关请求参数
type GtWReq struct {
	Timestamp  string `json:"timestamp"`
	Sign       string `json:"sign"`
	MemberId   string `json:"memberId"`
	ServUrl    string `json:"servUrl"`
	BizContent string `json:"bizContent"`
}
