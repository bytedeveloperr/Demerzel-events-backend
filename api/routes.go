package api

import (
	"demerzel-events/internal/handlers"
	"demerzel-events/internal/routes"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func BuildRoutesHandler() *gin.Engine {
	r := gin.New()

	if os.Getenv("APP_ENV") == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.Default())

	r.GET("/health", handlers.HealthHandler)

	routes.CommentRoutes(r)

	// OAuth routes
	oauthRoutes := r.Group("/oauth")

	oauthRoutes.GET("/initialize", handlers.InitalizeOAuthSignIn)
	oauthRoutes.GET("/callback", handlers.HandleOAuthCallBack)

	// All other API routes should be mounted on this route group
	apiRoutes := r.Group("/api")

	// mount the API routes auth middleware
	apiRoutes.Use(AuthMiddleware())

	// User routes
	apiRoutes.GET("/users/:id", handlers.GetUserById)
	apiRoutes.PUT("/users/:id", handlers.UpdateUser)
	apiRoutes.PUT("/user/:id", handlers.UpdateComment)

	//comment routes
	apiRoutes.POST("/events/:event_id/comments", handlers.CreateComment)
	apiRoutes.DELETE("/events/:event_id/comments/:comment_id", handlers.DeleteComment)

	return r
}
