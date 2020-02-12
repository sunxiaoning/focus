package servcontrl

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	ppayserv "focus/serv/ppay"
	servserv "focus/serv/serv"
	"focus/tx"
	"focus/types"
	ppaytype "focus/types/ppay"
	servtype "focus/types/serv"
	strutil "focus/util/strs"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strconv"
	"strings"
)

const (
	ServiceModule = types.ApiV1 + "/service"
	BufferSize    = 256 * 1024
)

func url(path string) string {
	return ServiceModule + path
}

var QueryLatest = types.NewController(url("/queryLatest"), http.MethodPost, queryLatest)

func queryLatest(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
	reqParam := servtype.NewQueryLatestReq()
	if err := json.NewDecoder(req.Body).Decode(reqParam); err != nil {
		types.InvalidParamPanic("invalid json!")
	}
	types.NewRestRestResponse(rw, servserv.QueryLatest(ctx, reqParam))
}

var GetById = types.NewController(url("/getById"), http.MethodGet, getById)

func getById(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
	serviceIdStr := req.URL.Query().Get("serviceId")
	if strings.TrimSpace(serviceIdStr) == "" {
		types.InvalidParamPanic("serviceId can't be empty!")
	}
	serviceId, err := strconv.Atoi(serviceIdStr)
	if err != nil {
		types.InvalidParamPanic("serviceId is invalid!")
	}
	types.NewRestRestResponse(rw, servserv.GetById(ctx, serviceId))
}

var QueryPrice = types.NewController(url("/queryPrice"), http.MethodGet, queryPrice)

func queryPrice(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
	serviceIdStr := req.URL.Query().Get("serviceId")
	if strings.TrimSpace(serviceIdStr) == "" {
		types.InvalidParamPanic("serviceId can't be empty!")
	}
	serviceId, err := strconv.Atoi(serviceIdStr)
	if err != nil {
		types.InvalidParamPanic("serviceId is invalid!")
	}
	types.NewRestRestResponse(rw, servserv.QueryPrice(ctx, serviceId))
}

var CalculatePrice = types.NewController(url("/calculatePrice"), http.MethodPost, calculatePrice)

func calculatePrice(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
	reqParam := servtype.NewCalculatePriceReq()
	if err := json.NewDecoder(req.Body).Decode(reqParam); err != nil {
		types.InvalidParamPanic("invalid json!")
	}
	types.NewRestRestResponse(rw, servserv.CalculatePrice(ctx, reqParam))
}

var CreateOrder = types.NewController(url("/createOrder"), http.MethodPost, createOrder)

func createOrder(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
	var reqParam = servtype.NewCreateOrderReq()
	if err := json.NewDecoder(req.Body).Decode(reqParam); err != nil {
		types.InvalidParamPanic("invalid json format!")
	}
	ctx = context.WithValue(ctx, "reqParam", reqParam)
	types.NewRestRestResponse(rw, tx.NewTxManager().RunTx(ctx, servserv.CreateOrderTx))
}

var Cashier = types.NewController(url("/cashier"), http.MethodPost, cashier)

func cashier(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
	var reqParam = &ppaytype.CashierReq{}
	if err := json.NewDecoder(req.Body).Decode(reqParam); err != nil {
		types.InvalidParamPanic("invalid json format!")
	}
	types.NewRestRestResponse(rw, ppayserv.Cashier(ctx, reqParam))
}

var GetReceiptCode = types.NewController(url("/getReceiptCode/{qrCodeUrl:\\S+}"), http.MethodGet, getReceiptCode)

func getReceiptCode(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	io.Copy(rw, bytes.NewBuffer(ppayserv.GetReceiptCode(ctx, vars["qrCodeUrl"])))
}

var UploadReceiptCode = types.NewController(url("/uploadReceiptCode"), http.MethodPost, uploadReceiptCode)

func uploadReceiptCode(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
	if err := req.ParseMultipartForm(BufferSize); err != nil {
		types.SystemPanic(fmt.Sprintf("parse multipart/form err! reason: %v", err))
	}
	operator := req.FormValue("operator")
	if strutil.IsBlank(operator) {
		types.InvalidParamPanic("operator can't be empty!")
	}
	payeeAccountIdStr := req.FormValue("payeeAccountId")
	if strutil.IsBlank(payeeAccountIdStr) {
		types.InvalidParamPanic("payeeAccountId can't be empty!")
	}
	payeeAccountId, err := strconv.Atoi(payeeAccountIdStr)
	if err != nil {
		types.InvalidParamPanic("payeeAccountId is invalid!")
	}
	payeeAmount := req.FormValue("payeeAmount")
	if strutil.IsBlank(payeeAmount) {
		types.InvalidParamPanic("payeeAmount can't be empty!")
	}
	if !strutil.IsValidMoney(payeeAmount) {
		types.InvalidParamPanic("payeeAmount is invalid!")
	}
	file, header, err := req.FormFile("receiptCode")
	if err != nil {
		types.SystemPanic(fmt.Sprintf("read receiptcode err! reason: %v", err))
	}
	reqParam := &ppaytype.UploadReceiptCodeReq{
		Operator:       operator,
		PayeeAccountId: payeeAccountId,
		PayeeAmount:    payeeAmount,
		File:           file,
		FileHeader:     header,
	}
	types.NewRestRestResponse(rw, ppayserv.UploadReceiptCode(ctx, reqParam))
}

// 支付系统回调通知支付结果
var PayResultNotify = types.NewController(url("/payResultNotify"), http.MethodPost, payResultNotify)

func payResultNotify(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
	payResult := &ppaytype.BizPayResultReq{}
	if err := json.NewDecoder(req.Body).Decode(payResult); err != nil {
		types.InvalidParamPanic("json invalid!")
	}
	logrus.Infof("payResult: %v", payResult)
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(servserv.PayResultNotify(ctx, payResult)))
}
