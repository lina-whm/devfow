package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/devflow/devflow-backend/internal/api/docs"
	"github.com/devflow/devflow-backend/internal/api/handler"
	"github.com/devflow/devflow-backend/internal/api/middleware"
)

func Setup(
	ah *handler.AuthHandler,
	hh *handler.HealthHandler,
	th *handler.TaskHandler,
	oh *handler.OrgHandler,
	ph *handler.ProjectHandler,
	ch *handler.CommentHandler,
	bh *handler.BoardHandler,
	jwtSecret string,
) *gin.Engine {
	r := gin.New()

	r.Use(middleware.CORS("*"))
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger())
	r.Use(gin.Recovery())

	api := r.Group("/api/v1")
	{
		health := api.Group("/health")
		{
			health.GET("/live", hh.Live)
			health.GET("/ready", hh.Ready)
		}

		auth := api.Group("/auth")
		{
			auth.POST("/register", ah.Register)
			auth.POST("/login", ah.Login)
			auth.POST("/refresh", ah.Refresh)
			auth.POST("/logout", ah.Logout)
			auth.POST("/verify-email", ah.VerifyEmail)
			auth.POST("/forgot-password", ah.ForgotPassword)
			auth.POST("/reset-password", ah.ResetPassword)
		}

		authed := api.Group("")
		authed.Use(middleware.Auth(jwtSecret))
		{
			authed.GET("/auth/me", ah.Me)
			authed.PATCH("/auth/me", ah.UpdateMe)
			orgs := authed.Group("/organizations")
			{
				orgs.POST("", oh.Create)
				orgs.GET("", oh.List)
				orgs.GET("/:orgId", oh.GetByID)
				orgs.PUT("/:orgId", oh.Update)
				orgs.DELETE("/:orgId", oh.Delete)
			}

			orgProjects := authed.Group("/organizations/:orgId/projects")
			{
				orgProjects.POST("", ph.Create)
				orgProjects.GET("", ph.List)
			}

			projects := authed.Group("/projects")
			{
				projects.GET("/:projectId", ph.GetByID)
				projects.PUT("/:projectId", ph.Update)
				projects.DELETE("/:projectId", ph.Delete)
				projects.GET("/:projectId/board", bh.GetByProject)
			}

			tasks := authed.Group("/projects/:projectId/tasks")
			{
				tasks.GET("", th.List)
				tasks.POST("/create", th.Create)
				tasks.GET("/:taskId", th.GetByID)
				tasks.PUT("/:taskId", th.Update)
				tasks.PUT("/:taskId/move", th.Move)
				tasks.DELETE("/:taskId", th.Delete)
				tasks.PUT("/:taskId/restore", th.Restore)

				tasks.POST("/:taskId/comments", ch.Create)
				tasks.GET("/:taskId/comments", ch.ListByTask)
			}

			authed.GET("/tasks/:taskId", th.GetByID)
			authed.PATCH("/tasks/:taskId", th.Update)
			authed.PUT("/tasks/:taskId", th.Update)
			authed.DELETE("/tasks/:taskId", th.Delete)
			authed.PATCH("/tasks/:taskId/move", th.Move)
			authed.PUT("/tasks/:taskId/move", th.Move)
			authed.GET("/tasks/:taskId/comments", ch.ListByTask)
			authed.POST("/tasks/:taskId/comments", ch.Create)

			boards := authed.Group("/boards")
			{
				boards.PATCH("/:boardId/columns", bh.UpdateColumns)
			}

			comments := authed.Group("/comments")
			{
				comments.GET("/:commentId", ch.GetByID)
				comments.PUT("/:commentId", ch.Update)
				comments.DELETE("/:commentId", ch.Delete)
			}

			authed.GET("/notifications", func(c *gin.Context) {
				c.JSON(200, gin.H{"data": []interface{}{}, "total": 0, "page": 1, "pageSize": 20, "totalPages": 0})
			})
			authed.GET("/notifications/unread-count", func(c *gin.Context) {
				c.JSON(200, 0)
			})
			authed.PATCH("/notifications/:id/read", func(c *gin.Context) {
				c.JSON(200, gin.H{"status": "ok"})
			})
			authed.PATCH("/notifications/read-all", func(c *gin.Context) {
				c.JSON(200, gin.H{"status": "ok"})
			})
		}
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
