package v1

import (
	"net/http"

	"github.com/S1riyS/wildberries-techschool/L0/server/internal/service"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	service *service.OrderService
}

func NewOrderHandler(service *service.OrderService) *OrderHandler {
	return &OrderHandler{
		service: service,
	}
}

func (h *OrderHandler) GetOne(ctx *gin.Context) {
	id := ctx.Param("id")

	order, err := h.service.GetOne(ctx.Request.Context(), id)
	if err != nil {
		ctx.Error(err) //nolint:errcheck // golangcli issue
		return
	}

	ctx.JSON(http.StatusOK, order)
}
