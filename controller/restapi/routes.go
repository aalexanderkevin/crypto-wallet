package restapi

// setupRouting contains REST path and handler configuration
// @title webhook API
// @version 1.0
// @description Jobadder service REST API specification
// @host localhost:9004
// @BasePath /v1/jobadder
// @securityDefinitions.apikey BasicAuth
// @in header
// @name Authorization
func (h *httpServer) setupRouting() {
	router := h.engine
	v1 := router.Group(h.config.Service.Path.V1)

	// public API
	btcAPI := v1.Group(h.config.Service.Path.Btc)
	// Use(middleware.Auth(h.config.JobAdder.WebhookUsername, h.config.JobAdder.WebhookPassword))
	{
		btcAPI.POST("/webhook/transaction", h.controller.webhook.Transaction)
		// btcAPI.POST("/project/update/:tenant-id", h.controller.webhook.ProjectUpdate)
	}
}
