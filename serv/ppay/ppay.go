package ppayserv

import (
	"context"
	"encoding/json"
	"fmt"
	"focus/cfg"
	ppaynotifyrepo "focus/repo/ppaynotify"
	ppayorderrepo "focus/repo/ppayorder"
	"focus/repo/preceiptaccount"
	preceiptcoderepo "focus/repo/preceiptcode"
	"focus/tx"
	"focus/types"
	notifystatusconst "focus/types/consts/notifystatus"
	orderstatusconst "focus/types/consts/orderstatus"
	userconsts "focus/types/consts/user"
	ppaytype "focus/types/ppay"
	dbutil "focus/util/db"
	fileutil "focus/util/file"
	httputil "focus/util/http"
	"focus/util/idgen"
	strutil "focus/util/strs"
	timutil "focus/util/tim"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
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
	Success             = "success"
)

var ReceiptCodeMimeTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
}

// 生成支付订单
func CreateOrder(ctx context.Context, reqParam *ppaytype.CreateOrderReq) *ppaytype.CreateOrderRes {
	receiptCodeEntity, payOrder := createOrderValidation(ctx, reqParam)
	payOrderNo, err := idgenutil.NextId()
	if err != nil {
		types.SystemPanic("generate payOrderNo failed!")
	}
	payOrder = &ppaytype.PPayOrderEntity{
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
	ppayorderrepo.Create(ctx, payOrder)
	res := &ppaytype.CreateOrderRes{PayOrderNo: payOrder.PayOrderNo}
	return res
}

func createOrderValidation(ctx context.Context, reqParam *ppaytype.CreateOrderReq) (*ppaytype.PReceiptCodeEntity, *ppaytype.PPayOrderEntity) {
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
	receiptAccountEntity := preceiptaccountrepo.GetById(ctx, reqParam.PayeeAccountId)
	if receiptAccountEntity.ID == 0 {
		types.NotFoundPanic(fmt.Sprintf("receiptAccount id =%d not exists!", reqParam.PayeeAccountId))
	}
	receiptCodeEntity := preceiptcoderepo.GetByAccountIdAndAmount(ctx, reqParam.PayAmount, reqParam.PayeeAccountId)
	if receiptCodeEntity.ID == 0 {

		// 兜底code
		receiptCodeEntity = preceiptcoderepo.GetByAccountIdAndAmount(ctx, preceiptcoderepo.MaxAmount, reqParam.PayeeAccountId)
		if receiptCodeEntity.ID == 0 {
			types.NotFoundPanic("not find PayeeAccount qrcode!")
		}
	}
	payOrder := ppayorderrepo.GetByOutTradeNo(ctx, reqParam.OutTradeNo)
	if payOrder.ID != 0 {
		types.RepeatRequestPanic("payOrder already exists!")
	}
	return receiptCodeEntity, payOrder
}

// 根据订单号查询支付订单详情
func GetOrderDetailByOrderNo(ctx context.Context, payOrderNo string) *ppaytype.OrderDetail {
	payOrderEntity := ppayorderrepo.GetByOrderNo(ctx, payOrderNo)
	receiptCodeEntity := preceiptcoderepo.GetById(ctx, payOrderEntity.ReceiptCodeId)
	if receiptCodeEntity.ID == 0 {
		types.ErrPanic(types.DataDirty, fmt.Sprintf("receiptcode id=%v not exists!", payOrderEntity.ReceiptCodeId))
	}
	receiptAccountEntity := preceiptaccountrepo.GetById(ctx, receiptCodeEntity.PayeeAccountId)
	if receiptAccountEntity.ID == 0 {
		types.ErrPanic(types.DataDirty, fmt.Sprintf("receiptaccount id=%v not exists!", receiptCodeEntity.PayeeAccountId))
	}
	return &ppaytype.OrderDetail{
		PPayOrderEntity:       payOrderEntity,
		ReceiptCodeUrl:        receiptCodeEntity.ReceiptCodeUrl,
		PReceiptAccountEntity: receiptAccountEntity,
	}
}

// 支付收银台
func Cashier(ctx context.Context, reqParam *ppaytype.CashierReq) tx.TFunRes {
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
	payOrderDetail := GetOrderDetailByOrderNo(ctx, reqParam.PayOrderNo)
	checkReqParamVal(reqParam, payOrderDetail)
	if payOrderDetail.PayStatus != orderstatusconst.I && payOrderDetail.PayStatus != orderstatusconst.P {
		return &ppaytype.CashierRes{
			OrderStatus: payOrderDetail.PayStatus,
		}
	}
	payOrderUpdated := submitPayOrder(ctx, reqParam.PayOrderNo)
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

func submitPayOrder(ctx context.Context, payOrderNo string) *ppaytype.PPayOrderEntity {
	updated := ppayorderrepo.Submit(ctx, payOrderNo)
	logrus.Infof("update payOrder count=%v", updated)
	return ppayorderrepo.GetByOrderNo(ctx, payOrderNo)
}

func GetReceiptCode(ctx context.Context, qrCodeUrl string) []byte {
	if !fileutil.PathExist(path.Join(cfg.Cfg.Server.RootFilePath, ReceiptCodeRootPath, qrCodeUrl)) {
		types.NotFoundPanic(fmt.Sprintf("%s not find!", qrCodeUrl))
	}
	content, err := ioutil.ReadFile(path.Join(cfg.Cfg.Server.RootFilePath, ReceiptCodeRootPath, qrCodeUrl))
	if err != nil {
		types.SystemPanic("read qrcode file err!")
	}
	return content
}

func UploadReceiptCode(ctx context.Context, req *ppaytype.UploadReceiptCodeReq) *ppaytype.UploadReceiptCodeRes {
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
	receiptCodeEntity := preceiptcoderepo.GetByAccountIdAndAmount(ctx, req.PayeeAmount, req.PayeeAccountId)
	if receiptCodeEntity.ID == 0 {
		receiptCodeEntity = &ppaytype.PReceiptCodeEntity{
			ReceiptCodeUrl: path.Join(savePath, saveName),
			PayeeAmount:    req.PayeeAmount,
			PayeeAccountId: req.PayeeAccountId,
			Operator:       req.Operator,
		}
		preceiptcoderepo.Create(ctx, receiptCodeEntity)
	} else {
		preceiptcoderepo.UpdateReceiptCodeUrl(ctx, receiptCodeEntity.ID, path.Join(savePath, saveName))
	}
	return &ppaytype.UploadReceiptCodeRes{
		ReceiptCodeId:  receiptCodeEntity.ID,
		ReceiptCodeUrl: receiptCodeEntity.ReceiptCodeUrl,
	}
}

func ResultNotifyTx(ctx context.Context) tx.TFunRes {
	payNotifyReq := ctx.Value("payNotifyReq").(*ppaytype.PayResultNotifyReq)
	if strutil.IsBlank(payNotifyReq.PayChannel) {
		types.InvalidParamPanic("payChannel can't be empty!")
	}
	if types.PayChannels()[payNotifyReq.PayChannel] == nil {
		types.ErrPanic(types.PayChannelNotSupport, fmt.Sprintf("payChannel=%s not Support!", payNotifyReq.PayChannel))
	}
	if payNotifyReq.PayeeAccountId <= 0 {
		types.InvalidParamPanic("payeeAccountId is invalid!")
	}
	if strutil.IsBlank(payNotifyReq.PayAmount) {
		types.InvalidParamPanic("payAmount can't be empty!")
	}
	if !strutil.IsValidMoney(payNotifyReq.PayAmount) {
		types.InvalidParamPanic("payAmount is invalid!")
	}
	if strutil.IsBlank(payNotifyReq.SuccessTime) {
		types.InvalidParamPanic("successTime can't be empty!")
	}
	if tim, err := timutil.DefParse(payNotifyReq.SuccessTime); err != nil || tim.After(time.Now()) {
		types.InvalidParamPanic("successTime is invalid!")
	}
	logrus.Infof("payNotifyReq: %v", payNotifyReq)
	receiptAccountEntity := preceiptcoderepo.GetById(ctx, payNotifyReq.PayeeAccountId)
	if receiptAccountEntity.ID == 0 {
		types.NotFoundPanic(fmt.Sprintf("receiptAccount id =%s not exists!", payNotifyReq.PayeeAccountId))
	}
	payOrderEntity := ppayorderrepo.GetWaittingByPayChannelAndAmount(ctx, payNotifyReq.PayAmount, payNotifyReq.PayChannel)
	if payOrderEntity.ID == 0 {
		types.NotFoundPanic("payOrder not exists!")
	}
	receiptCodeEntity := preceiptcoderepo.GetById(ctx, payOrderEntity.ReceiptCodeId)
	if receiptCodeEntity.ID == 0 || receiptCodeEntity.PayeeAccountId != payNotifyReq.PayeeAccountId {
		types.NotFoundPanic("payOrder not exists!")
	}
	result := &ppaytype.PayResultNotifyRes{
		PayStatus: orderstatusconst.S,
	}
	updateResult := ppayorderrepo.UpdatePayResult(ctx, payOrderEntity.ID, orderstatusconst.S)
	if updateResult != 1 {
		logrus.Infof("payOrder payOrderNo=%s has been processed!", payOrderEntity.PayOrderNo)
		return result
	}
	payOrderEntity.PayStatus = orderstatusconst.S
	InsertPayResultNotify(ctx, payOrderEntity)
	return result
}

func InsertPayResultNotify(ctx context.Context, payOrderEntity *ppaytype.PPayOrderEntity) {

	// 发送消息通知业务方支付成功
	notifyContent := ppaytype.BizPayResultReq{
		PayOrderNo:  payOrderEntity.PayOrderNo,
		PayReason:   payOrderEntity.PayReason,
		OrderAmount: payOrderEntity.OrderAmount,
		PayAmount:   payOrderEntity.PayAmount,
		PayStatus:   payOrderEntity.PayStatus,
	}
	notifyContentResult, _ := json.Marshal(notifyContent)
	ppayNotifyEntity := &ppaytype.PPayNotifyEntity{
		NotifyUrl:     payOrderEntity.NotifyUrl,
		NotifyStatus:  notifystatusconst.I,
		NotifyContent: string(notifyContentResult),
		CreatedTime:   time.Now(),
	}
	ppaynotifyrepo.Create(ctx, ppayNotifyEntity)
}

func NotifyBiz() {
	var payNotifyResults []*ppaytype.PPayNotifyEntity
	dbutil.NewDbExecutor(cfg.FocusCtx.DB.Table("personal_pay_notify").Where("notify_status = ? and status = 1", notifystatusconst.I).Find(&payNotifyResults))
	if payNotifyResults == nil || len(payNotifyResults) == 0 {
		logrus.Info("payNotifyResult is empty,end!")
		return
	}
	for _, payNotifyResult := range payNotifyResults {
		updateResult := dbutil.NewDbExecutor(cfg.FocusCtx.DB.Table("personal_pay_notify").Where("id = ? and notify_status = ? and status = 1", payNotifyResult.ID, notifystatusconst.I).Update("notify_status", notifystatusconst.P)).RowsAffected()
		if updateResult != 1 {
			logrus.Infof("payNotify ID = ? has been processed!", payNotifyResult.ID)
			return
		}
		reqHeaders := make(map[string]string)
		reqHeaders[userconsts.AccessToken] = "JRCASt7GYl0d5g5OAKFgiA=="
		code, content, _ := httputil.PostJsonWithHeader(payNotifyResult.NotifyUrl, reqHeaders, payNotifyResult.NotifyContent, time.Second*10)
		notifyStatus := notifystatusconst.I

		// 回调成功，统一返回 success
		if code == http.StatusOK && string(content) == Success {
			notifyStatus = notifystatusconst.S
		} else if time.Now().Add(time.Hour * 2).Before(payNotifyResult.CreatedTime) {
			logrus.Warn("notify biz timeout, stop notify!")
			notifyStatus = "F"
		}
		dbutil.NewDbExecutor(cfg.FocusCtx.DB.Table("personal_pay_notify").Where("id = ? and notify_status = ? and status = 1", payNotifyResult.ID, orderstatusconst.P).Update("notify_status", notifyStatus)).RowsAffected()
	}
}
