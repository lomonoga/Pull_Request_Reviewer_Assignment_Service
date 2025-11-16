package repository

import (
	"database/sql"
	"pull_requests_service/internal/domain"
)

type userRepository struct {
	BaseRepository
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &userRepository{BaseRepository{db: db}}
}

func (r *userRepository) CreateOrUpdateUser(user *domain.TeamMember) error {
	_, err := r.db.Exec(`
        INSERT INTO "user" (user_id, username, is_active) 
        VALUES ($1, $2, $3)
        ON CONFLICT (user_id) 
        DO UPDATE SET username = EXCLUDED.username, is_active = EXCLUDED.is_active`,
		user.UserID, user.Username, user.IsActive,
	)

	return err
}

func (r *userRepository) GetUserTeam(userID string) (string, error) {
	var teamName string
	err := r.db.QueryRow(`
        SELECT team_name FROM team_member WHERE user_id = $1`,
		userID,
	).Scan(&teamName)
	if err == sql.ErrNoRows {
		return "", domain.ErrUserNotFound
	}
	return teamName, err
}

func (r *userRepository) SetUserActive(userID string, isActive bool) error {
	_, err := r.db.Exec(`
        UPDATE "user" SET is_active = $1 WHERE user_id = $2`,
		isActive, userID,
	)

	return err
}

func (r *userRepository) GetUser(userID string) (*domain.TeamMember, error) {
	var user domain.TeamMember
	err := r.db.QueryRow(`
        SELECT user_id, username, is_active FROM "user" WHERE user_id = $1`,
		userID,
	).Scan(&user.UserID, &user.Username, &user.IsActive)
	if err == sql.ErrNoRows {
		return nil, domain.ErrUserNotFound
	}
	return &user, err
}

func (r *userRepository) GetActiveTeamMembers(teamName string) ([]domain.TeamMember, error) {
	rows, err := r.db.Query(`
        SELECT u.user_id, u.username, u.is_active 
        FROM "user" u 
        JOIN team_member tm ON u.user_id = tm.user_id 
        WHERE tm.team_name = $1 AND u.is_active = true`,
		teamName,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []domain.TeamMember
	for rows.Next() {
		var member domain.TeamMember
		if err := rows.Scan(&member.UserID, &member.Username, &member.IsActive); err != nil {
			return nil, err
		}
		members = append(members, member)
	}

	return members, nil
}
