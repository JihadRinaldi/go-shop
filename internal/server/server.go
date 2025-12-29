package server

import (
	"net/http"

	"github.com/JihadRinaldi/go-shop/internal/config"
	"github.com/JihadRinaldi/go-shop/internal/handler"
	"github.com/JihadRinaldi/go-shop/internal/interfaces"
	"github.com/JihadRinaldi/go-shop/internal/providers"
	"github.com/JihadRinaldi/go-shop/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type Server struct {
	config         *config.Config
	db             *gorm.DB
	logger         zerolog.Logger
	authHandler    *handler.AuthHandler
	userHandler    *handler.UserHandler
	productHandler *handler.ProductHandler
}

func New(cfg *config.Config, db *gorm.DB, logger *zerolog.Logger) *Server {
	var uploadProvider interfaces.UploadProvider
	if cfg.Upload.UploadProvider == "s3" {
		uploadProvider = providers.NewS3Provider(cfg)
	} else {
		uploadProvider = providers.NewLocalUploadProvider(cfg.Upload.Path)
	}

	authService := services.NewAuthService(db, cfg)
	userService := services.NewUserService(db, cfg)
	productService := services.NewProductService(db, cfg)
	uploadService := services.NewUploadService(uploadProvider)

	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	productHandler := handler.NewProductHandler(productService, uploadService)

	return &Server{
		config:         cfg,
		db:             db,
		logger:         *logger,
		authHandler:    authHandler,
		userHandler:    userHandler,
		productHandler: productHandler,
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

			categories := protected.Group("/categories")
			{
				categories.POST("/", s.adminMiddleware(), s.productHandler.CreateCategory)
				categories.PUT("/:id", s.adminMiddleware(), s.productHandler.UpdateCategory)
				categories.DELETE("/:id", s.adminMiddleware(), s.productHandler.DeleteCategory)
			}

			products := protected.Group("/products")
			{
				productRoutes := products
				productRoutes.POST("/", s.adminMiddleware(), s.productHandler.CreateProduct)
				productRoutes.PUT("/:id", s.adminMiddleware(), s.productHandler.UpdateProduct)
				productRoutes.DELETE("/:id", s.adminMiddleware(), s.productHandler.DeleteProduct)
				productRoutes.POST("/:id/images", s.adminMiddleware(), s.productHandler.UploadProductImage)
			}
		}

		api.GET("/categories", s.productHandler.GetCategories)
		api.GET("/products", s.productHandler.GetProducts)
		api.GET("/products/:id", s.productHandler.GetProduct)
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
