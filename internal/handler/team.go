package handler

import (
	"avito_test_task/internal/entity"
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"strings"
)

// CreateTeam POST /team/add
func (h *Handler) CreateTeam(c *gin.Context) {
	var req entity.Team
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("failed to parse request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    entity.InvalidRequest,
				"message": fmt.Errorf("invalid request body: %w", err).Error(),
			},
		})
		return
	}

	team, err := h.uc.CreateTeam(c.Request.Context(), &req)
	if err != nil {
		if strings.Contains(err.Error(), "TEAM_EXISTS") {
			slog.Error("team already exists", "error", err)
			c.JSON(http.StatusConflict, gin.H{
				"error": gin.H{
					"code":    entity.TeamExists,
					"message": fmt.Errorf("team already exists").Error(),
				},
			})
			return
		}

		slog.Error("failed to create team", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    entity.InternalServer,
				"message": fmt.Errorf("internal server error: %w", err).Error(),
			},
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"team": team,
	})
}

// GetTeam GET /team/get
func (h *Handler) GetTeam(c *gin.Context) {
	teamName := c.Query("team_name")
	if teamName == "" {
		slog.Error("team_name must be required")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    entity.InvalidRequest,
				"message": fmt.Errorf("invalid request body").Error(),
			},
		})
		return
	}

	team, err := h.uc.GetTeam(c.Request.Context(), teamName)
	if err != nil {
		if strings.Contains(err.Error(), "NOT_FOUND") {
			slog.Error("team not found", "error", err)
			c.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{
					"code":    entity.NotFound,
					"message": fmt.Errorf("team not found: %w", err).Error(),
				},
			})
			return
		}

		slog.Error("failed to get team", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    entity.InternalServer,
				"message": fmt.Errorf("internal server error: %w", err).Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"team": *team,
	})
}
