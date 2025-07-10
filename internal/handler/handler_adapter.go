package handler

import (
	"QuickBrick/internal/service"
	"github.com/gin-gonic/gin"
)

type FuncWithHistoryService func(*gin.Context, *service.PipelineHistoryService)

func Adapt(fn FuncWithHistoryService, svc *service.PipelineHistoryService) gin.HandlerFunc {
	return func(c *gin.Context) {
		fn(c, svc)
	}
}
