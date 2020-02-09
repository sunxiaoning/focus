package servserv

import (
	"context"
	"fmt"
	"focus/cfg"
	servrepo "focus/repo/serv"
	ppayserv "focus/serv/ppay"
	"focus/tx"
	"focus/types"
	"focus/types/consts/orderstatus"
	usertype "focus/types/member"
	pagetype "focus/types/page"
	ppaytype "focus/types/ppay"
	servicetype "focus/types/service"
	dbutil "focus/util/db"
	strutil "focus/util/strs"
	timutil "focus/util/tim"
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	CashierUrl = "/api/v1/service/cashier"
	Success    = "success"
)

// 最新服务查询
func QueryLatest(ctx context.Context) *types.PageResponse {
	reqParam := ctx.Value("reqParam").(*servicetype.QueryLatestReq)
	if reqParam.PageIndex < 1 {
		types.InvalidParamPanic("pageIndex is invalid!")
	}
	if reqParam.PageSize < 1 || reqParam.PageSize > 1000 {
		types.InvalidParamPanic("pageSize is invalid!")
	}
	queryParams := map[string]interface{}{}
	if reqParam.ChineseName != "" {
		queryParams["chinese_name"] = reqParam.ChineseName
	}
	if reqParam.ServiceType != 0 {
		queryParams["service_type"] = reqParam.ServiceType
	}
	pageQuery := pagetype.PageQuery{
		Page:   pagetype.NewPage(reqParam.PageIndex, reqParam.PageSize),
		Params: queryParams,
	}
	var services []*servicetype.ServiceEntity
	var total int
	servrepo.QueryLatest(pageQuery, &services, &total)
	var results []*servicetype.QueryLatestRes
	for _, service := range services {
		result := &servicetype.QueryLatestRes{ServiceId: service.ID, ServiceType: service.ServiceType, ChineseName: service.ChineseName, ServiceDesc: service.ServiceDesc, PublishTime: timutil.DefFormat(service.PublishTime)}
		results = append(results, result)
	}
	return types.NewPageResponse(total, results)
}

// 服务详情查询
func GetServiceById(ctx context.Context) *servicetype.GetByIdRes {
	serviceId := ctx.Value("serviceId").(int)
	if serviceId <= 0 {
		types.InvalidParamPanic("serviceId is invalided!")
	}
	service := &servicetype.ServiceEntity{}
	dbutil.NewDbExecutor(cfg.FocusCtx.DB.Table("service").Where("id = ? and service_status = 'FWZ' and status = 1", serviceId).Find(service))
	if service.ID == 0 {
		types.NotFoundPanic("service not exists!")
	}
	return &servicetype.GetByIdRes{ServiceId: service.ID, ServiceType: service.ServiceType, ChineseName: service.ChineseName, ServiceDesc: service.ServiceDesc, PublishTime: timutil.DefFormat(service.PublishTime)}
}

// 服务套餐查询
func QueryPrice(ctx context.Context) []*servicetype.QueryPriceRes {
	serviceId := ctx.Value("serviceId").(int)
	if serviceId <= 0 {
		types.InvalidParamPanic("serviceId is invalided!")
	}
	service := GetServiceById(ctx)
	if service == nil {
		types.NotFoundPanic("service not found!")
	}
	var prices []*servicetype.PriceEntity
	dbutil.NewDbExecutor(cfg.FocusCtx.DB.Table("service_price").Where("service_id = ? and status = 1", serviceId).Find(&prices))
	var results []*servicetype.QueryPriceRes
	for _, price := range prices {
		result := &servicetype.QueryPriceRes{ID: price.ID, PriceName: price.Price, ServiceId: price.ServiceId, Price: price.Price, ServiceAmount: price.ServiceAmount}
		results = append(results, result)
	}
	return results
}

// 服务价格计算器
func CalculatePrice(ctx context.Context) *servicetype.CalculatePriceRes {
	reqParam := ctx.Value("reqParam").(*servicetype.CalculatePriceReq)
	if reqParam.PriceId <= 0 {
		types.InvalidParamPanic("priceId param is invalid!")
	}
	if reqParam.Amount <= 0 {
		types.InvalidParamPanic("amount param is invalid!")
	}
	var priceEntity servicetype.PriceEntity
	dbutil.NewDbExecutor(cfg.FocusCtx.DB.Table("service_price").Where("id = ? and status = 1", reqParam.PriceId).First(&priceEntity))
	if priceEntity.ID == 0 {
		types.NotFoundPanic("price not exists!")
	}
	price, err := decimal.NewFromString(priceEntity.Price)
	if err != nil {
		types.ErrPanic(types.DataDirty, fmt.Sprintf("invalid price=%s", priceEntity.Price))
	}
	decimal.DivisionPrecision = 2
	price = price.Mul(decimal.NewFromInt(reqParam.Amount))
	return &servicetype.CalculatePriceRes{Price: price.StringFixedBank(2)}
}

// 生成服务订单
func CreateOrderTx(ctx context.Context) tx.TFunRes {
	tx := ctx.Value("tx").(*gorm.DB)
	reqParam := ctx.Value("reqParam").(*servicetype.CreateOrderRequest)
	if reqParam.OrderNo == "" {
		types.InvalidParamPanic("orderNo can't be empty!")
	}
	if reqParam.MemberId <= 0 {
		types.InvalidParamPanic("memberId is invalid!")
	}
	if reqParam.ServicePriceId <= 0 {
		types.InvalidParamPanic("servicePriceId is invalid!")
	}
	if reqParam.PurchaseAmount <= 1 {
		types.InvalidParamPanic("purchaseAmount is invalid!!")
	}
	if strutil.IsBlank(reqParam.PayChannel) {
		types.InvalidParamPanic("payChannel can't be empty!")
	}
	if types.PayChannels()[reqParam.PayChannel] == nil {
		types.ErrPanic(types.PayChannelNotSupport, "payChannel not support!")
	}
	var memberEntity usertype.MemberEntity
	dbutil.NewDbExecutor(cfg.FocusCtx.DB.Table("member").Where("id = ? and status = 1", reqParam.MemberId).Find(&memberEntity))
	if memberEntity.ID == 0 {
		types.NotFoundPanic(fmt.Sprintf("memberId %v not exists!", reqParam.MemberId))
	}
	var orderEntity servicetype.OrderEntity
	dbutil.NewDbExecutor(tx.Table("service_order").Where("order_no = ? and status = 1", reqParam.OrderNo).Find(&orderEntity))
	if orderEntity.ID != 0 {
		types.RepeatRequestPanic(fmt.Sprintf("order orderNo=%s already exists!", reqParam.OrderNo))
	}
	var priceEntity servicetype.PriceEntity
	dbutil.NewDbExecutor(tx.Table("service_price").Where("id = ? and status = 1", reqParam.ServicePriceId).Find(&priceEntity))
	if priceEntity.ID == 0 {
		types.NotFoundPanic(fmt.Sprintf("price priceId=%d not exists!", reqParam.ServicePriceId))
	}
	var serviceEntity servicetype.ServiceEntity
	dbutil.NewDbExecutor(tx.Table("service").Where("id = ? and service_status = 'FWZ' and status = 1", priceEntity.ServiceId).Find(&serviceEntity))
	if serviceEntity.ID == 0 {
		types.NotFoundPanic(fmt.Sprintf("service %d not found!", priceEntity.ServiceId))
	}
	orderEntity.OrderNo = reqParam.OrderNo
	orderEntity.MemberId = reqParam.MemberId
	orderEntity.ServicePriceId = reqParam.ServicePriceId
	orderEntity.PurchaseAmount = reqParam.PurchaseAmount
	orderEntity.StartTime = time.Now()
	orderEntity.FinishedTime = timutil.ZERO
	orderEntity.OrderStatus = orderstatusconst.I
	price, err := decimal.NewFromString(priceEntity.Price)
	if err != nil {
		types.ErrPanic(types.DataDirty, fmt.Sprintf("invalid price=%s", priceEntity.Price))
	}
	price = price.Mul(decimal.NewFromInt(int64(reqParam.PurchaseAmount))).Round(2)
	orderEntity.OrderAmount = price.StringFixedBank(2)
	orderEntity.PayAmount = price.StringFixedBank(2)
	orderEntity.CouponNo = reqParam.CouponNo
	orderEntity.PayChannel = reqParam.PayChannel
	dbutil.NewDbExecutor(tx.Table("service_order").Create(&orderEntity))
	createPayOrderReq := &ppaytype.CreateOrderReq{
		OutTradeNo:     orderEntity.OrderNo,
		OrderAmount:    orderEntity.OrderAmount,
		PayAmount:      orderEntity.PayAmount,
		PayChannel:     orderEntity.PayChannel,
		PayeeAccountId: types.PayChannels()[orderEntity.PayChannel].AccountId,
		PayReason:      fmt.Sprintf("购买服务:%v,套餐:%v,数量：%v", serviceEntity.ChineseName, priceEntity.PriceName, reqParam.PurchaseAmount),
		NotifyUrl:      fmt.Sprintf("http://localhost:%d/api/v1/notify/ppay", cfg.FocusCtx.Cfg.Server.ListenPort),
	}
	ctx = context.WithValue(ctx, "reqParam", createPayOrderReq)
	createPayOrderRes := ppayserv.CreateOrder(ctx)
	dbutil.NewDbExecutor(tx.Table("service_order").Where("id = ? and status = 1", orderEntity.ID).Update("out_order_no", createPayOrderRes.PayOrderNo))
	cashierParams := &servicetype.CashierReq{
		PayOrderNo:     createPayOrderRes.PayOrderNo,
		ServiceOrderNo: orderEntity.OrderNo,
		OrderAmount:    createPayOrderReq.OrderAmount,
		PayAmount:      createPayOrderReq.PayAmount,
		PayChannel:     createPayOrderReq.PayChannel,
		PayReason:      createPayOrderReq.PayReason,
	}
	res := &servicetype.CreateOrderRes{
		CashierUrl:    CashierUrl,
		CashierParams: cashierParams,
	}
	return res
}

// 支付系统通知服务订单支付结果
func PayResultNotify(ctx context.Context) string {
	payResult := ctx.Value("payResult").(*ppaytype.BizPayResultReq)
	if strutil.IsBlank(payResult.PayStatus) || (payResult.PayStatus != orderstatusconst.S && payResult.PayStatus != orderstatusconst.F) {
		types.InvalidParamPanic("payStatus is invalid!")
	}
	if strutil.IsBlank(payResult.OrderAmount) || !strutil.IsValidMoney(payResult.OrderAmount) {
		types.InvalidParamPanic("orderAmount is invalid!")
	}
	if strutil.IsBlank(payResult.PayOrderNo) {
		types.InvalidParamPanic("payOrderNo is invalid!")
	}
	logrus.Infof("payResult:%v", payResult)
	return tx.NewTxManager().RunTx(ctx, puchaseServiceTx).(string)
}

// 支付成功，生成服务权益
func puchaseServiceTx(ctx context.Context) tx.TFunRes {
	tx := ctx.Value("tx").(*gorm.DB)
	payResult := ctx.Value("payResult").(*ppaytype.BizPayResultReq)
	var serviceOrderEntity servicetype.OrderEntity
	dbutil.NewDbExecutor(tx.Table("service_order").Where("out_order_no = ? and order_status = 'P' and status = 1", payResult.PayOrderNo).Find(&serviceOrderEntity))
	if serviceOrderEntity.ID == 0 {
		logrus.Warnf("serviceOrder not exists! outTradeNo:%s", payResult.PayOrderNo)
		return Success
	}
	updateResult := dbutil.NewDbExecutor(tx.Table("service_order").Where("out_order_no = ? and order_status = 'P' and status = 1", payResult.PayOrderNo).Updates(servicetype.OrderEntity{OrderStatus: payResult.PayStatus, FinishedTime: time.Now()})).RowsAffected()
	if updateResult != 1 {
		logrus.Warn("serviceOrder has been processed! outTradeNo:%s", payResult.PayOrderNo)
		return Success
	}
	if payResult.PayStatus == orderstatusconst.S {
		var memberServiceEntity servicetype.MemberServiceEntity
		dbutil.NewDbExecutor(tx.Table("member_service").Where("member_id = ? and service_price_id = ? and member_service_status = 1 and status = 1", serviceOrderEntity.MemberId, serviceOrderEntity.ServicePriceId).Find(&memberServiceEntity))
		if memberServiceEntity.ID == 0 {
			memberServiceEntity = servicetype.MemberServiceEntity{
				MemberId:            serviceOrderEntity.MemberId,
				ServicePriceId:      serviceOrderEntity.ServicePriceId,
				OrderId:             serviceOrderEntity.ID,
				DeadlineTime:        time.Now().AddDate(0, serviceOrderEntity.PurchaseAmount, 0),
				MemberServiceStatus: 1,
			}
			dbutil.NewDbExecutor(tx.Table("member_service").Create(&memberServiceEntity))
		} else {
			memberServiceEntity.OrderId = serviceOrderEntity.ID
			memberServiceEntity.DeadlineTime = memberServiceEntity.DeadlineTime.AddDate(0, serviceOrderEntity.PurchaseAmount, 0)
			dbutil.NewDbExecutor(tx.Table("member_service").Where("id = ? and status = 1", memberServiceEntity.ID).Update(map[string]interface{}{
				"order_id":      memberServiceEntity.OrderId,
				"deadline_time": memberServiceEntity.DeadlineTime,
			}))
		}
	}
	return Success
}
