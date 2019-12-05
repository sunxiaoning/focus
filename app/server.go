package app

import (
	"focus/cfg"
	"focus/controller"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	"strconv"
)

type Server struct {
	Router *mux.Router
}

func InitServer() {
	router := controller.InitRouter()
	httpServer := &http.Server{
		Addr:    net.JoinHostPort("127.0.0.1", strconv.Itoa(cfg.FocusCtx.Cfg.Server.ListenPort)),
		Handler: router,
	}
	cfg.FocusCtx.HttpServer = httpServer
}

func StartServer() {
	if err := cfg.FocusCtx.HttpServer.ListenAndServe(); err != nil {
		logrus.Fatal("Fail to start http server: ", err)
	}
}
