package usecase

import (
	"avito_test_task/internal/entity"
	"avito_test_task/internal/repository/pg"
	"context"
)

// UseCase содержит всю бизнес-логику
type UseCase struct {
	repo *pg.Repository
}

// New создаёт новый use case
func New(repo *pg.Repository) *UseCase {
	return &UseCase{
		repo: repo,
	}
}

type Repository interface {
	// Teams
	CreateTeam(ctx context.Context, team *entity.Team) error
	GetTeam(ctx context.Context, teamName string) (*entity.Team, error)
	TeamExists(ctx context.Context, teamName string) (bool, error)

	// Users
	GetUser(ctx context.Context, userID string) (*entity.User, error)
	SetIsActive(ctx context.Context, userID string, isActive bool) (*entity.User, error)
	GetActiveCandidates(ctx context.Context, teamName string, excludeUserIDs []string) ([]*entity.User, error)

	// PRs
	CreatePR(ctx context.Context, pr *entity.PullRequest) error
	GetPR(ctx context.Context, prID string) (*entity.PullRequest, error)
	MergePR(ctx context.Context, prID string) (*entity.PullRequest, error)
	ReassignReviewer(ctx context.Context, prID, oldReviewerID, newReviewerID string) (*entity.PullRequest, error)
	GetReview(ctx context.Context, userID string) ([]*entity.PullRequestShort, error)
	PRExists(ctx context.Context, prID string) (bool, error)
}
