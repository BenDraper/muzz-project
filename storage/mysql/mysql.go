package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"muzz-project/storage"
	"time"
)

type MysqlStorage struct {
	db          *sql.DB
	maxPageSize int
}

var _ storage.Storage = (*MysqlStorage)(nil)

func NewMysqlStorage(db *sql.DB, maxPageSize int) *MysqlStorage {
	return &MysqlStorage{
		db:          db,
		maxPageSize: maxPageSize,
	}
}

func (m *MysqlStorage) GetLikesForUser(ctx context.Context, userId string, paginationToken int) ([]*storage.Decision, error) {
	query := fmt.Sprintf("SELECT id, actor_id, recipient_id, liked, created_at FROM Decisions WHERE recipient_id = ? AND liked = TRUE ORDER BY created_at DESC LIMIT %d OFFSET %d", m.maxPageSize, paginationToken)
	return m.getLikesHandler(ctx, userId, query)
}

func (m *MysqlStorage) GetNewLikesForUser(ctx context.Context, userId string, paginationToken int) ([]*storage.Decision, error) {
	query := fmt.Sprintf("SELECT id, actor_id, recipient_id, liked, created_at FROM Decisions d1 WHERE d1.recipient_id = ? AND d1.liked = TRUE AND NOT EXISTS (SELECT 1 FROM Decisions d2 WHERE d2.actor_id = d1.recipient_id  AND d2.recipient_id = d1.actor_id) ORDER BY created_at DESC LIMIT %d OFFSET %d;", m.maxPageSize, paginationToken)
	return m.getLikesHandler(ctx, userId, query)
}

func (m *MysqlStorage) getLikesHandler(ctx context.Context, userId string, query string) ([]*storage.Decision, error) {
	var decisions []*storage.Decision

	rows, err := m.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var decision storage.Decision
		if err := rows.Scan(&decision.ID, &decision.ActorID, &decision.RecipientID, &decision.Liked, &decision.CreatedAt); err != nil {
			return nil, err
		}
		decisions = append(decisions, &decision)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return decisions, nil
}

func (m *MysqlStorage) GetLikesCountForUser(ctx context.Context, userId string) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM Decisions WHERE recipient_id = ? AND liked = TRUE`

	err := m.db.QueryRowContext(ctx, query, userId).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (m *MysqlStorage) AddDecision(ctx context.Context, actorId string, recipientId string, liked bool) (bool, error) {
	//Check to see if recipient has already liked actor - if so it's a match!
	reciprocal := false

	if liked {
		var count int64

		matchQuery := `SELECT COUNT(*) FROM Decisions WHERE recipient_id = ? AND actor_id = ? AND liked = TRUE`

		err := m.db.QueryRowContext(ctx, matchQuery, actorId, recipientId).Scan(&count)
		if err != nil {
			return false, err
		}

		reciprocal = count > 0

	}

	query := `INSERT INTO Decisions (actor_id, recipient_id, liked, created_at) VALUES (?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, query, actorId, recipientId, liked, time.Now())

	return reciprocal, err
}
