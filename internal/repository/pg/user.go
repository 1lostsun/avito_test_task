package pg

import (
	"avito_test_task/internal/entity"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log/slog"
)

// SetIsActive обновляет активность пользователя
func (r *Repository) SetIsActive(ctx context.Context, userID string, isActive bool) (*entity.User, error) {
	user := &entity.User{}
	err := r.pg.QueryRow(ctx, `
		UPDATE users SET is_active = $1, updated_at = now()
		WHERE user_id = $2
		RETURNING user_id, username, 
			(SELECT team_name FROM teams WHERE id = users.team_id),
			 is_active, created_at, updated_at
		`,
		isActive, userID).Scan(&user.UserID,
		&user.Name,
		&user.TeamName,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.Error("user not found", "error", err, "user_id", userID)
			return nil, fmt.Errorf("NOT_FOUND: user not found")
		}
		slog.Error(fmt.Sprintf("error updating user: %v", err))
		return nil, err
	}

	return user, nil
}

// GetUser получает пользователя с именем команды
func (r *Repository) GetUser(ctx context.Context, userID string) (*entity.User, error) {
	user := &entity.User{}
	err := r.pg.QueryRow(ctx, `
		SELECT u.user_id, u.username, t.team_name, u.is_active, u.created_at, u.updated_at
		FROM users u 
		JOIN teams t ON t.id = u.team_id
        WHERE user_id = $1
		`,
		userID).Scan(&user.UserID,
		&user.Name,
		&user.TeamName,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt)
	if err != nil {
		slog.Error(fmt.Sprintf("error getting user: %v", err))
		return nil, err
	}

	return user, nil
}

// GetActiveCandidates получает активных пользователей для назначения ревьюверов
func (r *Repository) GetActiveCandidates(ctx context.Context, teamName string, excludeUserIDs []string) ([]*entity.User, error) {
	rows, err := r.pg.Query(ctx, `
			SELECT u.user_id, u.username, t.team_name, u.is_active 
			FROM users u
			JOIN teams t ON t.id = u.team_id
			WHERE t.team_name = $1 
				AND u.is_active = true
				AND u.user_id != ALL($2)
		`,
		teamName, excludeUserIDs)

	if err != nil {
		slog.Error(fmt.Sprintf("error getting active candidates: %v", err))
		return nil, err
	}

	defer rows.Close()

	var users []*entity.User
	for rows.Next() {
		user := &entity.User{}
		if err := rows.Scan(&user.UserID, &user.Name, &user.TeamName, &user.IsActive); err != nil {
			slog.Error(fmt.Sprintf("error scanning user row: %v", err))
			return nil, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		slog.Error(fmt.Sprintf("error scanning user rows: %v", err))
		return nil, err
	}

	return users, nil
}

// GetReview получает PR, где пользователь назначен ревьюером
func (r *Repository) GetReview(ctx context.Context, userID string) ([]*entity.PullRequestShort, error) {
	rows, err := r.pg.Query(ctx, `
		SELECT pr.pull_request_id, pr.pull_request_name, pr.author_id, pr.status 
		FROM pull_requests pr
		JOIN pull_request_reviewers prr ON prr.pr_id = pr.pull_request_id
		WHERE prr.user_id = $1
		ORDER BY pr.created_at DESC
	`, userID)

	if err != nil {
		slog.Error(fmt.Sprintf("error getting pull requests: %v", err))
		return nil, err
	}

	defer rows.Close()

	prs := make([]*entity.PullRequestShort, 0)
	for rows.Next() {
		pr := &entity.PullRequestShort{}
		slog.Info(userID)
		if err := rows.Scan(&pr.ID, &pr.Name, &pr.AuthorID, &pr.Status); err != nil {
			slog.Error(fmt.Sprintf("error scanning pull request: %v", err))
			return nil, err
		}

		prs = append(prs, pr)
	}

	if err := rows.Err(); err != nil {
		slog.Error("error iterating rows", "error", err)
		return nil, err
	}

	return prs, nil
}
