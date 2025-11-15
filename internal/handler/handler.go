package handler

import (
	"avito_test_task/internal/middleware"
	"avito_test_task/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Handler struct {
	uc *usecase.UseCase
}

func New(uc *usecase.UseCase) *Handler {
	return &Handler{
		uc: uc,
	}
}

func (h *Handler) InitRoutes(r *gin.Engine) {
	r.Use(middleware.PrometheusMiddleware())

	requireAdmin := middleware.RequireAdmin()
	requireUserOrAdmin := middleware.RequireAdminOrUser()

	// Teams
	team := r.Group("avito-test-task/team")
	{
		team.POST("/add", requireAdmin, h.CreateTeam)
		team.GET("/get", requireUserOrAdmin, h.GetTeam)
	}

	// Users
	users := r.Group("avito-test-task/users")
	{
		users.POST("/setIsActive", requireAdmin, h.SetIsActive)
		users.GET("/getReview", requireUserOrAdmin, h.GetReviews)
	}

	// Pull Requests
	pullRequests := r.Group("avito-test-task/pullRequest")
	{
		pullRequests.POST("/create", requireAdmin, h.CreatePR)
		pullRequests.POST("/merge", requireAdmin, h.MergePR)
		pullRequests.POST("/reassign", requireAdmin, h.ReassignPR)
	}

	// metrics endpoint
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
}
