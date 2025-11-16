package repository

import (
	"database/sql"
	"pull_requests_service/internal/domain"
)

type prRepository struct {
	BaseRepository
}

func NewPRRepository(db *sql.DB) domain.PRRepository {
	return &prRepository{BaseRepository{db: db}}
}

func (r *prRepository) CreatePR(pr *domain.PullRequest) error {
	var reviewer1, reviewer2 sql.NullString
	if len(pr.AssignedReviewers) > 0 {
		reviewer1 = sql.NullString{String: pr.AssignedReviewers[0], Valid: true}
	}
	if len(pr.AssignedReviewers) > 1 {
		reviewer2 = sql.NullString{String: pr.AssignedReviewers[1], Valid: true}
	}

	_, err := r.db.Exec(`
        INSERT INTO pull_request (pull_request_id, pull_request_name, author_id, status, reviewer_1, reviewer_2)
        VALUES ($1, $2, $3, $4, $5, $6)`,
		pr.PullRequestID, pr.PullRequestName, pr.AuthorID, pr.Status, reviewer1, reviewer2,
	)
	return err
}

func (r *prRepository) GetPR(prID string) (*domain.PullRequest, error) {
	var pr domain.PullRequest
	var reviewer1, reviewer2 sql.NullString
	var mergedAt sql.NullTime

	err := r.db.QueryRow(`
        SELECT pull_request_id, pull_request_name, author_id, status, reviewer_1, reviewer_2, "mergedAt"
        FROM pull_request WHERE pull_request_id = $1`,
		prID,
	).Scan(&pr.PullRequestID, &pr.PullRequestName, &pr.AuthorID, &pr.Status, &reviewer1, &reviewer2, &mergedAt)

	if err == sql.ErrNoRows {
		return nil, domain.ErrPRNotFound
	}
	if err != nil {
		return nil, err
	}

	if reviewer1.Valid {
		pr.AssignedReviewers = append(pr.AssignedReviewers, reviewer1.String)
	}
	if reviewer2.Valid {
		pr.AssignedReviewers = append(pr.AssignedReviewers, reviewer2.String)
	}

	if mergedAt.Valid {
		pr.MergedAt = &mergedAt.Time
	}

	return &pr, nil
}

func (r *prRepository) UpdatePR(pr *domain.PullRequest) error {
	var reviewer1, reviewer2 sql.NullString
	if len(pr.AssignedReviewers) > 0 {
		reviewer1 = sql.NullString{String: pr.AssignedReviewers[0], Valid: true}
	}
	if len(pr.AssignedReviewers) > 1 {
		reviewer2 = sql.NullString{String: pr.AssignedReviewers[1], Valid: true}
	}

	_, err := r.db.Exec(`
        UPDATE pull_request 
        SET pull_request_name = $1, status = $2, reviewer_1 = $3, reviewer_2 = $4, "mergedAt" = $5
        WHERE pull_request_id = $6`,
		pr.PullRequestName, pr.Status, reviewer1, reviewer2, pr.MergedAt, pr.PullRequestID,
	)
	return err
}

func (r *prRepository) GetPRsByReviewer(userID string) ([]*domain.PullRequest, error) {
	rows, err := r.db.Query(`
        SELECT pull_request_id, pull_request_name, author_id, status, reviewer_1, reviewer_2, "mergedAt"
        FROM pull_request 
        WHERE reviewer_1 = $1 OR reviewer_2 = $1`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prs []*domain.PullRequest
	for rows.Next() {
		var pr domain.PullRequest
		var reviewer1, reviewer2 sql.NullString
		var mergedAt sql.NullTime

		if err := rows.Scan(&pr.PullRequestID, &pr.PullRequestName, &pr.AuthorID, &pr.Status, &reviewer1, &reviewer2, &mergedAt); err != nil {
			return nil, err
		}

		if reviewer1.Valid {
			pr.AssignedReviewers = append(pr.AssignedReviewers, reviewer1.String)
		}
		if reviewer2.Valid {
			pr.AssignedReviewers = append(pr.AssignedReviewers, reviewer2.String)
		}
		if mergedAt.Valid {
			pr.MergedAt = &mergedAt.Time
		}

		prs = append(prs, &pr)
	}

	return prs, nil
}

func (r *prRepository) PRExists(prID string) (bool, error) {
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM pull_request WHERE pull_request_id = $1)", prID).Scan(&exists)
	return exists, err
}
