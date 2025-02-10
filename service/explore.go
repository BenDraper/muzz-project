package service

import (
	"context"
	"fmt"
	"muzz-project/service/protos"
	"muzz-project/storage"
	"strconv"
)

var (
	badTokenError = fmt.Errorf("Token must be positive integer")
)

type ExploreService struct {
	storage     storage.Storage
	maxPageSize int
}

func NewExploreService(storage storage.Storage, maxPageSize int) *ExploreService {
	return &ExploreService{
		storage:     storage,
		maxPageSize: maxPageSize,
	}
}

func (e ExploreService) ListLikedYou(ctx context.Context, in *protos.ListLikedYouRequest) (*protos.ListLikedYouResponse, error) {
	return e.listLikesHandler(ctx, in, e.storage.GetLikesForUser)
}
func (e ExploreService) ListNewLikedYou(ctx context.Context, in *protos.ListLikedYouRequest) (*protos.ListLikedYouResponse, error) {
	return e.listLikesHandler(ctx, in, e.storage.GetNewLikesForUser)
}

func (e ExploreService) listLikesHandler(ctx context.Context, in *protos.ListLikedYouRequest, dbFunction func(context.Context, string, int) ([]*storage.Decision, error)) (*protos.ListLikedYouResponse, error) {
	var token = 0
	var err error

	if in.GetPaginationToken() != "" {
		token, err = strconv.Atoi(in.GetPaginationToken())
		if err != nil {
			return nil, badTokenError
		}
	}
	likes, err := dbFunction(ctx, in.GetRecipientUserId(), token)
	if err != nil {
		return nil, err
	}

	nextPaginationToken := ""

	if len(likes) == e.maxPageSize {
		nextPaginationToken = fmt.Sprintf("%d", token+e.maxPageSize)
	}

	out := &protos.ListLikedYouResponse{
		Likers:              []*protos.ListLikedYouResponse_Liker{},
		NextPaginationToken: &nextPaginationToken,
	}
	for _, l := range likes {
		out.Likers = append(out.Likers, l.ToProto())
	}
	return out, nil
}

func (e ExploreService) CountLikedYou(ctx context.Context, in *protos.CountLikedYouRequest) (*protos.CountLikedYouResponse, error) {

	count, err := e.storage.GetLikesCountForUser(ctx, in.GetRecipientUserId())
	if err != nil {
		return nil, err
	}

	return &protos.CountLikedYouResponse{
		Count: uint64(count),
	}, nil
}
func (e ExploreService) PutDecision(ctx context.Context, in *protos.PutDecisionRequest) (*protos.PutDecisionResponse, error) {
	match, err := e.storage.AddDecision(ctx, in.GetActorUserId(), in.GetRecipientUserId(), in.GetLikedRecipient())
	if err != nil {
		return nil, err
	}
	return &protos.PutDecisionResponse{
		MutualLikes: match,
	}, nil
}
