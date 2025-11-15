package usecase

import (
	"avito_test_task/internal/entity"
	"context"
	"fmt"
	"log/slog"
)

// CreateTeam валидирует команду и создает ее
func (uc *UseCase) CreateTeam(ctx context.Context, team *entity.Team) (*entity.Team, error) {
	exists, err := uc.repo.TeamExists(ctx, team.Name)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, fmt.Errorf("TEAM_EXISTS: team with name %s already exists", team.Name)
	}

	if err := validateTeam(team); err != nil {
		return nil, err
	}

	if err := uc.repo.CreateTeam(ctx, team); err != nil {
		slog.Error("error to create team", "error", err, "team", team.Name)
		return nil, err
	}

	slog.Info("team successfully created", "team", team.Name, "members", team.Members)
	return team, nil
}

// GetTeam получает команду по ее имени
func (uc *UseCase) GetTeam(ctx context.Context, teamName string) (*entity.Team, error) {
	team, err := uc.repo.GetTeam(ctx, teamName)
	if err != nil {
		slog.Error("failed to get team", "error", err, "teamName", teamName)
		return nil, fmt.Errorf("NOT_FOUND: team with name %s not found", teamName)
	}

	slog.Info("team successfully retrieved", "teamName", teamName)
	return team, nil
}

func validateTeam(team *entity.Team) error {
	if team.Name == "" {
		return fmt.Errorf("team name is required")
	}

	if len(team.Members) == 0 {
		return fmt.Errorf("team must have at least one member")
	}

	seen := make(map[string]bool)
	for _, member := range team.Members {
		if member.UserID == "" {
			return fmt.Errorf("member user_id is required")
		}

		if member.Name == "" {
			return fmt.Errorf("member username must be required")
		}

		if seen[member.UserID] {
			return fmt.Errorf("team user %s is already member of team %s", member.UserID, team.Name)
		}

		seen[member.UserID] = true
	}

	return nil
}
