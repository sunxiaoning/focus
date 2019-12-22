package contrl

import (
	"context"
	"encoding/json"
	"fmt"
	servcontrl "focus/contrl/serv"
	"focus/contrl/user"
	"focus/filter"
	"focus/types"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"reflect"
)

var controllers = []*types.Controller{
	Hi, Hello, Err, usercontrl.Login, servcontrl.QueryLatest,
	servcontrl.GetById, servcontrl.QueryPrice,
	servcontrl.CalculatePrice, servcontrl.CreateOrder,
}

func InitRouter() *mux.Router {
	router := mux.NewRouter()
	for _, controller := range controllers {
		router.Path(controller.Path).Methods(controller.Method).Handler(handle(controller))
	}
	return router
}

func handle(controller *types.Controller) http.HandlerFunc {
	pctx := context.Background()
	return func(rw http.ResponseWriter, req *http.Request) {
		ctx, cancel := context.WithCancel(pctx)
		defer func() {
			cancel()
			if r := recover(); r != nil {
				err, ok := r.(*types.FocusError)
				if ok {
					handleErrResponse(rw, err)
				} else {
					handleErrResponse(rw, types.SystemErr(fmt.Sprintf("%v", r)))
				}
			}
		}()

		// 过滤器执行
		if err := filter.Process(ctx, rw, req); err != nil {
			handleErrResponse(rw, err)
			return
		}

		// 执行Controller逻辑
		if err := controller.Handle(ctx, rw, req); err != nil {
			handleErrResponse(rw, err)
		}
		return
	}

}

func handleErrResponse(rw http.ResponseWriter, err error) {
	logrus.Error("handle error: ", err)
	var httpstatus int
	var errcode int
	var errmsg string
	rw.Header().Set("Content-Type", "application/json")
	if (reflect.TypeOf(err) != reflect.TypeOf(&types.FocusError{})) {
		logrus.Errorf("unexpected error! err: %s", err)
		httpstatus = http.StatusInternalServerError
		errcode = types.SystemError
		errmsg = fmt.Sprintf("unexpected error! err: %s", err)
	} else {
		appError := err.(*types.FocusError)
		if appError.Code == types.SystemError {
			httpstatus = http.StatusInternalServerError
		} else if appError.Code == types.InvalidParamError {
			httpstatus = http.StatusBadRequest
		} else {
			httpstatus = http.StatusOK
		}
		errcode = appError.Code
		errmsg = appError.Message
	}

	rw.WriteHeader(httpstatus)
	encoder := json.NewEncoder(rw)
	encoder.SetEscapeHTML(false)
	res := make(map[string]interface{})
	res[types.ErrCode] = errcode
	res[types.ErrMsg] = errmsg
	encoder.Encode(res)
}
