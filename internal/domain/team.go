package domain

type Team struct {
	TeamName string
	Members  []TeamMember
}

type TeamMember struct {
	UserID   string
	Username string
	IsActive bool
}

type TeamRepository interface {
	CreateTeam(team *Team) error
	GetTeam(teamName string) (*Team, error)
	TeamExists(teamName string) (bool, error)
}

type UserRepository interface {
	CreateOrUpdateUser(user *TeamMember) error
	GetUserTeam(userID string) (string, error)
	SetUserActive(userID string, isActive bool) error
	GetUser(userID string) (*TeamMember, error)
	GetActiveTeamMembers(teamName string) ([]TeamMember, error)
}
