package server

import (
	"strings"

	"github.com/JihadRinaldi/go-shop/internal/models"
	"github.com/JihadRinaldi/go-shop/internal/utils"
	"github.com/gin-gonic/gin"
)

const (
	AuthorizationHeader = "Authorization"
	AuthorizationBearer = "Bearer"
)

func (s *Server) authMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader(AuthorizationHeader)
		if authHeader == "" {
			utils.UnauthorizedResponse(ctx, "Authorization header required")
			ctx.Abort()
			return
		}

		tokenString := strings.Split(authHeader, " ")
		if len(tokenString) != 2 || tokenString[0] != AuthorizationBearer {
			utils.UnauthorizedResponse(ctx, "Authorization header required")
			ctx.Abort()
			return
		}

		claims, err := utils.ValidateToken(tokenString[1], s.config.JWT.SecretKey)
		if err != nil {
			utils.UnauthorizedResponse(ctx, "Authorization header required")
			ctx.Abort()
			return
		}

		ctx.Set("user_id", claims.UserID)
		ctx.Set("email", claims.Email)
		ctx.Set("role", claims.Role)

		ctx.Next()
	}
}

func (s *Server) adminMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role, exists := ctx.Get("role")
		if !exists || role != models.UserRoleAdmin {
			utils.ForbiddenResponse(ctx, "Admin access required")
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
