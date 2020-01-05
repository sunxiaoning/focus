package ppayserv

import (
	"context"
	"fmt"
	"focus/cfg"
	"focus/tx"
	"focus/types"
	orderstatusconst "focus/types/consts/orderstatus"
	ppaytype "focus/types/ppay"
	dbutil "focus/util/db"
	fileutil "focus/util/file"
	"focus/util/idgen"
	strutil "focus/util/strs"
	timutil "focus/util/tim"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

const (
	PayOrderTimeOut     = 5
	MaxReceiptCodeSize  = 5 * 1024 * 1024
	ReceiptCodeRootPath = "/ReceiptCodes"
)

var ReceiptCodeMimeTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
}

func CreateOrder(ctx context.Context) *ppaytype.CreateOrderRes {
	tx, ok := ctx.Value("tx").(*gorm.DB)
	if !ok {
		tx = cfg.FocusCtx.DB
	}
	reqParam := ctx.Value("reqParam").(*ppaytype.CreateOrderReq)
	if strutil.IsBlank(reqParam.OutTradeNo) {
		types.InvalidParamPanic("outTradeNo can't be empty!")
	}
	if strutil.IsBlank(reqParam.OrderAmount) {
		types.InvalidParamPanic("orderAmount can't be empty!")
	}
	if !strutil.IsValidMoney(reqParam.OrderAmount) {
		types.InvalidParamPanic("orderAmount is invalid!")
	}
	if strutil.IsBlank(reqParam.PayAmount) {
		types.InvalidParamPanic("payAmount can't be empty!")
	}
	if !strutil.IsValidMoney(reqParam.PayAmount) {
		types.InvalidParamPanic("payAmount is invalid!")
	}
	if strutil.IsBlank(reqParam.PayChannel) {
		types.InvalidParamPanic("payChannel can't be empty!")
	}
	if types.PayChannels()[reqParam.PayChannel] == nil {
		types.ErrPanic(types.PayChannelNotSupport, fmt.Sprintf("payChannel=%s not Support!", reqParam.PayChannel))
	}
	if reqParam.PayeeAccountId <= 0 {
		types.ErrPanic(types.InvalidParamError, "payeeAccountId is invalid!")
	}
	if strutil.IsBlank(reqParam.PayReason) {
		types.InvalidParamPanic("payReason can't be empty!")
	}
	if strutil.IsBlank(reqParam.NotifyUrl) {
		types.InvalidParamPanic("notifyUrl can't be empty!")
	}
	var receiptAccountEntity ppaytype.PReceiptAccountEntity
	dbutil.NewDbExecutor(tx.Table("personal_receipt_account").Where("id = ? and status = 1", reqParam.PayeeAccountId).Find(&receiptAccountEntity))
	if receiptAccountEntity.ID == 0 {
		types.NotFoundPanic(fmt.Sprintf("receiptAccount id =%s not exists!", reqParam.PayeeAccountId))
	}
	var receiptCodeEntity ppaytype.PReceiptCodeEntity
	dbutil.NewDbExecutor(tx.Table("personal_receipt_code").Where("payee_channel = ? and payee_amount = ? and payee_account_id = ? and status = 1", reqParam.PayChannel, reqParam.PayAmount, reqParam.PayeeAccountId).Find(&receiptCodeEntity))
	if receiptCodeEntity.ID == 0 {
		dbutil.NewDbExecutor(tx.Table("personal_receipt_code").Where(`payee_channel = ? and payee_amount = 9999.99
		and payee_account_id = ? and status = 1`, reqParam.PayChannel, reqParam.PayeeAccountId).Find(&receiptCodeEntity))
	}
	if receiptCodeEntity.ID == 0 {
		types.NotFoundPanic("not find PayeeAccount qrcode!")
	}
	var payOrder ppaytype.PPayOrderEntity
	dbutil.NewDbExecutor(tx.Table("personal_pay_order").Where("out_trade_no = ? and status = 1", reqParam.OutTradeNo).Find(&payOrder))
	if payOrder.ID != 0 {
		types.RepeatRequestPanic("payOrder already exists!")
	}
	payOrderNo, err := idgenutil.NextId()
	if err != nil {
		types.SystemPanic("generate payOrderNo failed!")
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
	return res
}

func GetOrderDetailByOrderNo(ctx context.Context) *ppaytype.OrderDetail {
	payOrderNo := ctx.Value("payOrderNo").(string)
	payOrderEntity := getPayOrderEntity(payOrderNo)
	var receiptCodeEntity ppaytype.PReceiptCodeEntity
	dbutil.NewDbExecutor(cfg.FocusCtx.DB.Table("personal_receipt_code").Where("id = ? and status = 1", payOrderEntity.ReceiptCodeId).Find(&receiptCodeEntity))
	if receiptCodeEntity.ID == 0 {
		types.ErrPanic(types.DataDirty, fmt.Sprintf("receiptcode id=%v not exists!", payOrderEntity.ReceiptCodeId))
	}
	var receiptAccountEntity ppaytype.PReceiptAccountEntity
	dbutil.NewDbExecutor(cfg.FocusCtx.DB.Table("personal_receipt_account").Where("id = ? and status = 1", receiptCodeEntity.PayeeAccountId).Find(&receiptAccountEntity))
	if receiptAccountEntity.ID == 0 {
		types.ErrPanic(types.DataDirty, fmt.Sprintf("receiptaccount id=%v not exists!", receiptCodeEntity.PayeeAccountId))
	}
	return &ppaytype.OrderDetail{
		PPayOrderEntity:       payOrderEntity,
		ReceiptCodeUrl:        receiptCodeEntity.ReceiptCodeUrl,
		PReceiptAccountEntity: &receiptAccountEntity,
	}
}

func getPayOrderEntity(payOrderNo string) *ppaytype.PPayOrderEntity {
	var payOrderEntity ppaytype.PPayOrderEntity
	dbutil.NewDbExecutor(cfg.FocusCtx.DB.Table("personal_pay_order").Where("pay_order_no = ? and status = 1", payOrderNo).Find(&payOrderEntity))
	if payOrderEntity.ID == 0 {
		types.InvalidParamPanic(fmt.Sprintf("payOrder no=%v not exists!", payOrderNo))
	}
	return &payOrderEntity
}

func Cashier(ctx context.Context) tx.TFunRes {
	reqParam := ctx.Value("reqParam").(*ppaytype.CashierReq)
	if strutil.IsBlank(reqParam.PayOrderNo) {
		types.InvalidParamPanic("payOrderNo can't be empty!")
	}
	if strutil.IsBlank(reqParam.OutOrderNo) {
		types.InvalidParamPanic("serviceOrderNo can't be empty!")
	}
	if strutil.IsBlank(reqParam.OrderAmount) {
		types.InvalidParamPanic("orderAmount can't be empty!")
	}
	if !strutil.IsValidMoney(reqParam.OrderAmount) {
		types.InvalidParamPanic("orderAmount is invalid!")
	}
	if strutil.IsBlank(reqParam.PayAmount) {
		types.InvalidParamPanic("payOrderNo can't be empty!")
	}
	if !strutil.IsValidMoney(reqParam.PayAmount) {
		types.InvalidParamPanic("payAmount can't be empty!")
	}
	if strutil.IsBlank(reqParam.PayChannel) {
		types.InvalidParamPanic("payChannel can't be empty!")
	}
	if types.PayChannels()[reqParam.PayChannel] == nil {
		types.ErrPanic(types.PayChannelNotSupport, "payChannel not support!")
	}
	ctx = context.WithValue(ctx, "payOrderNo", reqParam.PayOrderNo)
	payOrderDetail := GetOrderDetailByOrderNo(ctx)
	checkReqParamVal(reqParam, payOrderDetail)
	if payOrderDetail.PayStatus != orderstatusconst.I && payOrderDetail.PayStatus != orderstatusconst.P {
		return &ppaytype.CashierRes{
			OrderStatus: payOrderDetail.PayStatus,
		}
	}
	payOrderUpdated := SubmitPayOrder(context.WithValue(ctx, "payOrderNo", reqParam.PayOrderNo))
	return &ppaytype.CashierRes{
		OrderStatus: payOrderUpdated.PayStatus,
		MaxTimeout:  timutil.DefFormat(payOrderDetail.StartTime.Add(time.Minute * PayOrderTimeOut)),
		QrCodeUrl:   payOrderDetail.ReceiptCodeUrl,
	}
}

func checkReqParamVal(param *ppaytype.CashierReq, order *ppaytype.OrderDetail) {
	isLegal := param.PayAmount == order.PayAmount && param.OrderAmount == order.OrderAmount && param.OutOrderNo == order.OutTradeNo && param.PayChannel == order.PayChannel && param.PayReason == order.PayReason
	if !isLegal {
		types.InvalidParamPanic("invalid cashier req!")
	}
}

func SubmitPayOrder(ctx context.Context) *ppaytype.PPayOrderEntity {
	payOrderNo := ctx.Value("payOrderNo").(string)
	payOrderEntity := getPayOrderEntity(payOrderNo)
	updated := dbutil.NewDbExecutor(cfg.FocusCtx.DB.Table("personal_pay_order").Where("pay_order_no = ? and pay_status = 'I' and status = 1", payOrderEntity.PayOrderNo).Update("pay_status", orderstatusconst.P)).RowsAffected()
	logrus.Infof("update payOrder count=%v", updated)
	return getPayOrderEntity(payOrderNo)
}

func GetReceiptCode(ctx context.Context) []byte {
	qrCodeUrl := ctx.Value("qrCodeUrl").(string)
	if !fileutil.PathExist(path.Join(cfg.Cfg.Server.RootFilePath, ReceiptCodeRootPath, qrCodeUrl)) {
		types.NotFoundPanic(fmt.Sprintf("%s not find!", qrCodeUrl))
	}
	content, err := ioutil.ReadFile(path.Join(cfg.Cfg.Server.RootFilePath, ReceiptCodeRootPath, qrCodeUrl))
	if err != nil {
		types.SystemPanic("read qrcode file err!")
	}
	return content
}

func UploadReceiptCode(ctx context.Context) *ppaytype.UploadReceiptCodeRes {
	req := ctx.Value("UploadReceiptCodeReq").(*ppaytype.UploadReceiptCodeReq)
	var destFile *os.File
	defer func() {
		req.File.Close()
		if r := recover(); r != nil {
			if destFile != nil {
				os.Remove(destFile.Name())
			}
			panic(r)
		}
	}()
	if req.FileHeader.Size > MaxReceiptCodeSize {
		types.ErrPanic(types.FileSizeTooLarge, "ReceiptCode Size too Large!")
	}
	if !ReceiptCodeMimeTypes[req.FileHeader.Header.Get("Content-Type")] {
		types.InvalidParamPanic("ReceiptCode mimetype not support!")
	}
	t := time.Now()
	savePath := path.Join(strconv.Itoa(t.Year()), strconv.Itoa(int(t.Month())), strconv.Itoa(t.Day()))
	id, err := idgenutil.NextId()
	if err != nil {
		types.ErrPanic(types.GenUUIDError, err.Error())
	}
	saveName := strings.Join([]string{
		fmt.Sprintf("%v", id),
		req.FileHeader.Filename[strings.LastIndex(req.FileHeader.Filename, ".")+1:],
	}, ".")
	destFile, err = fileutil.OpenFile(path.Join(cfg.Cfg.Server.RootFilePath, ReceiptCodeRootPath, savePath, saveName), os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		types.SystemPanic(fmt.Sprintf("save receiptcode err! reason: %v", err))
	}
	defer destFile.Close()
	_, err = io.Copy(destFile, req.File)
	if err != nil {
		types.SystemPanic(fmt.Sprintf("save receiptcode err! reason: %v", err))
	}
	var receiptCodeEntity ppaytype.PReceiptCodeEntity
	dbutil.NewDbExecutor(cfg.FocusCtx.DB.Table("personal_receipt_code").Where("payee_account_id = ? and payee_amount = ? and status = 1", req.PayeeAccountId, req.PayeeAmount).Find(&receiptCodeEntity))
	if receiptCodeEntity.ID == 0 {
		receiptCodeEntity = ppaytype.PReceiptCodeEntity{
			ReceiptCodeUrl: path.Join(savePath, saveName),
			PayeeAmount:    req.PayeeAmount,
			PayeeAccountId: req.PayeeAccountId,
			Operator:       req.Operator,
		}
		dbutil.NewDbExecutor(cfg.FocusCtx.DB.Table("personal_receipt_code").Create(&receiptCodeEntity))
	} else {
		dbutil.NewDbExecutor(cfg.FocusCtx.DB.Table("personal_receipt_code").Where("id = ? and status = 1", receiptCodeEntity.ID).Update("receipt_code_url", path.Join(savePath, saveName)))
	}
	return &ppaytype.UploadReceiptCodeRes{
		ReceiptCodeId:  receiptCodeEntity.ID,
		ReceiptCodeUrl: receiptCodeEntity.ReceiptCodeUrl,
	}
}
