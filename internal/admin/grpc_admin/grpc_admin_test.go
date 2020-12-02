package grpc_admin_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/admin/grpc_admin"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/artist/mock_artist"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGRPCAdmin_Create_Passed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_artist.NewMockArtistRep(ctrl)
	mockGRPC := grpc_admin.NewAdminGRPCUsecase(mockRepo)

	name := "Some name"
	artist := &models.Artist{
		Name: name,
	}

	mockRepo.
		EXPECT().
		Insert(gomock.Eq(artist)).
		Return(nil)

	_, err := mockGRPC.CreateArtist(context.Background(), grpc_admin.ArtistToGRPCArtist(artist))
	assert.Nil(t, err)
}

func TestGRPCAdmin_Create_Failed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_artist.NewMockArtistRep(ctrl)
	mockGRPC := grpc_admin.NewAdminGRPCUsecase(mockRepo)

	name := "Some name"
	artist := &models.Artist{
		Name: name,
	}

	mockRepo.
		EXPECT().
		Insert(gomock.Eq(artist)).
		Return(sql.ErrNoRows)

	_, err := mockGRPC.CreateArtist(context.Background(), grpc_admin.ArtistToGRPCArtist(artist))
	assert.Equal(t, err, sql.ErrNoRows)
}
