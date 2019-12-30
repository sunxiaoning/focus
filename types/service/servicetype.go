package servicetype

import (
	"focus/types"
	"time"
)

type ServiceEntity struct {
	ID            int
	ServiceType   int
	ServiceName   string
	ChineseName   string
	ServiceDesc   string
	ServiceStatus string
	PublishTime   time.Time
}

type QueryLatestReq struct {
	types.PageRequest
	ServiceType int    `json:"serviceType"`
	ChineseName string `json:"chineseName"`
}

func NewQueryLatestReq() *QueryLatestReq {
	return &QueryLatestReq{
		PageRequest: types.PageRequest{
			PageIndex: 1,
			PageSize:  10,
		},
		ServiceType: 0,
		ChineseName: "",
	}
}

type QueryLatestRes struct {
	ServiceId   int    `json:"serviceId"`
	ServiceType int    `json:"serviceType"`
	ChineseName string `json:"chineseName"`
	ServiceDesc string `json:"serviceDesc"`
	PublishTime string `json:"publishTime"`
}

type GetByIdRes QueryLatestRes

type PriceEntity struct {
	ID                int
	PriceName         string
	ServiceId         int
	ConcurrencyNumber int
	Price             string
	PriceType         string
	ServiceAmount     int64
}

type QueryPriceRes struct {
	ID            int    `json:"id"`
	PriceName     string `json:"priceName"`
	ServiceId     int    `json:"serviceId"`
	Price         string `json:"price"`
	ServiceAmount int64  `json:"serviceAmount"`
}

type CalculatePriceReq struct {
	PriceId int   `json:"priceId"`
	Amount  int64 `json:"amount"`
}

type CalculatePriceRes struct {
	Price string `json:"price"`
}

func NewCalculatePriceReq() *CalculatePriceReq {
	return &CalculatePriceReq{0, 1}
}

type OrderEntity struct {
	ID             int
	OrderNo        string
	MemberId       int
	ServicePriceId int
	PurchaseAmount int64
	StartTime      time.Time
	FinishedTime   time.Time
	OrderStatus    string
	OutOrderNo     string
	OrderAmount    string
	PayAmount      string
	CouponNo       string
	PayChannel     string
}
type CreateOrderRequest struct {
	OrderNo        string `json:"orderNo"`
	MemberId       int    `json:"memberId"`
	ServicePriceId int    `json:"servicePriceId"`
	PurchaseAmount int64  `json:"purchaseAmount"`
	CouponNo       string `json:"couponNo"`
	PayChannel     string `json:"payChannel"`
}

func NewCreateOrderReq() *CreateOrderRequest {
	return &CreateOrderRequest{
		OrderNo:        "",
		MemberId:       0,
		ServicePriceId: 0,
		PurchaseAmount: 0,
		CouponNo:       "",
	}
}

type CreateOrderRes struct {
	CashierUrl    string      `json:"cashierUrl"`
	CashierParams *CashierReq `json:"cashierParams"`
}

type CashierReq struct {
	PayOrderNo     string `json:"payOrderNo"`
	ServiceOrderNo string `json:"serviceOrderNo"`
	OrderAmount    string `json:"orderAmount"`
	PayAmount      string `json:"payAmount"`
	PayChannel     string `json:"payChannel"`
	PayReason      string `json:"payReason"`
}
