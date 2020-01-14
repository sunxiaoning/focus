package contrl

import (
	"context"
	"encoding/json"
	"fmt"
	"focus/filter"
	"focus/types"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"reflect"
)

func InitRouter(controllers []*types.Controller) *mux.Router {
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

		// 执行Controller逻辑
		controller.Handle(filter.Process(ctx, rw, req), rw, req)
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
