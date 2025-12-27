package server

import (
	"net/http"

	"github.com/JihadRinaldi/go-shop/internal/config"
	"github.com/JihadRinaldi/go-shop/internal/handler"
	"github.com/JihadRinaldi/go-shop/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type Server struct {
	config      *config.Config
	db          *gorm.DB
	logger      zerolog.Logger
	authHandler *handler.AuthHandler
	userHandler *handler.UserHandler
}

func New(cfg *config.Config, db *gorm.DB, logger *zerolog.Logger) *Server {
	authService := services.NewAuthService(db, cfg)
	userService := services.NewUserService(db, cfg)

	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)

	return &Server{
		config:      cfg,
		db:          db,
		logger:      *logger,
		authHandler: authHandler,
		userHandler: userHandler,
	}
}

func (s *Server) SetupRoutes() *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(s.corsMiddleware())

	router.GET("/healthz", s.healthCheck)

	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", s.authHandler.Register)
			auth.POST("/login", s.authHandler.Login)
			auth.POST("/refresh", s.authHandler.RefreshToken)
			auth.POST("/logout", s.authHandler.Logout)
		}

		protected := api.Group("/")
		protected.Use(s.authMiddleware())
		{
			user := protected.Group("/user")
			{
				user.GET("/profile", s.userHandler.GetProfile)
				user.PUT("/profile", s.userHandler.UpdateProfile)
			}
		}
	}

	return router
}

func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (s *Server) corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
