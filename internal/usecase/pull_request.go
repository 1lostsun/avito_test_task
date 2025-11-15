package usecase

import (
	"avito_test_task/internal/entity"
	"context"
	"fmt"
	"log/slog"
	"sort"
	"time"
)

// CreatePR создает pull request и назначает наименее загруженных ревьюеров
func (uc *UseCase) CreatePR(ctx context.Context, prID, prName, authorID string) (*entity.PullRequest, error) {
	exists, err := uc.repo.PRExists(ctx, prID)
	if err != nil {
		slog.Error("failed to check existence of PR", "error", err)
		return nil, err
	}

	if exists {
		slog.Error("PR already exists")
		return nil, fmt.Errorf("PR_EXISTS: PR id already exists")
	}

	user, err := uc.repo.GetUser(ctx, authorID)
	if err != nil {
		slog.Error("author not found", "error", err, "authorID", authorID)
		return nil, fmt.Errorf("NOT_FOUND: author not found")
	}

	candidates, err := uc.repo.GetActiveCandidates(ctx, user.TeamName, []string{authorID})
	if err != nil {
		slog.Error("failed to get active candidates", "error", err)
		return nil, fmt.Errorf("get candidates: %w", err)
	}

	reviewers, err := uc.selectLeastBusyReviewers(ctx, candidates, 2)
	if err != nil {
		slog.Error("failed to get reviewers", "error", err)
		return nil, err
	}

	now := time.Now()
	pr := &entity.PullRequest{
		ID:              prID,
		Name:            prName,
		AuthorID:        authorID,
		Status:          entity.OPEN,
		AssignReviewers: reviewers,
		CreatedAt:       &now,
		MergedAt:        nil,
	}

	if err := uc.repo.CreatePR(ctx, pr); err != nil {
		slog.Error("failed to create PR", "error", err)
		return nil, err
	}

	slog.Info("PR created successfully", "prID", prID, "reviewersCount", len(reviewers))

	return pr, nil
}

// selectLeastBusyReviewers выбирает N наименее загруженных ревьюеров
func (uc *UseCase) selectLeastBusyReviewers(ctx context.Context, candidates []*entity.User, countReviewers int) ([]string, error) {
	if len(candidates) == 0 {
		slog.Error("no candidates on review")
		return []string{}, nil
	}

	userIDs := make([]string, len(candidates))
	for i, c := range candidates {
		userIDs[i] = c.UserID
	}

	workload, err := uc.repo.GetReviewersWorkload(ctx, userIDs)
	if err != nil {
		slog.Error("failed to get reviewers workload", "error", err)
		return nil, err
	}

	workloads := make([]reviewerWorkload, 0, len(workload))
	for id, w := range workload {
		workloads = append(workloads, reviewerWorkload{
			UserID: id,
			Count:  w,
		})
	}

	sort.Slice(workloads, func(i, j int) bool {
		return workloads[i].Count < workloads[j].Count
	})

	count := minimum(countReviewers, len(candidates))
	reviewers := make([]string, count)
	for i := 0; i < count; i++ {
		reviewers[i] = candidates[i].UserID
	}

	return reviewers, nil
}

type reviewerWorkload struct {
	UserID string
	Count  int
}

func minimum(a, b int) int {
	if a < b {
		return a
	}

	return b
}

// MergePR помечает pull request как MERGED
func (uc *UseCase) MergePR(ctx context.Context, prID string) (*entity.PullRequest, error) {
	pr, err := uc.repo.MergePR(ctx, prID)
	if err != nil {
		slog.Error("failed to merge PR", "error", err)
		return nil, fmt.Errorf("NOT_FOUND: PR not found")
	}

	slog.Info("PR merged successfully", "prID", prID)
	return pr, nil
}

// ReassignReviewer переназначает ревьюера
func (uc *UseCase) ReassignReviewer(ctx context.Context, prID, oldReviewerID string) (*entity.PullRequest, string, error) {
	pr, err := uc.repo.GetPR(ctx, prID)
	if err != nil {
		slog.Error("failed to get PR", "error", err)
		return nil, "", fmt.Errorf("NOT_FOUND: PR not found")
	}

	user, err := uc.repo.GetUser(ctx, oldReviewerID)
	if err != nil {
		slog.Error("failed to get user", "error", err)
		return nil, "", fmt.Errorf("NOT_FOUND: author not found")
	}

	excludeIDs := append(pr.AssignReviewers, pr.AuthorID)
	candidates, err := uc.repo.GetActiveCandidates(ctx, user.TeamName, excludeIDs)
	if err != nil {
		slog.Error("failed to get active candidates", "error", err)
		return nil, "", fmt.Errorf("get active candidates error: %v", err)
	}

	if len(candidates) == 0 {
		slog.Error("no candidates on review")
		return nil, "", fmt.Errorf("NO_CANDIDATE: no active replacement candidate in team")
	}

	newReviewerID, err := uc.selectLeastBusyReviewer(ctx, candidates)
	if err != nil {
		slog.Error("failed to select reviewer", "error", err)
		return nil, "", fmt.Errorf("select reviewer error: %v", err)
	}

	updatedPR, err := uc.repo.ReassignReviewer(ctx, prID, oldReviewerID, newReviewerID)
	if err != nil {
		slog.Error("failed to reassign reviewer", "error", err)
		return nil, "", fmt.Errorf("reviewer reassign error: %v", err)
	}

	slog.Info("reviewer reassigned",
		"prID", prID,
		"oldReviewer", oldReviewerID,
		"newReviewer", newReviewerID)

	return updatedPR, newReviewerID, nil
}

// selectLeastBusyReviewer выбирает одного наименее загруженного ревьюера
func (uc *UseCase) selectLeastBusyReviewer(ctx context.Context, candidates []*entity.User) (string, error) {
	if len(candidates) == 0 {
		slog.Error("no available candidates on review")
		return "", fmt.Errorf("NOT_FOUND: no available candidates on review")
	}

	userIDs := make([]string, len(candidates))
	for i, c := range candidates {
		userIDs[i] = c.UserID
	}

	workload, err := uc.repo.GetReviewersWorkload(ctx, userIDs)
	if err != nil {
		slog.Error("failed to get reviewers workload", "error", err)
		return "", fmt.Errorf("NOT_FOUND: reviewers workload not found")
	}

	minLoad := -1
	reviewerID := ""

	for _, c := range candidates {
		load := workload[c.UserID]
		if minLoad == -1 || load < minLoad {
			minLoad = load
			reviewerID = c.UserID
		}
	}

	return reviewerID, nil
}
