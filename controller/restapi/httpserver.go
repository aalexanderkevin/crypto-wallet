package restapi

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/aalexanderkevin/crypto-wallet/config"
	"github.com/aalexanderkevin/crypto-wallet/container"
	"github.com/aalexanderkevin/crypto-wallet/controller/middleware"
	"github.com/aalexanderkevin/crypto-wallet/controller/restapi/handler"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type HttpServer interface {
	Start() error
	GetHandler() (http.Handler, error)
}

type httpServer struct {
	config     config.Config
	engine     *gin.Engine
	controller controllers
}

type controllers struct {
	webhook handler.Webhook
}

func NewHttpServer(container *container.Container) *httpServer {
	gin.SetMode(gin.ReleaseMode)
	if strings.ToLower(container.Config().LogLevel) == gin.DebugMode {
		gin.SetMode(gin.DebugMode)
	}

	engine := newGinEngine()

	controllers := controllers{
		*handler.NewWebhook(container),
	}
	requestHandler := &httpServer{container.Config(), engine, controllers}
	requestHandler.setupRouting()

	return requestHandler
}

func (h *httpServer) Start() error {
	return h.engine.Run(fmt.Sprintf("%s:%s", h.config.Service.Host, h.config.Service.Port))
}

func (h *httpServer) GetHandler() (http.Handler, error) {
	return h.engine, nil
}

func newGinEngine() *gin.Engine {
	r := gin.New()

	r.Use(middleware.LogrusLogger(logrus.StandardLogger()), gin.Recovery())

	return r
}
