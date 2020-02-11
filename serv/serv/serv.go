package servserv

import (
	"context"
	"fmt"
	memberrepo "focus/repo/member"
	servrepo "focus/repo/serv"
	servorderrepo "focus/repo/servorder"
	servpricerepo "focus/repo/servprice"
	ppayserv "focus/serv/ppay"
	"focus/tx"
	"focus/types"
	"focus/types/consts/orderstatus"
	pagetype "focus/types/page"
	ppaytype "focus/types/ppay"
	servtype "focus/types/serv"
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
	NotifyUrl  = "http://localhost:7001/api/v1/notify/ppay"
)

// 最新服务查询
func QueryLatest(ctx context.Context, reqParam *servtype.QueryLatestReq) *types.PageResponse {
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
	services, total := servrepo.QueryLatest(ctx, pageQuery)
	var results []*servtype.QueryLatestRes
	for _, service := range services {
		result := &servtype.QueryLatestRes{ServiceId: service.ID, ServiceType: service.ServiceType,
			ChineseName: service.ChineseName, ServiceDesc: service.ServiceDesc, PublishTime: timutil.DefFormat(service.PublishTime)}
		results = append(results, result)
	}
	return types.NewPageResponse(total, results)
}

// 服务详情查询
func GetById(ctx context.Context, serviceId int) *servtype.GetByIdRes {
	if serviceId <= 0 {
		types.InvalidParamPanic("serviceId is invalided!")
	}
	service := servrepo.GetById(ctx, serviceId)
	if service.ID == 0 {
		types.NotFoundPanic("service not exists!")
	}
	return &servtype.GetByIdRes{ServiceId: service.ID, ServiceType: service.ServiceType,
		ChineseName: service.ChineseName, ServiceDesc: service.ServiceDesc, PublishTime: timutil.DefFormat(service.PublishTime)}
}

// 服务套餐查询
func QueryPrice(ctx context.Context, serviceId int) []*servtype.QueryPriceRes {
	if serviceId <= 0 {
		types.InvalidParamPanic("serviceId is invalided!")
	}
	service := GetById(ctx, serviceId)
	if service == nil {
		types.NotFoundPanic("service not found!")
	}
	prices := servpricerepo.QueryByServiceId(ctx, serviceId)
	var results []*servtype.QueryPriceRes
	for _, price := range prices {
		result := &servtype.QueryPriceRes{ID: price.ID, PriceName: price.Price,
			ServiceId: price.ServiceId, Price: price.Price, ServiceAmount: price.ServiceAmount}
		results = append(results, result)
	}
	return results
}

// 服务价格计算器
func CalculatePrice(ctx context.Context, reqParam *servtype.CalculatePriceReq) *servtype.CalculatePriceRes {
	if reqParam.PriceId <= 0 {
		types.InvalidParamPanic("priceId param is invalid!")
	}
	if reqParam.Amount <= 0 {
		types.InvalidParamPanic("amount param is invalid!")
	}
	priceEntity := servpricerepo.GetById(ctx, reqParam.PriceId)
	if priceEntity.ID == 0 {
		types.NotFoundPanic("price not exists!")
	}
	price, err := decimal.NewFromString(priceEntity.Price)
	if err != nil {
		types.ErrPanic(types.DataDirty, fmt.Sprintf("invalid price=%s", priceEntity.Price))
	}
	decimal.DivisionPrecision = 2
	price = price.Mul(decimal.NewFromInt(reqParam.Amount))
	return &servtype.CalculatePriceRes{Price: price.StringFixedBank(2)}
}

// 生成服务订单
func CreateOrderTx(ctx context.Context) tx.TFunRes {
	reqParam := ctx.Value("reqParam").(*servtype.CreateOrderRequest)

	// 下单校验
	orderEntity, priceEntity, serviceEntity := createOrderValidation(ctx, reqParam)
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

	// 价格计算
	price = price.Mul(decimal.NewFromInt(int64(reqParam.PurchaseAmount))).Round(2)
	orderEntity.OrderAmount = price.StringFixedBank(2)
	orderEntity.PayAmount = price.StringFixedBank(2)
	orderEntity.CouponNo = reqParam.CouponNo
	orderEntity.PayChannel = reqParam.PayChannel

	// 生成服务订单
	servorderrepo.Create(ctx, orderEntity)

	// 生成支付订单
	createPayOrderReq := &ppaytype.CreateOrderReq{
		OutTradeNo:     orderEntity.OrderNo,
		OrderAmount:    orderEntity.OrderAmount,
		PayAmount:      orderEntity.PayAmount,
		PayChannel:     orderEntity.PayChannel,
		PayeeAccountId: types.PayChannels()[orderEntity.PayChannel].AccountId,
		PayReason:      fmt.Sprintf("购买服务:%v,套餐:%v,数量：%v", serviceEntity.ChineseName, priceEntity.PriceName, reqParam.PurchaseAmount),
		NotifyUrl:      NotifyUrl,
	}
	createPayOrderRes := ppayserv.CreateOrder(ctx, createPayOrderReq)

	// 更新服务订单的支付订单
	servorderrepo.Submit(ctx, orderEntity.ID, createPayOrderRes.PayOrderNo)

	// 支付收银台请求参数
	cashierParams := &servtype.CashierReq{
		PayOrderNo:     createPayOrderRes.PayOrderNo,
		ServiceOrderNo: orderEntity.OrderNo,
		OrderAmount:    createPayOrderReq.OrderAmount,
		PayAmount:      createPayOrderReq.PayAmount,
		PayChannel:     createPayOrderReq.PayChannel,
		PayReason:      createPayOrderReq.PayReason,
	}
	res := &servtype.CreateOrderRes{
		CashierUrl:    CashierUrl,
		CashierParams: cashierParams,
	}
	return res
}

func createOrderValidation(ctx context.Context, reqParam *servtype.CreateOrderRequest) (*servtype.OrderEntity, *servtype.PriceEntity, *servtype.ServiceEntity) {
	if reqParam.OrderNo == "" {
		types.InvalidParamPanic("orderNo can't be empty!")
	}
	if reqParam.MemberId <= 0 {
		types.InvalidParamPanic("memberId is invalid!")
	}
	if reqParam.ServicePriceId <= 0 {
		types.InvalidParamPanic("servicePriceId is invalid!")
	}
	if reqParam.PurchaseAmount < 1 {
		types.InvalidParamPanic("purchaseAmount is invalid!!")
	}
	if strutil.IsBlank(reqParam.PayChannel) {
		types.InvalidParamPanic("payChannel can't be empty!")
	}
	if types.PayChannels()[reqParam.PayChannel] == nil {
		types.ErrPanic(types.PayChannelNotSupport, "payChannel not support!")
	}
	memberEntity := memberrepo.GetById(ctx, reqParam.MemberId)
	if memberEntity.ID == 0 {
		types.NotFoundPanic(fmt.Sprintf("memberId %v not exists!", reqParam.MemberId))
	}
	orderEntity := servorderrepo.GetByOrderNo(ctx, reqParam.OrderNo)
	if orderEntity.ID != 0 {
		types.RepeatRequestPanic(fmt.Sprintf("order orderNo=%s already exists!", reqParam.OrderNo))
	}
	priceEntity := servpricerepo.GetById(ctx, reqParam.ServicePriceId)
	if priceEntity.ID == 0 {
		types.NotFoundPanic(fmt.Sprintf("price priceId=%d not exists!", reqParam.ServicePriceId))
	}
	serviceEntity := servrepo.GetById(ctx, priceEntity.ServiceId)
	if serviceEntity.ID == 0 {
		types.NotFoundPanic(fmt.Sprintf("service %d not found!", priceEntity.ServiceId))
	}
	return orderEntity, priceEntity, serviceEntity
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
	var serviceOrderEntity servtype.OrderEntity
	dbutil.NewDbExecutor(tx.Table("service_order").Where("out_order_no = ? and order_status = 'P' and status = 1", payResult.PayOrderNo).Find(&serviceOrderEntity))
	if serviceOrderEntity.ID == 0 {
		logrus.Warnf("serviceOrder not exists! outTradeNo:%s", payResult.PayOrderNo)
		return Success
	}
	updateResult := dbutil.NewDbExecutor(tx.Table("service_order").Where("out_order_no = ? and order_status = 'P' and status = 1", payResult.PayOrderNo).Updates(servtype.OrderEntity{OrderStatus: payResult.PayStatus, FinishedTime: time.Now()})).RowsAffected()
	if updateResult != 1 {
		logrus.Warn("serviceOrder has been processed! outTradeNo:%s", payResult.PayOrderNo)
		return Success
	}
	if payResult.PayStatus == orderstatusconst.S {
		var memberServiceEntity servtype.MemberServiceEntity
		dbutil.NewDbExecutor(tx.Table("member_service").Where("member_id = ? and service_price_id = ? and member_service_status = 1 and status = 1", serviceOrderEntity.MemberId, serviceOrderEntity.ServicePriceId).Find(&memberServiceEntity))
		if memberServiceEntity.ID == 0 {
			memberServiceEntity = servtype.MemberServiceEntity{
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
