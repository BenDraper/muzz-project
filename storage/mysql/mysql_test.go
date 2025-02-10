package mysql

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/assert"
	"muzz-project/storage"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestMysqlStorage_GetLikesCountForUser(t *testing.T) {
	ctx := context.Background()

	tests := map[string]struct {
		dbOutcomes func(mock sqlmock.Sqlmock)
		userId     string
		want       int64
		wantErr    error
	}{
		"user with likes": {
			dbOutcomes: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM Decisions WHERE recipient_id = \? AND liked = TRUE`).
					WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(5))
			},
			userId:  "1",
			want:    5,
			wantErr: nil,
		},
		"user without likes": {
			dbOutcomes: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM Decisions WHERE recipient_id = \? AND liked = TRUE`).
					WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))
			},
			userId:  "1",
			want:    0,
			wantErr: nil,
		},
		"database error": {
			dbOutcomes: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM Decisions WHERE recipient_id = \? AND liked = TRUE`).
					WithArgs("2").
					WillReturnError(sql.ErrConnDone)
			},
			userId:  "2",
			want:    0,
			wantErr: sql.ErrConnDone,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer mockDB.Close()

			tt.dbOutcomes(mock)

			m := &MysqlStorage{
				db: mockDB,
			}
			got, err := m.GetLikesCountForUser(ctx, tt.userId)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestMysqlStorage_GetLikesForUser(t *testing.T) {
	ctx := context.Background()
	arbitraryTime := time.Now()

	tests := map[string]struct {
		dbOutcomes  func(mock sqlmock.Sqlmock)
		userId      string
		query       string
		maxPageSize int
		want        []*storage.Decision
		wantErr     error
	}{
		"user with likes": {
			dbOutcomes: func(mock sqlmock.Sqlmock) {
				query := `SELECT id, actor_id, recipient_id, liked, created_at FROM Decisions WHERE recipient_id = \? AND liked = TRUE ORDER BY created_at DESC LIMIT 10 OFFSET 0`
				rows := sqlmock.NewRows([]string{"id", "actor_id", "recipient_id", "liked", "created_at"}).
					AddRow("1", "2", "1", true, arbitraryTime).
					AddRow("2", "3", "1", true, arbitraryTime)
				mock.ExpectQuery(query).WithArgs("1").WillReturnRows(rows)

			},
			userId:      "1",
			query:       `SELECT id, actor_id, recipient_id, liked, created_at FROM Decisions WHERE recipient_id = ? AND liked = TRUE ORDER BY created_at DESC LIMIT 10 OFFSET 0`,
			maxPageSize: 10,
			want: []*storage.Decision{
				{ID: 1, ActorID: 2, RecipientID: 1, Liked: true, CreatedAt: arbitraryTime},
				{ID: 2, ActorID: 3, RecipientID: 1, Liked: true, CreatedAt: arbitraryTime},
			},
			wantErr: nil,
		},
		"user without likes": {
			dbOutcomes: func(mock sqlmock.Sqlmock) {
				query := `SELECT id, actor_id, recipient_id, liked, created_at FROM Decisions WHERE recipient_id = \? AND liked = TRUE ORDER BY created_at DESC LIMIT 10 OFFSET 0`
				rows := sqlmock.NewRows([]string{"id", "actor_id", "recipient_id", "liked", "created_at"})
				mock.ExpectQuery(query).WithArgs("1").WillReturnRows(rows)

			},
			userId:      "1",
			query:       `SELECT id, actor_id, recipient_id, liked, created_at FROM Decisions WHERE recipient_id = ? AND liked = TRUE ORDER BY created_at DESC LIMIT 10 OFFSET 0`,
			maxPageSize: 10,
			want:        nil,
			wantErr:     nil,
		},
		"database error": {
			dbOutcomes: func(mock sqlmock.Sqlmock) {
				query := `SELECT id, actor_id, recipient_id, liked, created_at FROM Decisions WHERE recipient_id = \? AND liked = TRUE ORDER BY created_at DESC LIMIT 10 OFFSET 0`
				mock.ExpectQuery(query).
					WithArgs("2").
					WillReturnError(sql.ErrConnDone)
			},
			userId:      "2",
			query:       `SELECT id, actor_id, recipient_id, liked, created_at FROM Decisions WHERE recipient_id = ? AND liked = TRUE ORDER BY created_at DESC LIMIT 10 OFFSET 0`,
			maxPageSize: 10,
			want:        nil,
			wantErr:     sql.ErrConnDone,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer mockDB.Close()

			tt.dbOutcomes(mock)

			m := &MysqlStorage{
				db:          mockDB,
				maxPageSize: tt.maxPageSize,
			}

			got, err := m.getLikesHandler(ctx, tt.userId, tt.query)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestMysqlStorage_AddDecision(t *testing.T) {
	ctx := context.Background()

	tests := map[string]struct {
		dbOutcomes  func(mock sqlmock.Sqlmock)
		actorId     string
		recipientId string
		liked       bool
		want        bool
		wantErr     error
	}{
		"users match": {
			dbOutcomes: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM Decisions WHERE recipient_id = \? AND actor_id = \? AND liked = TRUE`).
					WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO Decisions (actor_id, recipient_id, liked, created_at) VALUES (?, ?, ?, ?)")).
					WithArgs("1", "2", true, sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(0, 1))

			},
			actorId:     "1",
			recipientId: "2",
			liked:       true,
			want:        true,
			wantErr:     nil,
		},
		"users don't match": {
			dbOutcomes: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM Decisions WHERE recipient_id = \? AND actor_id = \? AND liked = TRUE`).
					WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0))
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO Decisions (actor_id, recipient_id, liked, created_at) VALUES (?, ?, ?, ?)")).
					WithArgs("1", "2", true, sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(0, 1))

			},
			actorId:     "1",
			recipientId: "2",
			liked:       true,
			want:        false,
			wantErr:     nil,
		},
		"actor doesn't like recipient": {
			dbOutcomes: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO Decisions (actor_id, recipient_id, liked, created_at) VALUES (?, ?, ?, ?)")).
					WithArgs("1", "2", false, sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(0, 1))

			},
			actorId:     "1",
			recipientId: "2",
			liked:       false,
			want:        false,
			wantErr:     nil,
		},
		"database error 1": {
			dbOutcomes: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM Decisions WHERE recipient_id = \? AND actor_id = \? AND liked = TRUE`).
					WillReturnError(sql.ErrConnDone)
			},
			actorId:     "1",
			recipientId: "2",
			liked:       true,
			want:        false,
			wantErr:     sql.ErrConnDone,
		},
		"database error 2": {
			dbOutcomes: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM Decisions WHERE recipient_id = \? AND actor_id = \? AND liked = TRUE`).
					WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO Decisions (actor_id, recipient_id, liked, created_at) VALUES (?, ?, ?, ?)")).
					WithArgs("1", "2", true, sqlmock.AnyArg()).
					WillReturnError(sql.ErrConnDone)

			},
			actorId:     "1",
			recipientId: "2",
			liked:       true,
			want:        true,
			wantErr:     sql.ErrConnDone,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer mockDB.Close()

			tt.dbOutcomes(mock)

			m := &MysqlStorage{
				db: mockDB,
			}

			got, err := m.AddDecision(ctx, tt.actorId, tt.recipientId, tt.liked)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
