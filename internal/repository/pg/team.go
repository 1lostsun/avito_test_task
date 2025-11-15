package pg

import (
	"avito_test_task/internal/entity"
	"context"
	"fmt"
	"log/slog"
	"time"
)

// CreateTeam создает команду и добавляет в нее участников
func (r *Repository) CreateTeam(ctx context.Context, team *entity.Team) error {
	var teamID int

	tx, err := r.pg.Begin(ctx)
	if err != nil {
		slog.Error(fmt.Sprintf("error starting transaction: %v", err))
		return err
	}

	defer tx.Rollback(ctx)
	err = tx.QueryRow(ctx, `
		INSERT INTO teams (team_name, created_at)
		VALUES($1, now())
		RETURNING id
		`, team.Name).Scan(&teamID)
	if err != nil {
		slog.Error(fmt.Sprintf("error inserting team: %v", err))
		return err
	}

	for _, member := range team.Members {
		_, err := tx.Exec(ctx, `
			INSERT INTO users (user_id, username, team_id, is_active)
			VALUES($1, $2, $3, $4)
			ON CONFLICT (user_id) DO UPDATE
				SET username = EXCLUDED.username,
					team_id = EXCLUDED.team_id,
					is_active = EXCLUDED.is_active,
					updated_at = now()
			`,
			member.UserID,
			member.Name,
			teamID,
			true,
		)
		if err != nil {
			slog.Error("error inserting user",
				"error", err,
				"user_id", member.UserID)
			return fmt.Errorf("insert user %s: %w", member.UserID, err)
		}
	}

	slog.Info("team created successfully")
	return tx.Commit(ctx)
}

// GetTeam получает команду с участниками
func (r *Repository) GetTeam(ctx context.Context, teamName string) (*entity.Team, error) {
	rows, err := r.pg.Query(ctx, `
		SELECT t.team_name, t.created_at,
			u.user_id, u.username, u.is_active
		FROM teams t
		LEFT JOIN users u ON u.team_id = t.id
		WHERE t.team_name = $1
		ORDER BY u.username`, teamName)
	if err != nil {
		slog.Error(fmt.Sprintf("error getting team: %v", err))
		return nil, err
	}

	defer rows.Close()

	var team *entity.Team
	var members []*entity.TeamMember

	for rows.Next() {
		var (
			teamName  string
			createdAt time.Time
			userID    *string
			username  *string
			isActive  *bool
		)

		err = rows.Scan(&teamName, &createdAt, &userID, &username, &isActive)
		if err != nil {
			slog.Error(fmt.Sprintf("error scanning row: %v", err))
			return nil, err
		}

		if team == nil {
			team = &entity.Team{
				Name:      teamName,
				CreatedAt: &createdAt,
				Members:   members,
			}
		}

		if userID != nil {
			members = append(members, &entity.TeamMember{
				UserID:   *userID,
				Name:     *username,
				IsActive: *isActive,
			})
		}
	}

	if err = rows.Err(); err != nil {
		slog.Error(fmt.Sprintf("rows iteration: %v", err))
		return nil, err
	}

	if team == nil {
		slog.Error("team not found")
		return nil, err
	}

	team.Members = members
	return team, nil
}

// TeamExists проверяет существует ли команда с определенным именем
func (r *Repository) TeamExists(ctx context.Context, teamName string) (bool, error) {
	var exists bool

	err := r.pg.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM teams WHERE team_name = $1)", teamName).Scan(&exists)
	if err != nil {
		slog.Error(fmt.Sprintf("error checking existence of team: %v", err))
		return false, err
	}

	return exists, nil
}
