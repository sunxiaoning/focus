package app

import (
	"focus/cfg"
	"focus/contrl"
	gtwctl "focus/contrl/gtw"
	ppayctl "focus/contrl/ppay"
	servcontrl "focus/contrl/serv"
	usercontrl "focus/contrl/user"
	"focus/filter"
	"focus/types"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	"strconv"
)

type Server struct {
	Router *mux.Router
}

var Apis = []*types.Controller{
	usercontrl.Login, servcontrl.QueryLatest,
	servcontrl.GetById, servcontrl.QueryPrice,
	servcontrl.CalculatePrice, servcontrl.CreateOrder, servcontrl.Cashier,
	servcontrl.GetReceiptCode, servcontrl.UploadReceiptCode,
	servcontrl.PayResultNotify,
	ppayctl.PPay,
}

var ApiFilters = []*types.Filter{
	filter.UserIdentityAuthor,
}

var Gtw = []*types.Controller{
	gtwctl.Gtw,
}

var GtwFilters = []*types.Filter{
	filter.UserIdentityAuthor, filter.SignCheck, filter.ServiceAuth, filter.VisiterLimiter,
}

func InitServer(serverPort int, controllers []*types.Controller, filters []*types.Filter) {
	filter.InitFilter(filters)
	router := contrl.InitRouter(controllers)
	if serverPort == 0 {
		serverPort = cfg.FocusCtx.Cfg.Server.ListenPort
	}
	logrus.Infof("serverPort:%d", serverPort)
	httpServer := &http.Server{
		Addr:    net.JoinHostPort("127.0.0.1", strconv.Itoa(serverPort)),
		Handler: router,
	}
	cfg.FocusCtx.HttpServer = httpServer
}

func StartServer() {
	if err := cfg.FocusCtx.HttpServer.ListenAndServe(); err != nil {
		logrus.Fatal("Fail to start http server: ", err)
	}
}
