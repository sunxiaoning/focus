package ppaytype

import "time"

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
	ReceiptCode    string
	ReceiptCodeUrl string
	PayeeChannel   string
	PayeeAmount    string
	PayeeAccountId int
}

type PReceiptAccountEntity struct {
	ID          int
	AccountType string
	AccountName string
	NickName    string
	RealName    string
}
