package routes

import (
	"onlineChat/internal/users"
	"onlineChat/internal/ws"
	"onlineChat/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func SetupRoutes(userHandler *users.Handler, wsHandler *ws.Handler, config *Config, logger *logrus.Logger) *gin.Engine {
	r := gin.New()

	r.Use(middleware.RequestLogger(logger))
	r.Use(middleware.ErrorHandler(logger))
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.CORS(config.Security.CORSOrigin))
	r.Use(middleware.RateLimit(config.Security.RateLimitRequests, 10))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	auth := r.Group("/auth")
	{
		auth.POST("/register", userHandler.Register)
		auth.POST("/login", userHandler.Login)
	}

	protected := r.Group("/")
	protected.Use(middleware.NewAuthMiddleware(config.JWT.Secret, logger).RequireAuth())
	{
		protected.GET("/auth/profile", userHandler.GetProfile)
		protected.PUT("/auth/profile", userHandler.UpdateProfile)
		protected.DELETE("/auth/account", userHandler.DeleteAccount)

		users := protected.Group("/users")
		{
			users.GET("/:id", userHandler.GetUserByID)
		}

		chats := protected.Group("/chats")
		{
			chats.POST("/", wsHandler.CreateChat)
			chats.GET("/", wsHandler.GetAllChats)
			chats.GET("/search", wsHandler.SearchPublicChats)
			chats.POST("/:chatID/join", wsHandler.JoinChat)
			chats.POST("/:chatID/leave", wsHandler.LeaveChat)
			chats.GET("/:chatID/clients", wsHandler.GetClientsByChatID)
			chats.GET("/:chatID/messages", wsHandler.GetChatMessages)
			chats.GET("/:chatID/ws", wsHandler.ServeWS)
		}
	}

	return r
}

type Config struct {
	JWT      JWTConfig
	Security SecurityConfig
}

type JWTConfig struct {
	Secret string
}

type SecurityConfig struct {
	CORSOrigin        string
	RateLimitRequests int
}
