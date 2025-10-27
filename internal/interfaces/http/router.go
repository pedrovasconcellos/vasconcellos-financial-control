package http

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	"github.com/vasconcellos/financial-control/internal/interfaces/http/handler"
	"github.com/vasconcellos/financial-control/internal/interfaces/http/middleware"
)

type RouterParams struct {
	AuthHandler        *handler.AuthHandler
	AccountHandler     *handler.AccountHandler
	CategoryHandler    *handler.CategoryHandler
	TransactionHandler *handler.TransactionHandler
	BudgetHandler      *handler.BudgetHandler
	GoalHandler        *handler.GoalHandler
	ReportHandler      *handler.ReportHandler
	HealthHandler      *handler.HealthHandler
	AuthMiddleware     *middleware.AuthMiddleware
	AllowedOrigins     []string
	Logger             *zap.Logger
	ForceHTTPS         bool
	Environment        string
}

func NewRouter(params RouterParams) *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Recovery())

	// Middleware de redirecionamento HTTPS (deve ser aplicado primeiro)
	if params.ForceHTTPS {
		engine.Use(middleware.NewHTTPSRedirectMiddleware(true, params.Environment).Handle())
	}

	if params.Logger != nil {
		engine.Use(middleware.NewRequestLoggerMiddleware(params.Logger).Handle())
	}
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     params.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Swagger documentation endpoint (open, no auth required)
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/swagger/doc.json")))

	api := engine.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			if params.HealthHandler != nil {
				v1.GET("/health", params.HealthHandler.Status)
			}
			v1.POST("/auth/login", params.AuthHandler.Login)

			protected := v1.Group("/")
			protected.Use(params.AuthMiddleware.Handle())

			protected.GET("/accounts", params.AccountHandler.List)
			protected.POST("/accounts", params.AccountHandler.Create)
			protected.PATCH("/accounts/:id", params.AccountHandler.Update)
			protected.DELETE("/accounts/:id", params.AccountHandler.Delete)

			protected.GET("/categories", params.CategoryHandler.List)
			protected.POST("/categories", params.CategoryHandler.Create)
			protected.DELETE("/categories/:id", params.CategoryHandler.Delete)

			protected.GET("/transactions", params.TransactionHandler.List)
			protected.POST("/transactions", params.TransactionHandler.Create)
			protected.PATCH("/transactions/:id", params.TransactionHandler.Update)
			protected.POST("/transactions/:id/receipt", params.TransactionHandler.AttachReceipt)

			protected.GET("/budgets", params.BudgetHandler.List)
			protected.POST("/budgets", params.BudgetHandler.Create)

			protected.GET("/goals", params.GoalHandler.List)
			protected.POST("/goals", params.GoalHandler.Create)
			protected.POST("/goals/:id/progress", params.GoalHandler.UpdateProgress)

			protected.GET("/reports/summary", params.ReportHandler.Summary)
		}
	}

	return engine
}
