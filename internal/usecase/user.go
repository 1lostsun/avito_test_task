package usecase

import (
	"avito_test_task/internal/entity"
	"context"
	"log/slog"
)

// SetIsActive устанавливает активность пользователя
func (uc *UseCase) SetIsActive(ctx context.Context, userID string, active bool) (*entity.User, error) {
	user, err := uc.repo.SetIsActive(ctx, userID, active)
	if err != nil {
		slog.Error("failed to set user active status", "error", err, "userID", userID)
		return nil, err
	}

	slog.Info("user active status updated", "status", active, "userID", userID)
	return user, nil
}

// GetUserReviews получает все pull request, назначенные пользователю на ревью
func (uc *UseCase) GetUserReviews(ctx context.Context, userID string) ([]*entity.PullRequestShort, error) {
	prs, err := uc.repo.GetReview(ctx, userID)
	if err != nil {
		slog.Error("failed to get user reviews", "error", err, "userID", userID)
		return nil, err
	}

	slog.Info("successfully get user reviews", "pull requests", prs)
	return prs, nil
}
