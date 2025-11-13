package entity

type User struct {
	UserID   string
	Name     string
	TeamName string
	IsActive bool
}

type TeamMember struct {
	UserID   string
	Name     string
	IsActive bool
}
