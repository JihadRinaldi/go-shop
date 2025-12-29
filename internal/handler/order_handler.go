package handler

import (
	"strconv"

	"github.com/JihadRinaldi/go-shop/internal/services"
	"github.com/JihadRinaldi/go-shop/internal/utils"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService *services.OrderService
}

func NewOrderHandler(orderService *services.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	userID := c.GetUint("user_id")
	orderResponse, err := h.orderService.CreateOrder(userID)
	if err != nil {
		utils.BadRequestResponse(c, "Failed to create order", err)
		return
	}

	utils.SuccessResponse(c, "Order created successfully", orderResponse)
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	userID := c.GetUint("user_id")
	orderID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid order ID", err)
		return
	}

	orderResponse, err := h.orderService.GetOrder(userID, uint(orderID))
	if err != nil {
		utils.NotFoundResponse(c, "Order not found")
		return
	}

	utils.SuccessResponse(c, "Order fetched successfully", orderResponse)
}

func (h *OrderHandler) GetOrders(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	orders, meta, err := h.orderService.GetOrders(userID, page, limit)

	if err != nil {
		utils.BadRequestResponse(c, "Failed to fetch orders", err)
		return
	}

	utils.PaginatedSuccessResponse(c, "Orders fetched successfully", orders, *meta)
}
