package pg

import (
	"avito_test_task/internal/entity"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log/slog"
)

// CreatePR создает pull request
func (r *Repository) CreatePR(ctx context.Context, pr *entity.PullRequest) error {
	tx, err := r.pg.Begin(ctx)
	if err != nil {
		slog.Error(fmt.Sprintf("error "))
		return err
	}

	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		INSERT INTO pull_requests (
		    pull_request_id, 
            pull_request_name, 
            author_id, 
            status,
            created_at
		) 
		VALUES ($1, $2, $3, $4, $5)
	`, pr.ID, pr.Name, pr.AuthorID, pr.Status, pr.CreatedAt)

	if err != nil {
		slog.Error(fmt.Sprintf("error insert pull request: %v", err))
		return err
	}

	if len(pr.AssignReviewers) > 0 {
		batch := pgx.Batch{}
		for _, reviewerID := range pr.AssignReviewers {
			batch.Queue(`INSERT INTO pull_request_reviewers(
            pr_id,
            user_id,
            assigned_at)
            VALUES ($1, $2, now())
            `, pr.ID, reviewerID)
		}

		br := tx.SendBatch(ctx, &batch)

		for i := 0; i < len(pr.AssignReviewers); i++ {
			if _, err := br.Exec(); err != nil {
				br.Close()
				slog.Error("error inserting reviewer",
					"error", err,
					"index", i,
					"pr_id", pr.ID)
				return fmt.Errorf("insert reviewer: %w", err)
			}
		}

		if err := br.Close(); err != nil {
			slog.Error("error closing batch results", "error", err)
			return fmt.Errorf("close batch results: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		slog.Error(fmt.Sprintf("error committing PR: %v", err))
		return err
	}

	slog.Info("PR created successfully", "pr_id", pr.ID, "reviewers", len(pr.AssignReviewers))

	return nil
}

// MergePR помечает PR как merged
func (r *Repository) MergePR(ctx context.Context, prID string) (*entity.PullRequest, error) {
	pr := &entity.PullRequest{}

	err := r.pg.QueryRow(ctx, `
		UPDATE pull_requests
		SET status = 'MERGED', merged_at = now()
		WHERE pull_request_id = $1
		RETURNING pull_request_id, pull_request_name, author_id, status, created_at, merged_at
	`, prID).Scan(&pr.ID, &pr.Name, &pr.AuthorID, &pr.Status, &pr.CreatedAt, &pr.MergedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.Error(fmt.Sprintf("PR not found: %v", err))
			return nil, errors.New("PR not found")
		}

		slog.Error(fmt.Sprintf("error merge pull request: %v", err))
		return nil, err
	}

	rows, err := r.pg.Query(ctx, `
		SELECT user_id FROM pull_request_reviewers 
		WHERE pr_id = $1
		ORDER BY assigned_at DESC
		`, prID)
	if err != nil {
		slog.Error(fmt.Sprintf("error getting reviewers: %v", err))
		return nil, err
	}

	defer rows.Close()

	pr.AssignReviewers = make([]string, 0)
	for rows.Next() {
		var reviewerID string
		if err := rows.Scan(&reviewerID); err != nil {
			slog.Error(fmt.Sprintf("error scanning reviewer id: %v", err))
			return nil, err
		}
		pr.AssignReviewers = append(pr.AssignReviewers, reviewerID)
	}

	if err := rows.Err(); err != nil {
		slog.Error(fmt.Sprintf("error rows iteration: %v", err))
		return nil, err
	}

	return pr, nil
}

func (r *Repository) ReassignReviewer(ctx context.Context, prID, oldReviewerID, newReviewerID string) (*entity.PullRequest, error) {
	pr, err := r.GetPR(ctx, prID)
	if err != nil {
		return nil, err
	}

	if pr.Status == entity.MERGED {
		slog.Error("PR is already merged", "prID", prID)
		return nil, fmt.Errorf("PR_MERGED: cannot reassign on merged PR")
	}

	if !contains(pr.AssignReviewers, oldReviewerID) {
		slog.Error("Reviewer is not assigned to PR", "prID", prID, "reviewerID", oldReviewerID)
		return nil, fmt.Errorf("NOT_ASSIGNED: reviewer is not assigned to this PR")
	}

	_, err = r.pg.Exec(ctx, `
		UPDATE pull_request_reviewers
		SET user_id = $1, assigned_at = now()
		WHERE pr_id = $2 AND user_id = $3
	`, newReviewerID, prID, oldReviewerID)

	if err != nil {
		slog.Error(fmt.Sprintf("failed to update reviewer: %v", err))
		return nil, err
	}

	for i, reviewerID := range pr.AssignReviewers {
		if reviewerID == oldReviewerID {
			pr.AssignReviewers[i] = newReviewerID
			break
		}
	}

	return pr, nil
}

// GetPR вспомогательная функция получения pull request
func (r *Repository) GetPR(ctx context.Context, prID string) (*entity.PullRequest, error) {
	pr := &entity.PullRequest{}

	err := r.pg.QueryRow(ctx, `
		SELECT pull_request_id, pull_request_name, author_id, status, created_at, merged_at
		FROM pull_requests
		WHERE pull_request_id = $1
	`, prID).Scan(&pr.ID, &pr.Name, &pr.AuthorID, &pr.Status, &pr.CreatedAt, &pr.MergedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.Error(fmt.Sprintf("PR not found: %v", err))
			return nil, errors.New("PR not found")
		}

		slog.Error(fmt.Sprintf("error getting pull request: %v", err))
		return nil, err
	}

	rows, err := r.pg.Query(ctx, `
		SELECT user_id FROM pull_request_reviewers
		WHERE pr_id = $1
	`, prID)

	if err != nil {
		slog.Error(fmt.Sprintf("error getting reviewers: %v", err))
		return nil, err
	}

	defer rows.Close()

	pr.AssignReviewers = make([]string, 0)
	for rows.Next() {
		var reviewerID string
		if err := rows.Scan(&reviewerID); err != nil {
			slog.Error(fmt.Sprintf("error scanning reviewer id: %v", err))
			return nil, err
		}

		pr.AssignReviewers = append(pr.AssignReviewers, reviewerID)
	}

	if err := rows.Err(); err != nil {
		slog.Error(fmt.Sprintf("error rows iteration: %v", err))
		return nil, err
	}

	return pr, nil
}

// contains проверяет, содержит ли слайс строку
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// PRExists проверяет pull request на существование
func (r *Repository) PRExists(ctx context.Context, prID string) (bool, error) {
	var exists bool
	err := r.pg.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM pull_requests WHERE pull_request_id = $1)", prID).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.Error(fmt.Sprintf("PR not found: %v", err))
			return false, nil
		}

		slog.Error(fmt.Sprintf("error checking if PR exists: %v", err))
		return false, err
	}

	return exists, nil
}

// GetReviewersWorkload возвращает количество открытых PR для каждого ревьювера
func (r *Repository) GetReviewersWorkload(ctx context.Context, userIDs []string) (map[string]int, error) {
	if len(userIDs) == 0 {
		return map[string]int{}, nil
	}

	rows, err := r.pg.Query(ctx, `
		SELECT prr.pr_id, COUNT(*) 
		FROM pull_request_reviewers prr
		JOIN pull_requests pr ON pr.pull_request_id = prr.pr_id
		WHERE prr.user_id = ANY($1)
		AND pr.status = 'OPEN'
		GROUP BY prr.pr_id
	`, userIDs)

	if err != nil {
		slog.Error(fmt.Sprintf("error getting reviewers workload: %v", err))
		return nil, err
	}

	defer rows.Close()

	workload := make(map[string]int)
	for _, reviewerID := range userIDs {
		workload[reviewerID] = 0
	}

	for rows.Next() {
		var reviewerID string
		var count int
		if err := rows.Scan(&reviewerID, &count); err != nil {
			slog.Error(fmt.Sprintf("error scanning reviewer workload: %v", err))
			return nil, err
		}

		workload[reviewerID] = count
	}

	return workload, nil
}
