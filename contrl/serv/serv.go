package servcontrl

import (
	"context"
	"encoding/json"
	servserv "focus/serv/serv"
	"focus/tx"
	"focus/types"
	servicetype "focus/types/service"
	"net/http"
	"strconv"
	"strings"
)

const (
	ServiceModule = types.ApiV1 + "/service"
)

func url(path string) string {
	return ServiceModule + path
}

var QueryLatest = types.NewController(url("/queryLatest"), http.MethodPost, queryLatest)

func queryLatest(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
	reqParam := servicetype.NewQueryLatestReq()
	if err := json.NewDecoder(req.Body).Decode(reqParam); err != nil {
		return types.NewErr(types.InvalidParamError, "invalid json!")
	}
	ctx = context.WithValue(ctx, "reqParam", reqParam)
	res, err := servserv.QueryLatest(ctx)
	if err != nil {
		return err
	}
	return types.NewRestRestResponse(rw, res)
}

var GetById = types.NewController(url("/getById"), http.MethodGet, getById)

func getById(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
	serviceIdStr := req.URL.Query().Get("serviceId")
	if strings.TrimSpace(serviceIdStr) == "" {
		return types.NewErr(types.InvalidParamError, "serviceId can't be empty!")
	}
	serviceId, err := strconv.Atoi(serviceIdStr)
	if err != nil {
		return types.NewErr(types.InvalidParamError, "serviceId is invalid!")
	}
	ctx = context.WithValue(ctx, "serviceId", serviceId)
	res, err := servserv.GetServiceById(ctx)
	if err != nil {
		return err
	}
	return types.NewRestRestResponse(rw, res)
}

var QueryPrice = types.NewController(url("/queryPrice"), http.MethodGet, queryPrice)

func queryPrice(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
	serviceIdStr := req.URL.Query().Get("serviceId")
	if strings.TrimSpace(serviceIdStr) == "" {
		return types.NewErr(types.InvalidParamError, "serviceId can't be empty!")
	}
	serviceId, err := strconv.Atoi(serviceIdStr)
	if err != nil {
		return types.NewErr(types.InvalidParamError, "serviceId is invalid!")
	}
	ctx = context.WithValue(ctx, "serviceId", serviceId)
	res, err := servserv.QueryPrice(ctx)
	if err != nil {
		return err
	}
	return types.NewRestRestResponse(rw, res)
}

var CalculatePrice = types.NewController(url("/calculatePrice"), http.MethodPost, calculatePrice)

func calculatePrice(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
	reqParam := servicetype.NewCalculatePriceReq()
	if err := json.NewDecoder(req.Body).Decode(reqParam); err != nil {
		return types.NewErr(types.InvalidParamError, "invalid json!")
	}
	ctx = context.WithValue(ctx, "reqParam", reqParam)
	res, err := servserv.CalculatePrice(ctx)
	if err != nil {
		return err
	}
	return types.NewRestRestResponse(rw, res)
}

var CreateOrder = types.NewController(url("/createOrder"), http.MethodPost, createOrder)

func createOrder(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
	var reqParam = servicetype.NewCreateOrderReq()
	if err := json.NewDecoder(req.Body).Decode(reqParam); err != nil {
		return types.InvalidParamErr("invalid json format!")
	}
	ctx = context.WithValue(ctx, "reqParam", reqParam)
	data, err := tx.NewTxManager().RunTx(ctx, servserv.CreateOrderTx)
	if err != nil {
		return err
	}
	return types.NewRestRestResponse(rw, data)
}

var Cashier = types.NewController(url("/cashier"), http.MethodPost, cashier)

func cashier(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
	return nil
}
