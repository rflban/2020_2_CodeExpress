package grpc_track_test

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/track/proto_track"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/track/grpc_track"
	"golang.org/x/net/context"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/track/mock_track"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
)

func TestTrackUsecase_CreateTrack(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_track.NewMockTrackRep(ctrl)
	mockGRPC := grpc_track.NewTrackGRPCUsecase(mockRepo)

	id := uint64(42)
	track := &models.Track{}

	mockRepo.
		EXPECT().
		Insert(gomock.Eq(track)).
		DoAndReturn(func(track *models.Track) error {
			track.ID = id
			return nil
		})

	_, err := mockGRPC.CreateTrack(context.Background(), grpc_track.TrackToTrackGRPC(track))
	assert.Equal(t, err, nil)
}

func TestTrackUsecase_CreateTrack_Failed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_track.NewMockTrackRep(ctrl)
	mockGRPC := grpc_track.NewTrackGRPCUsecase(mockRepo)

	track := &models.Track{}
	dbErr := errors.New("Some db error")

	mockRepo.
		EXPECT().
		Insert(gomock.Eq(track)).
		Return(dbErr)

	_, err := mockGRPC.CreateTrack(context.Background(), grpc_track.TrackToTrackGRPC(track))
	assert.Equal(t, err, dbErr)
}

func TestArtistUsecase_DeleteArtist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_track.NewMockTrackRep(ctrl)
	mockGRPC := grpc_track.NewTrackGRPCUsecase(mockRepo)

	id := uint64(5)
	mockRepo.
		EXPECT().
		Delete(gomock.Eq(id)).
		Return(nil)

	_, err := mockGRPC.DeleteTrack(context.Background(), &proto_track.TrackID{
		ID: id,
	})
	assert.Equal(t, err, nil)
}

func TestArtistUsecase_DeleteArtist_Failed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_track.NewMockTrackRep(ctrl)
	mockGRPC := grpc_track.NewTrackGRPCUsecase(mockRepo)

	id := uint64(5)

	mockRepo.
		EXPECT().
		Delete(gomock.Eq(id)).
		Return(sql.ErrNoRows)

	_, err := mockGRPC.DeleteTrack(context.Background(), &proto_track.TrackID{
		ID: id,
	})
	assert.Equal(t, err, sql.ErrNoRows)
}

func TestArtistUsecase_DeleteArtist_Internal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_track.NewMockTrackRep(ctrl)
	mockGRPC := grpc_track.NewTrackGRPCUsecase(mockRepo)

	id := uint64(5)
	dbErr := errors.New("Some db error")

	mockRepo.
		EXPECT().
		Delete(gomock.Eq(id)).
		Return(dbErr)

	_, err := mockGRPC.DeleteTrack(context.Background(), &proto_track.TrackID{
		ID: id,
	})
	assert.Equal(t, err, dbErr)
}

func TestArtistUsecase_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_track.NewMockTrackRep(ctrl)
	mockGRPC := grpc_track.NewTrackGRPCUsecase(mockRepo)

	expectedTrack := &models.Track{
		ID: 1,
	}

	mockRepo.
		EXPECT().
		SelectByID(gomock.Eq(expectedTrack.ID), uint64(0)).
		Return(expectedTrack, nil)

	_, err := mockGRPC.GetByID(context.Background(), &proto_track.GetByIdMessage{
		TrackId: expectedTrack.ID,
	})
	assert.Equal(t, err, nil)
}

func TestArtistUsecase_GetByID_Failed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_track.NewMockTrackRep(ctrl)
	mockGRPC := grpc_track.NewTrackGRPCUsecase(mockRepo)

	expectedTrack := &models.Track{
		ID: 1,
	}

	mockRepo.
		EXPECT().
		SelectByID(gomock.Eq(expectedTrack.ID), gomock.Eq(uint64(0))).
		Return(nil, sql.ErrNoRows)

	_, err := mockGRPC.GetByID(context.Background(), &proto_track.GetByIdMessage{
		TrackId: expectedTrack.ID,
	})
	assert.Equal(t, err, sql.ErrNoRows)
}

func TestArtistUsecase_GetByID_Internal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_track.NewMockTrackRep(ctrl)
	mockGRPC := grpc_track.NewTrackGRPCUsecase(mockRepo)

	dbErr := errors.New("Some db error")
	expectedTrack := &models.Track{
		ID: 1,
	}

	mockRepo.
		EXPECT().
		SelectByID(gomock.Eq(expectedTrack.ID), uint64(0)).
		Return(nil, dbErr)

	_, err := mockGRPC.GetByID(context.Background(), &proto_track.GetByIdMessage{
		TrackId: expectedTrack.ID,
	})
	assert.Equal(t, err, dbErr)
}

func TestArtistUsecase_GetByArtistID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_track.NewMockTrackRep(ctrl)
	mockGRPC := grpc_track.NewTrackGRPCUsecase(mockRepo)

	expectedTrack := &models.Track{
		ID:       1,
		ArtistID: 5,
	}

	expectedTracks := []*models.Track{expectedTrack}

	mockRepo.
		EXPECT().
		SelectByArtistId(gomock.Eq(expectedTrack.ArtistID), uint64(0)).
		Return(expectedTracks, nil)

	_, err := mockGRPC.GetByArtistId(context.Background(), &proto_track.GetByArtistIdMessage{
		UserID:   uint64(0),
		ArtistID: expectedTrack.ArtistID,
	})
	assert.Equal(t, err, nil)
}

func TestArtistUsecase_GetByArtistID_Failed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_track.NewMockTrackRep(ctrl)
	mockGRPC := grpc_track.NewTrackGRPCUsecase(mockRepo)

	expectedTrack := &models.Track{
		ID: 1,
	}

	mockRepo.
		EXPECT().
		SelectByID(gomock.Eq(expectedTrack.ID), gomock.Eq(uint64(0))).
		Return(nil, sql.ErrNoRows)

	_, err := mockGRPC.GetByID(context.Background(), &proto_track.GetByIdMessage{
		TrackId: expectedTrack.ID,
	})
	assert.Equal(t, err, sql.ErrNoRows)
}

func TestArtistUsecase_GetByArtistID_Internal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_track.NewMockTrackRep(ctrl)
	mockGRPC := grpc_track.NewTrackGRPCUsecase(mockRepo)

	dbErr := errors.New("Some db error")
	expectedTrack := &models.Track{
		ID: 1,
	}

	mockRepo.
		EXPECT().
		SelectByID(gomock.Eq(expectedTrack.ID), uint64(0)).
		Return(nil, dbErr)

	_, err := mockGRPC.GetByID(context.Background(), &proto_track.GetByIdMessage{
		TrackId: expectedTrack.ID,
	})
	assert.Equal(t, err, dbErr)
}

func TestArtistUsecase_GetByArtistParams(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_track.NewMockTrackRep(ctrl)
	mockGRPC := grpc_track.NewTrackGRPCUsecase(mockRepo)

	count, from := uint64(1), uint64(0)
	expectedTrack := &models.Track{
		ID:       1,
		ArtistID: 5,
	}

	expectedTracks := []*models.Track{expectedTrack}

	mockRepo.
		EXPECT().
		SelectByParams(gomock.Eq(count), gomock.Eq(from), uint64(0)).
		Return(expectedTracks, nil)

	_, err := mockGRPC.GetByParams(context.Background(), &proto_track.GetByParamsMessage{
		Count:  count,
		From:   from,
		UserID: uint64(0),
	})
	assert.Equal(t, err, nil)
}

func TestArtistUsecase_GetByArtistParams_Internal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_track.NewMockTrackRep(ctrl)
	mockGRPC := grpc_track.NewTrackGRPCUsecase(mockRepo)

	dbErr := errors.New("Some db error")
	count, from := uint64(1), uint64(0)

	mockRepo.
		EXPECT().
		SelectByParams(gomock.Eq(count), gomock.Eq(from), uint64(0)).
		Return(nil, dbErr)

	_, err := mockGRPC.GetByParams(context.Background(), &proto_track.GetByParamsMessage{
		Count:  count,
		From:   from,
		UserID: uint64(0),
	})
	assert.Equal(t, err, dbErr)
}

func TestTrackUsecase_UpdateTrack(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_track.NewMockTrackRep(ctrl)
	mockGRPC := grpc_track.NewTrackGRPCUsecase(mockRepo)

	id := uint64(42)
	track := &models.Track{ID: id}

	mockRepo.
		EXPECT().
		Update(gomock.Eq(track)).
		Return(nil)

	_, err := mockGRPC.UpdateTrack(context.Background(), grpc_track.TrackToTrackGRPC(track))
	assert.Equal(t, err, nil)
}

func TestTrackUsecase_UpdateTrackAudio(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_track.NewMockTrackRep(ctrl)
	mockGRPC := grpc_track.NewTrackGRPCUsecase(mockRepo)

	id := uint64(42)
	track := &models.Track{ID: id}

	mockRepo.
		EXPECT().
		UpdateAudio(gomock.Eq(track)).
		Return(nil)

	_, err := mockGRPC.UpdateTrackAudio(context.Background(), grpc_track.TrackToTrackGRPC(track))
	assert.Equal(t, err, nil)
}

func TestArtistUsecase_GetByAlbumID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_track.NewMockTrackRep(ctrl)
	mockGRPC := grpc_track.NewTrackGRPCUsecase(mockRepo)

	albumID := uint64(42)
	expectedTrack := &models.Track{
		AlbumID: albumID,
	}

	expectedTracks := []*models.Track{expectedTrack}

	mockRepo.
		EXPECT().
		SelectByAlbumID(albumID, uint64(0)).
		Return(expectedTracks, nil)

	_, err := mockGRPC.GetByAlbumID(context.Background(), &proto_track.GetByAlbumIdMessage{
		AlbumId: albumID,
	})
	assert.Equal(t, err, nil)
}

func TestArtistUsecase_GetFavoritesByUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_track.NewMockTrackRep(ctrl)
	mockGRPC := grpc_track.NewTrackGRPCUsecase(mockRepo)

	userID := uint64(42)
	expectedTrack := &models.Track{
		ID: 5,
	}

	expectedTracks := []*models.Track{expectedTrack}

	mockRepo.
		EXPECT().
		SelectFavoritesByUserID(userID).
		Return(expectedTracks, nil)

	_, err := mockGRPC.GetFavoritesByUserID(context.Background(), &proto_track.UserID{
		ID: userID,
	})
	assert.Equal(t, err, nil)
}

func TestArtistUsecase_AddToFavorites(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_track.NewMockTrackRep(ctrl)
	mockGRPC := grpc_track.NewTrackGRPCUsecase(mockRepo)

	userID := uint64(42)
	trackID := uint64(7)

	mockRepo.
		EXPECT().
		InsertTrackToUser(gomock.Eq(userID), gomock.Eq(trackID)).
		Return(nil)

	_, err := mockGRPC.AddToFavourites(context.Background(), &proto_track.Favorites{
		UserID:  userID,
		TrackID: trackID,
	})
	assert.Equal(t, err, nil)
}

func TestArtistUsecase_DeleteFromFavourites(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_track.NewMockTrackRep(ctrl)
	mockGRPC := grpc_track.NewTrackGRPCUsecase(mockRepo)

	userID := uint64(42)
	trackID := uint64(7)

	mockRepo.
		EXPECT().
		DeleteTrackFromUsersTracks(gomock.Eq(userID), gomock.Eq(trackID)).
		Return(nil)

	_, err := mockGRPC.DeleteFromFavourites(context.Background(), &proto_track.Favorites{
		UserID:  userID,
		TrackID: trackID,
	})
	assert.Equal(t, err, nil)
}
