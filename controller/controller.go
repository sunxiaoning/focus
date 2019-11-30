package controller

import (
	"context"
	"encoding/json"
	"focus/types"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Handle func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error

type Controller struct {
	Path   string
	Method string
	Handle Handle
}

func NewController(path string, method string, handle Handle) *Controller {
	return &Controller{
		Path:   path,
		Method: method,
		Handle: handle,
	}
}

var controllers = []*Controller{
	Hi, Hello, Err,
}

func InitRouter() *mux.Router {
	router := mux.NewRouter()
	for _, controller := range controllers {
		router.Path(controller.Path).Methods(controller.Method).Handler(filter(controller))
	}
	return router
}

func filter(controller *Controller) http.HandlerFunc {
	pctx := context.Background()
	return func(rw http.ResponseWriter, req *http.Request) {
		ctx, cancel := context.WithCancel(pctx)
		defer cancel()
		if err := controller.Handle(ctx, rw, req); err != nil {
			logrus.Error("handle error: ", err)
			handleErrResponse(rw, err)
		}
		return
	}

}

func handleErrResponse(rw http.ResponseWriter, err error) {
	appError := err.(*types.FocusError)
	var code int
	if appError.Code == types.SystemError {
		code = http.StatusInternalServerError
	} else if appError.Code == types.InvalidParamError {
		code = http.StatusBadRequest
	} else {
		code = http.StatusOK
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(code)
	encoder := json.NewEncoder(rw)
	encoder.SetEscapeHTML(false)
	res := make(map[string]interface{})
	res[types.ErrCode] = appError.Code
	res[types.ErrMsg] = appError.Message
	encoder.Encode(res)
}
