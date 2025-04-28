package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const metricsPath = "/metrics"

func NewMetricsRouter(router *gin.Engine) {
	router.GET(metricsPath, gin.WrapH(promhttp.Handler()))
}
