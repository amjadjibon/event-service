package main

import (
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"event-service/conf"
	"event-service/docs"
	"event-service/handler"
	"event-service/pkg/database"
	"event-service/pkg/httpserver"
	"event-service/pkg/logger"
	"event-service/repository"
)

const apiPathPrefix = "/api/v1"

func SetupServer(h *handler.EventHandler) *gin.Engine {
	r := gin.Default()
	docs.SwaggerInfo.BasePath = apiPathPrefix
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(requestid.New())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	apiGroup := r.Group(apiPathPrefix)

	eventGroup := apiGroup.Group("/events")
	{
		eventGroup.GET("/", h.EventList)
		eventGroup.GET("/:id", h.EventDetail)
		eventGroup.GET("/:id/workshops", h.EventWorkshopList)
	}

	workshopGroup := apiGroup.Group("/workshops")
	{
		workshopGroup.GET("/:id", h.WorkshopDetail)
		workshopGroup.POST("/:id/reservation", h.MakeReservation)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return r
}

func main() {
	cfg := conf.NewConfig()
	logger.InitLogger(cfg.LogLevel)
	db := database.GetDB(cfg)

	eventRepo := repository.NewEventRepository(db)
	eventHandler := handler.NewEventHandler(eventRepo)

	r := SetupServer(eventHandler)
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	httpserver.Run(addr, r)
}
