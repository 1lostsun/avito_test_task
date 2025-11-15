package handler

import (
	"avito_test_task/internal/entity"
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"strings"
)

type SetIsActiveRequest struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

// SetIsActive POST /users/setIsActive
func (h *Handler) SetIsActive(c *gin.Context) {
	var req SetIsActiveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("failed to parse SetIsActiveRequest:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    entity.InvalidRequest,
				"message": fmt.Errorf("invalid request body: %w", err).Error(),
			},
		})
		return
	}

	user, err := h.uc.SetIsActive(c.Request.Context(), req.UserID, req.IsActive)
	if err != nil {
		if strings.Contains(err.Error(), "NOT_FOUND") {
			slog.Error("user not found", "error", err, "userID", req.UserID)
			c.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{
					"code":    entity.NotFound,
					"message": fmt.Errorf("user not found: %w", err).Error(),
				},
			})
			return
		}

		slog.Error("failed to set IsActive:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    entity.InternalServer,
				"message": fmt.Errorf("internal server error: %w", err).Error(),
			},
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": &user,
	})
}

// GetReviews Get /users/getReview
func (h *Handler) GetReviews(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		slog.Error("user id must be required")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    entity.InvalidRequest,
				"message": "invalid request body",
			},
		})
	}

	prs, err := h.uc.GetUserReviews(c.Request.Context(), userID)
	if err != nil {
		slog.Error("failed to get user reviews:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    entity.InternalServer,
				"message": fmt.Errorf("internal server error: %w", err).Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"reviews": prs,
	})
}
