package servserv

import (
	"context"
	"fmt"
	"focus/cfg"
	ppayserv "focus/serv/ppay"
	"focus/tx"
	"focus/types"
	"focus/types/consts/orderstatus"
	usertype "focus/types/member"
	ppaytype "focus/types/ppay"
	servicetype "focus/types/service"
	dbutil "focus/util/db"
	strutil "focus/util/strs"
	timutil "focus/util/tim"
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
	"time"
)

const (
	CashierUrl = "/api/v1/service/cashier"
)

func QueryLatest(ctx context.Context) *types.PageResponse {
	reqParam := ctx.Value("reqParam").(*servicetype.QueryLatestReq)
	if reqParam.PageIndex < 1 {
		types.InvalidParamPanic("pageIndex is invalid!")
	}
	if reqParam.PageSize < 1 || reqParam.PageSize > 1000 {
		types.InvalidParamPanic("pageSize is invalid!")
	}
	query := map[string]interface{}{}
	if reqParam.ChineseName != "" {
		query["chinese_name"] = reqParam.ChineseName
	}
	if reqParam.ServiceType != 0 {
		query["service_type"] = reqParam.ServiceType
	}
	var services []*servicetype.ServiceEntity
	var total int
	dbutil.NewDbExecutor(cfg.FocusCtx.DB.Table("service").Where(query).Where("service_status = 'FWZ' and status = 1")).PageQuery(reqParam.PageIndex, reqParam.PageSize, &total, &services)
	var results []*servicetype.QueryLatestRes
	for _, service := range services {
		result := &servicetype.QueryLatestRes{ServiceId: service.ID, ServiceType: service.ServiceType, ChineseName: service.ChineseName, ServiceDesc: service.ServiceDesc, PublishTime: timutil.DefFormat(service.PublishTime)}
		results = append(results, result)
	}
	return types.NewPageResponse(total, results)
}

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
	price = price.Mul(decimal.NewFromInt(reqParam.PurchaseAmount)).Round(2)
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
