package handler

import (
	"avito_test_task/internal/entity"
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"strings"
)

type CreatePRRequest struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
}

type MergePRRequest struct {
	PullRequestID string `json:"pull_request_id"`
}

type ReassignPRRequest struct {
	PullRequestID string `json:"pull_request_id"`
	OldUserID     string `json:"old_user_id"`
}

// CreatePR POST /pullRequest/create
func (h *Handler) CreatePR(c *gin.Context) {
	var req CreatePRRequest
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

	pr, err := h.uc.CreatePR(c.Request.Context(), req.PullRequestID, req.PullRequestName, req.AuthorID)
	if err != nil {
		if strings.Contains(err.Error(), "PR_EXISTS") {
			slog.Error("pull request already exists", "error", err)
			c.JSON(http.StatusConflict, gin.H{
				"error": gin.H{
					"code":    entity.PRExists,
					"message": fmt.Sprintf("pull request already exists: %s", req.PullRequestID),
				},
			})
			return
		}

		if strings.Contains(err.Error(), "NOT_FOUND") {
			slog.Error("pull request author not found", "error", err)
			c.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{
					"code":    entity.NotFound,
					"message": fmt.Sprintf("pull request author not found: %s", req.AuthorID),
				},
			})
			return
		}

		slog.Error("failed to create pull request", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    entity.InternalServer,
				"message": fmt.Errorf("internal server error: %w", err).Error(),
			},
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"pr": *pr,
	})
}

// MergePR POST /pullRequest/merge
func (h *Handler) MergePR(c *gin.Context) {
	var req MergePRRequest
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

	pr, err := h.uc.MergePR(c.Request.Context(), req.PullRequestID)
	if err != nil {
		if strings.Contains(err.Error(), "NOT_FOUND") {
			slog.Error("pull request not found", "error", err)
			c.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{
					"code":    entity.NotFound,
					"message": fmt.Sprintf("pull request not found: %s", req.PullRequestID),
				},
			})
			return
		}

		slog.Error("failed to merge pull request", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    entity.InternalServer,
				"message": fmt.Errorf("internal server error: %w", err).Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"pr": *pr,
	})
}

// ReassignPR POST /pullRequest/reassign
func (h *Handler) ReassignPR(c *gin.Context) {
	var req ReassignPRRequest
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

	pr, newReviewerID, err := h.uc.ReassignReviewer(c.Request.Context(), req.PullRequestID, req.OldUserID)
	if err != nil {
		if strings.Contains(err.Error(), "NOT_FOUND") {
			slog.Error("pull request not found", "error", err)
			c.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{
					"code":    entity.NotFound,
					"message": fmt.Sprintf("pull request not found: %s", req.PullRequestID),
				},
			})
			return
		}

		if strings.Contains(err.Error(), "PR_MERGED") {
			slog.Error("pull request already merged", "error", err)
			c.JSON(http.StatusConflict, gin.H{
				"error": gin.H{
					"code":    entity.PRMerged,
					"message": fmt.Sprintf("pull request already merged: %s", req.PullRequestID),
				},
			})
			return
		}

		if strings.Contains(err.Error(), "NOT_ASSIGNED") {
			slog.Error("reviewer is not assigned to this PR", "error", err)
			c.JSON(http.StatusConflict, gin.H{
				"error": gin.H{
					"code":    entity.NotAssigned,
					"message": fmt.Sprintf("reviewer is not assigned to this PR: %s", req.PullRequestID),
				},
			})
			return
		}

		if strings.Contains(err.Error(), "NO_CANDIDATE") {
			slog.Error("no candidates to assign", "error", err)
			c.JSON(http.StatusConflict, gin.H{
				"error": gin.H{
					"code":    entity.NoCandidate,
					"message": fmt.Sprintf("no candidates to assign: %s", req.PullRequestID),
				},
			})
			return
		}

		slog.Error("failed to reassign pull request", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    entity.InternalServer,
				"message": fmt.Errorf("internal server error: %w", err).Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"pr":          *pr,
		"replaced_by": newReviewerID,
	})
}
