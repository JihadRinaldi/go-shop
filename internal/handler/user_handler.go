package handler

import (
	"github.com/JihadRinaldi/go-shop/internal/dto"
	"github.com/JihadRinaldi/go-shop/internal/services"
	"github.com/JihadRinaldi/go-shop/internal/utils"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	profile, err := h.userService.GetProfile(userID)
	if err != nil {
		utils.NotFoundResponse(c, "User not found")
		return
	}

	utils.SuccessResponse(c, "Profile fetched", profile)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}

	profile, err := h.userService.UpdateProfile(userID, &req)
	if err != nil {
		utils.BadRequestResponse(c, "Failed to update profile", err)
		return
	}

	utils.SuccessResponse(c, "Profile updated", profile)
}
