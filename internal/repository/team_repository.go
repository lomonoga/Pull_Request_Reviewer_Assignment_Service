package repository

import (
	"database/sql"
	"fmt"
	"pull_requests_service/internal/domain"
)

type teamRepository struct {
	BaseRepository
}

func NewTeamRepository(db *sql.DB) domain.TeamRepository {
	return &teamRepository{BaseRepository{db: db}}
}

func (r *teamRepository) CreateTeam(team *domain.Team) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("INSERT INTO team (team_name) VALUES ($1) ON CONFLICT (team_name) DO NOTHING", team.TeamName)
	if err != nil {
		return err
	}

	for _, member := range team.Members {
		var existingTeam string
		err = tx.QueryRow(`
            SELECT team_name FROM team_member WHERE user_id = $1 AND team_name != $2`,
			member.UserID, team.TeamName,
		).Scan(&existingTeam)

		if err == nil {
			return fmt.Errorf("user %s already belongs to team %s", member.UserID, existingTeam)
		}

		if err != sql.ErrNoRows {
			return err
		}

		_, err = tx.Exec(`
            INSERT INTO team_member (team_name, user_id) 
            VALUES ($1, $2) 
            ON CONFLICT (team_name, user_id) 
            DO UPDATE SET team_name = EXCLUDED.team_name`,
			team.TeamName, member.UserID,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *teamRepository) GetTeam(teamName string) (*domain.Team, error) {
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM team WHERE team_name = $1)", teamName).Scan(&exists)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, domain.ErrTeamNotFound
	}

	rows, err := r.db.Query(`
        SELECT u.user_id, u.username, u.is_active 
        FROM "user" u 
        JOIN team_member tm ON u.user_id = tm.user_id 
        WHERE tm.team_name = $1`,
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

	return &domain.Team{
		TeamName: teamName,
		Members:  members,
	}, nil
}

func (r *teamRepository) TeamExists(teamName string) (bool, error) {
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM team WHERE team_name = $1)", teamName).Scan(&exists)
	return exists, err
}
