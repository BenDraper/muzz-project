package storage

import (
	"context"
	"fmt"
	"muzz-project/service/protos"
	"time"
)

type Storage interface {
	GetLikesForUser(ctx context.Context, userId string, paginationToken int) ([]*Decision, error)
	GetNewLikesForUser(ctx context.Context, userId string, paginationToken int) ([]*Decision, error)
	GetLikesCountForUser(ctx context.Context, userId string) (int64, error)
	AddDecision(ctx context.Context, actorId string, recipientId string, liked bool) (bool, error)
}

type Decision struct {
	ID          int64     `db:"id"`
	ActorID     int64     `db:"actor_id"`
	RecipientID int64     `db:"recipient_id"`
	Liked       bool      `db:"liked"`
	CreatedAt   time.Time `db:"created_at"`
}

func (d Decision) ToProto() *protos.ListLikedYouResponse_Liker {
	return &protos.ListLikedYouResponse_Liker{
		ActorId:       fmt.Sprintf("%d", d.ActorID),
		UnixTimestamp: uint64(d.CreatedAt.Unix()),
	}
}
