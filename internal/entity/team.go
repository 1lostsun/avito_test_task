package entity

import "time"

type Team struct {
	Name      string        `json:"team_name"`
	Members   []*TeamMember `json:"team_members"`
	CreatedAt *time.Time
}

type TeamMember struct {
	UserID   string `json:"user_id"`
	Name     string `json:"username"`
	IsActive bool   `json:"is_active"`
}
