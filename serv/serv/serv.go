package servserv

import (
	"context"
	"fmt"
	"focus/cfg"
	ppayservice "focus/serv/ppay"
	"focus/tx"
	"focus/types"
	"focus/types/consts/orderstatus"
	ppaytype "focus/types/ppay"
	servicetype "focus/types/service"
	"focus/util/page"
	timutil "focus/util/tim"
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
	"time"
)

var (
	CashierUrl = "/api/v1/service/cashier"
)

func QueryLatest(ctx context.Context) (*types.PageResponse, error) {
	reqParam := ctx.Value("reqParam").(*servicetype.QueryLatestReq)
	if reqParam.PageIndex < 1 {
		return nil, types.NewErr(types.InvalidParamError, "pageIndex is invalid!")
	}
	if reqParam.PageSize < 1 || reqParam.PageSize > 1000 {
		return nil, types.NewErr(types.InvalidParamError, "pageSize is invalid!")
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
	dbQuery := cfg.FocusCtx.DB.Table("service").Where(query).Where("service_status = 'FWZ' and status = 1")
	pageutil.PageQuery(dbQuery, reqParam.PageIndex, reqParam.PageSize, &total, &services)
	var results []*servicetype.QueryLatestRes
	for _, service := range services {
		result := &servicetype.QueryLatestRes{ServiceId: service.ID, ServiceType: service.ServiceType, ChineseName: service.ChineseName, ServiceDesc: service.ServiceDesc, PublishTime: timutil.DefFormat(service.PublishTime)}
		results = append(results, result)
	}
	return types.NewPageResponse(total, results), nil
}

func GetServiceById(ctx context.Context) (*servicetype.GetByIdRes, error) {
	serviceId := ctx.Value("serviceId").(int)
	if serviceId <= 0 {
		return nil, types.NewErr(types.InvalidParamError, "serviceId is invalided!")
	}
	service := &servicetype.ServiceEntity{}
	cfg.FocusCtx.DB.Table("service").Where("id = ? and status = 1", serviceId).First(service)
	if service.ID == 0 {
		return nil, types.NewErr(types.NotFound, "service not exists!")
	}
	return &servicetype.GetByIdRes{ServiceId: service.ID, ServiceType: service.ServiceType, ChineseName: service.ChineseName, ServiceDesc: service.ServiceDesc, PublishTime: timutil.DefFormat(service.PublishTime)}, nil
}

func QueryPrice(ctx context.Context) ([]*servicetype.QueryPriceRes, error) {
	serviceId := ctx.Value("serviceId").(int)
	if serviceId <= 0 {
		return nil, types.NewErr(types.InvalidParamError, "serviceId is invalided!")
	}
	service, err := GetServiceById(ctx)
	if err != nil {
		return nil, err
	}
	if service == nil {
		return nil, types.NewErr(types.NotFound, "service not found!")
	}
	var prices []*servicetype.PriceEntity
	cfg.FocusCtx.DB.Table("service_price").Where("service_id = ? and status = 1", serviceId).Find(&prices)
	var results []*servicetype.QueryPriceRes
	for _, price := range prices {
		result := &servicetype.QueryPriceRes{ID: price.ID, PriceName: price.Price, ServiceId: price.ServiceId, Price: price.Price, ServiceAmount: price.ServiceAmount}
		results = append(results, result)
	}
	return results, nil
}

func CalculatePrice(ctx context.Context) (*servicetype.CalculatePriceRes, error) {
	reqParam := ctx.Value("reqParam").(*servicetype.CalculatePriceReq)
	if reqParam.PriceId <= 0 {
		return nil, types.NewErr(types.InvalidParamError, "priceId param is invalid!")
	}
	if reqParam.Amount <= 0 {
		return nil, types.NewErr(types.InvalidParamError, "amount param is invalid!")
	}
	var priceEntity servicetype.PriceEntity
	cfg.FocusCtx.DB.Table("service_price").Where("id = ? and status = 1", reqParam.PriceId).First(&priceEntity)
	if priceEntity.ID == 0 {
		return nil, types.NewErr(types.NotFound, "price not exists!")
	}
	price, err := decimal.NewFromString(priceEntity.Price)
	if err != nil {
		return nil, types.NewErr(types.DataDirty, fmt.Sprintf("invalid price=%s", priceEntity.Price))
	}
	decimal.DivisionPrecision = 2
	price = price.Mul(decimal.NewFromInt(reqParam.Amount))
	return &servicetype.CalculatePriceRes{Price: price.String()}, nil
}

func CreateOrderTx(ctx context.Context) (tx.TFunRes, error) {
	tx := ctx.Value("tx").(*gorm.DB)
	reqParam := ctx.Value("reqParam").(*servicetype.CreateOrderRequest)
	if reqParam.OrderNo == "" {
		return nil, types.InvalidParamErr("orderNo can't be empty!")
	}
	if reqParam.MemberId <= 0 {
		return nil, types.InvalidParamErr("memberId is invalid!")
	}
	if reqParam.ServicePriceId <= 0 {
		return nil, types.InvalidParamErr("servicePriceId is invalid!")
	}
	if reqParam.PurchaseAmount <= 1 {
		return nil, types.InvalidParamErr("purchaseAmount is invalid!!")
	}
	if reqParam.PayChannel == "" {
		return nil, types.InvalidParamErr("payChannel can't be empty!")
	}
	if types.PayChannels()[reqParam.PayChannel] == nil {
		return nil, types.NewErr(types.PayChannelNotSupport, "payChannel not support!")
	}
	var orderEntity servicetype.OrderEntity
	tx.Table("service_order").Where("order_no = ? and status = 1", reqParam.OrderNo).First(&orderEntity)
	if orderEntity.ID != 0 {
		return nil, types.RepeatRequestErr(fmt.Sprintf("order orderNo=%s already exists!", reqParam.OrderNo))
	}
	var priceEntity servicetype.PriceEntity
	tx.Table("service_price").Where("id = ? and status = 1", reqParam.ServicePriceId).First(&priceEntity)
	if priceEntity.ID == 0 {
		return nil, types.NotFoundErr(fmt.Sprintf("price priceId=%d not exists!", reqParam.ServicePriceId))
	}
	orderEntity.OrderNo = reqParam.OrderNo
	orderEntity.MemberId = reqParam.MemberId
	orderEntity.ServicePriceId = reqParam.ServicePriceId
	orderEntity.PurchaseAmount = reqParam.PurchaseAmount
	orderEntity.StartTime = time.Now()
	orderEntity.FinishedTime = timutil.ZERO
	orderEntity.OrderStatus = orderstatusconst.I
	price, err := decimal.NewFromString(priceEntity.Price)
	price.StringFixedBank(2)
	if err != nil {
		types.NewErr(types.DataDirty, fmt.Sprintf("invalid price=%s", priceEntity.Price))
	}
	price = price.Mul(decimal.NewFromInt(reqParam.PurchaseAmount)).Round(2)
	orderEntity.OrderAmount = price.StringFixedBank(2)
	orderEntity.PayAmount = price.StringFixedBank(2)
	orderEntity.CouponNo = reqParam.CouponNo
	orderEntity.PayChannel = reqParam.PayChannel
	err = tx.Table("service_order").Create(&orderEntity).Error
	if err != nil {
		return nil, types.DbErr(err)
	}
	createPayOrderReq := &ppaytype.CreateOrderReq{
		OutTradeNo:     orderEntity.OrderNo,
		OrderAmount:    orderEntity.OrderAmount,
		PayAmount:      orderEntity.PayAmount,
		PayChannel:     orderEntity.PayChannel,
		PayeeAccountId: types.PayChannels()[orderEntity.PayChannel].AccountId,
		PayReason:      fmt.Sprintf("购买%s", priceEntity.PriceName),
		NotifyUrl:      fmt.Sprintf("http://localhost:%d/api/v1/notify/ppay", cfg.FocusCtx.Cfg.Server.ListenPort),
	}
	ctx = context.WithValue(ctx, "reqParam", createPayOrderReq)
	createPayOrderRes, err := ppayservice.CreateOrder(ctx)
	if err != nil {
		return nil, err
	}
	cashierParams := &servicetype.CashierReq{
		OrderNo:     orderEntity.OrderNo,
		OutOrderNo:  createPayOrderRes.PayOrderNo,
		OrderAmount: createPayOrderReq.OrderAmount,
		RealAmount:  createPayOrderReq.PayAmount,
		PayChannel:  createPayOrderReq.PayChannel,
		PayReason:   createPayOrderReq.PayReason,
	}
	res := servicetype.CreateOrderRes{
		CashierUrl:    CashierUrl,
		CashierParams: cashierParams,
	}
	return res, nil
}
