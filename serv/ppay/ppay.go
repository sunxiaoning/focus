package ppayserv

import (
	"context"
	"fmt"
	"focus/cfg"
	"focus/types"
	orderstatusconst "focus/types/consts/orderstatus"
	ppaytype "focus/types/ppay"
	"focus/util"
	dbutil "focus/util/db"
	strutil "focus/util/strs"
	timutil "focus/util/tim"
	"github.com/jinzhu/gorm"
	"time"
)

func CreateOrder(ctx context.Context) (*ppaytype.CreateOrderRes, error) {
	tx, ok := ctx.Value("tx").(*gorm.DB)
	if !ok {
		tx = cfg.FocusCtx.DB
	}
	reqParam := ctx.Value("reqParam").(*ppaytype.CreateOrderReq)
	if strutil.IsBlank(reqParam.OutTradeNo) {
		return nil, types.InvalidParamErr("outTradeNo can't be empty!")
	}
	if strutil.IsBlank(reqParam.OrderAmount) {
		return nil, types.InvalidParamErr("orderAmount can't be empty!")
	}
	if !strutil.IsValidMoney(reqParam.OrderAmount) {
		return nil, types.InvalidParamErr("orderAmount is invalid!")
	}
	if strutil.IsBlank(reqParam.PayAmount) {
		return nil, types.InvalidParamErr("payAmount can't be empty!")
	}
	if !strutil.IsValidMoney(reqParam.PayAmount) {
		return nil, types.InvalidParamErr("payAmount is invalid!")
	}
	if strutil.IsBlank(reqParam.PayChannel) {
		return nil, types.InvalidParamErr("payChannel can't be empty!")
	}
	if types.PayChannels()[reqParam.PayChannel] == nil {
		return nil, types.NewErr(types.PayChannelNotSupport, fmt.Sprintf("payChannel=%s not Support!", reqParam.PayChannel))
	}
	if reqParam.PayeeAccountId <= 0 {
		return nil, types.NewErr(types.InvalidParamError, "payeeAccountId is invalid!")
	}
	if strutil.IsBlank(reqParam.PayReason) {
		return nil, types.InvalidParamErr("payReason can't be empty!")
	}
	if strutil.IsBlank(reqParam.NotifyUrl) {
		return nil, types.InvalidParamErr("notifyUrl can't be empty!")
	}
	var receiptAccountEntity ppaytype.PReceiptAccountEntity
	dbutil.NewDbExecutor(tx.Table("personal_receipt_account").Where("id = ? and status = 1", reqParam.PayeeAccountId).Find(&receiptAccountEntity))
	if receiptAccountEntity.ID == 0 {
		return nil, types.NotFoundErr(fmt.Sprintf("receiptAccount id =%s not exists!", reqParam.PayeeAccountId))
	}
	var receiptCodeEntity ppaytype.PReceiptCodeEntity
	dbutil.NewDbExecutor(tx.Table("personal_receipt_code").Where("payee_channel = ? and payee_amount = ? and payee_account_id = ? and status = 1", reqParam.PayChannel, reqParam.PayAmount, reqParam.PayeeAccountId).Find(&receiptCodeEntity))
	if receiptCodeEntity.ID == 0 {
		dbutil.NewDbExecutor(tx.Table("personal_receipt_code").Where(`payee_channel = ? and payee_amount = 9999.99
		and payee_account_id = ? and status = 1`, reqParam.PayChannel, reqParam.PayeeAccountId).Find(&receiptCodeEntity))
	}
	if receiptCodeEntity.ID == 0 {
		return nil, types.NotFoundErr("not find PayeeAccount qrcode!")
	}
	var payOrder ppaytype.PPayOrderEntity
	dbutil.NewDbExecutor(tx.Table("personal_pay_order").Where("out_trade_no = ? and status = 1", reqParam.OutTradeNo).Find(&payOrder))
	if payOrder.ID != 0 {
		return nil, types.RepeatRequestErr("payOrder already exists!")
	}
	payOrderNo, err := util.IdGenerator.NextID()
	if err != nil {
		return nil, types.SystemErr("generate payOrderNo failed!")
	}
	payOrder = ppaytype.PPayOrderEntity{
		PayOrderNo:    fmt.Sprintf("%v", payOrderNo),
		OutTradeNo:    reqParam.OutTradeNo,
		OrderAmount:   reqParam.OrderAmount,
		PayReason:     reqParam.PayReason,
		NotifyUrl:     reqParam.NotifyUrl,
		PayAmount:     reqParam.PayAmount,
		ReceiptCodeId: receiptCodeEntity.ID,
		PayChannel:    reqParam.PayChannel,
		PayStatus:     orderstatusconst.I,
		StartTime:     time.Now(),
		FinishTime:    timutil.ZERO,
	}
	dbutil.NewDbExecutor(tx.Table("personal_pay_order").Create(&payOrder))
	res := &ppaytype.CreateOrderRes{PayOrderNo: payOrder.PayOrderNo}
	return res, nil
}
