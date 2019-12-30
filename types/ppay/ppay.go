package ppaytype

import (
	"mime/multipart"
	"time"
)

type PPayOrderEntity struct {
	ID            int
	PayOrderNo    string
	OutTradeNo    string
	OrderAmount   string
	PayReason     string
	NotifyUrl     string
	PayAmount     string
	Payer         string
	ReceiptCodeId int
	PayChannel    string
	PayStatus     string
	StartTime     time.Time
	FinishTime    time.Time
}

type OrderDetail struct {
	*PPayOrderEntity
	ReceiptCodeUrl string
	*PReceiptAccountEntity
}

type CreateOrderReq struct {
	OutTradeNo     string
	OrderAmount    string
	PayAmount      string
	PayChannel     string
	PayeeAccountId int
	PayReason      string
	NotifyUrl      string
}

type CreateOrderRes struct {
	PayOrderNo string
}

type PReceiptCodeEntity struct {
	ID             int
	ReceiptCodeUrl string
	PayeeAmount    string
	PayeeAccountId int
	Operator       string
}

type PReceiptAccountEntity struct {
	ID          int
	AccountType string
	AccountName string
	NickName    string
	RealName    string
}

type CashierReq struct {
	PayOrderNo  string `json:"payOrderNo"`
	OutOrderNo  string `json:"outOrderNo"`
	OrderAmount string `json:"orderAmount"`
	PayAmount   string `json:"payAmount"`
	PayChannel  string `json:"payChannel"`
	PayReason   string `json:"payReason"`
}

type CashierRes struct {
	OrderStatus string `json:"orderStatus"`
	MaxTimeout  string `json:"maxTimeout"`
	QrCodeUrl   string `json:"qrCodeUrl"`
}

type UploadReceiptCodeReq struct {
	Operator       string
	PayeeAccountId int
	PayeeAmount    string
	multipart.File
	*multipart.FileHeader
}

type UploadReceiptCodeRes struct {
	ReceiptCodeId  int
	ReceiptCodeUrl string
}

type PayResultNotifyReq struct {
	PayChannel     string
	PayeeAccountId int
	Payer          string
	PayAmount      string
	SuccessTime    string
}

type PayResultNotifyRes struct {
	PayStatus string
}
