package v1

import (
	"log/slog"

	"github.com/S1riyS/wildberries-techschool/L0/server/internal/service"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	logger  *slog.Logger
	service service.OrderService
}

func (oh *OrderHandler) GetOne(ctx *gin.Context) {

}
