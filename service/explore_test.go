package service

import (
	"context"
	"fmt"
	"log"
	"muzz-project/service/protos"
	"muzz-project/storage"
	storageMock "muzz-project/storage/mocks"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestExploreService_listLikesHandler(t *testing.T) {
	arbitraryTime := time.Now()
	emptyString := ""

	ctx := context.Background()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStorage := storageMock.NewMockStorage(mockCtrl)

	tests := map[string]struct {
		mockStorageOutcomes func(storageMock *storageMock.MockStorage)
		maxPageSize         int
		in                  *protos.ListLikedYouRequest
		dbFunctionName      string
		want                *protos.ListLikedYouResponse
		wantErr             error
	}{
		"ListLikedYou returns results": {
			mockStorageOutcomes: func(storageMock *storageMock.MockStorage) {
				storageMock.EXPECT().GetLikesForUser(gomock.Any(), "1", 0).Times(1).Return([]*storage.Decision{
					{
						ID:          1,
						ActorID:     2,
						RecipientID: 1,
						Liked:       false,
						CreatedAt:   arbitraryTime,
					},
				}, nil)
			},
			maxPageSize: 10,
			in: &protos.ListLikedYouRequest{
				RecipientUserId: "1",
				PaginationToken: nil,
			},
			dbFunctionName: "GetLikesForUser",
			want: &protos.ListLikedYouResponse{
				Likers: []*protos.ListLikedYouResponse_Liker{
					{
						ActorId:       "2",
						UnixTimestamp: uint64(arbitraryTime.Unix()),
					},
				},
				NextPaginationToken: &emptyString,
			},
			wantErr: nil,
		},
		"ListNewLikedYou returns results": {
			mockStorageOutcomes: func(storageMock *storageMock.MockStorage) {
				storageMock.EXPECT().GetNewLikesForUser(gomock.Any(), "1", 0).Times(1).Return([]*storage.Decision{
					{
						ID:          1,
						ActorID:     2,
						RecipientID: 1,
						Liked:       false,
						CreatedAt:   arbitraryTime,
					},
				}, nil)
			},
			maxPageSize: 10,
			in: &protos.ListLikedYouRequest{
				RecipientUserId: "1",
				PaginationToken: nil,
			},
			dbFunctionName: "GetNewLikesForUser",
			want: &protos.ListLikedYouResponse{
				Likers: []*protos.ListLikedYouResponse_Liker{
					{
						ActorId:       "2",
						UnixTimestamp: uint64(arbitraryTime.Unix()),
					},
				},
				NextPaginationToken: &emptyString,
			},
			wantErr: nil,
		},
		"ListLikedYou returns results with pagination": {
			mockStorageOutcomes: func(storageMock *storageMock.MockStorage) {
				storageMock.EXPECT().GetLikesForUser(gomock.Any(), "1", 10).Times(1).Return([]*storage.Decision{
					{
						ID:          1,
						ActorID:     2,
						RecipientID: 1,
						Liked:       false,
						CreatedAt:   arbitraryTime,
					},
				}, nil)
			},
			maxPageSize: 10,
			in: &protos.ListLikedYouRequest{
				RecipientUserId: "1",
				PaginationToken: stringPtr("10"),
			},
			dbFunctionName: "GetLikesForUser",
			want: &protos.ListLikedYouResponse{
				Likers: []*protos.ListLikedYouResponse_Liker{
					{
						ActorId:       "2",
						UnixTimestamp: uint64(arbitraryTime.Unix()),
					},
				},
				NextPaginationToken: &emptyString,
			},
			wantErr: nil,
		},
		"ListLikedYou returns results with pagination error": {
			mockStorageOutcomes: func(storageMock *storageMock.MockStorage) {},
			maxPageSize:         10,
			in: &protos.ListLikedYouRequest{
				RecipientUserId: "1",
				PaginationToken: stringPtr("bad"),
			},
			dbFunctionName: "GetLikesForUser",
			want:           nil,
			wantErr:        badTokenError,
		},
		"ListLikedYou returns results with storage service error": {
			mockStorageOutcomes: func(storageMock *storageMock.MockStorage) {
				storageMock.EXPECT().GetLikesForUser(gomock.Any(), "1", 0).Times(1).Return(nil, fmt.Errorf("storage error"))
			},
			maxPageSize: 10,
			in: &protos.ListLikedYouRequest{
				RecipientUserId: "1",
				PaginationToken: nil,
			},
			dbFunctionName: "GetLikesForUser",
			want:           nil,
			wantErr:        fmt.Errorf("storage error"),
		},
		"ListLikedYou returns results with working pagination": {
			mockStorageOutcomes: func(storageMock *storageMock.MockStorage) {
				storageMock.EXPECT().GetLikesForUser(gomock.Any(), "1", 0).Times(1).
					Return([]*storage.Decision{
						{
							ID:          1,
							ActorID:     2,
							RecipientID: 1,
							Liked:       false,
							CreatedAt:   arbitraryTime,
						},
					}, nil)
			},
			maxPageSize: 1,
			in: &protos.ListLikedYouRequest{
				RecipientUserId: "1",
				PaginationToken: nil,
			},
			dbFunctionName: "GetLikesForUser",
			want: &protos.ListLikedYouResponse{
				Likers: []*protos.ListLikedYouResponse_Liker{
					{
						ActorId:       "2",
						UnixTimestamp: uint64(arbitraryTime.Unix()),
					},
				},
				NextPaginationToken: stringPtr("1"),
			},
			wantErr: nil,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tt.mockStorageOutcomes(mockStorage)

			e := ExploreService{
				storage:     mockStorage,
				maxPageSize: tt.maxPageSize,
			}

			var dbFunction func(context.Context, string, int) ([]*storage.Decision, error)

			if tt.dbFunctionName == "GetLikesForUser" {
				dbFunction = mockStorage.GetLikesForUser
			} else if tt.dbFunctionName == "GetNewLikesForUser" {
				dbFunction = mockStorage.GetNewLikesForUser
			} else {
				log.Fatal("bad function name")
			}

			got, err := e.listLikesHandler(ctx, tt.in, dbFunction)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func stringPtr(s string) *string {
	return &s
}
