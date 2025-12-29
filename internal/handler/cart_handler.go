package handler

import (
	"strconv"

	"github.com/JihadRinaldi/go-shop/internal/dto"
	"github.com/JihadRinaldi/go-shop/internal/services"
	"github.com/JihadRinaldi/go-shop/internal/utils"
	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	cartService *services.CartService
}

func NewCartHandler(cartService *services.CartService) *CartHandler {
	return &CartHandler{
		cartService: cartService,
	}
}

func (h *CartHandler) GetCart(c *gin.Context) {
	userID := c.GetUint("user_id")

	cart, err := h.cartService.GetCart(userID)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to fetch cart", err)
		return
	}

	utils.SuccessResponse(c, "Cart fetched", cart)
}

func (h *CartHandler) AddToCart(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}

	cart, err := h.cartService.AddToCart(userID, req)
	if err != nil {
		utils.BadRequestResponse(c, err.Error(), err)
		return
	}

	utils.SuccessResponse(c, "Item added to cart", cart)
}

func (h *CartHandler) UpdateCartItem(c *gin.Context) {
	userID := c.GetUint("user_id")

	itemID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid item ID", err)
		return
	}

	var req dto.UpdateCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}

	cart, err := h.cartService.UpdateCartItem(userID, uint(itemID), req)
	if err != nil {
		utils.BadRequestResponse(c, err.Error(), err)
		return
	}

	utils.SuccessResponse(c, "Cart item updated", cart)
}

func (h *CartHandler) RemoveCartItem(c *gin.Context) {
	userID := c.GetUint("user_id")

	itemID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid item ID", err)
		return
	}

	if err := h.cartService.RemoveCartItem(userID, uint(itemID)); err != nil {
		utils.InternalServerErrorResponse(c, "Failed to remove cart item", err)
		return
	}

	utils.SuccessResponse(c, "Cart item removed", nil)
}
