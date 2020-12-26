package usecase

import (
	"context"
	"testing"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/track/grpc_track"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/track/mock_track"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/track/proto_track"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
)

func TestTrackUsecase_CreateTrack(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_track.NewMockTrackServiceClient(ctrl)
	mockUsecase := NewTrackUsecase(mockClient)

	track := models.Track{
		ID:      1,
		Title:   "title",
		AlbumID: 1,
	}

	mockClient.
		EXPECT().
		CreateTrack(context.Background(), grpc_track.TrackToTrackGRPC(&track)).
		Return(grpc_track.TrackToTrackGRPC(&track), nil)

	argTrack := models.Track{
		Title:   "title",
		AlbumID: 1,
	}

	mockClient.
		EXPECT().
		GetByID(context.Background(), &proto_track.GetByIdMessage{
			TrackId: track.ID,
			UserId:  0,
		}).
		Return(grpc_track.TrackToTrackGRPC(&argTrack), nil)

	err := mockUsecase.CreateTrack(&track, 0)
	assert.Equal(t, err, nil)
}

func TestTrackUsecase_DeleteTrack(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_track.NewMockTrackServiceClient(ctrl)
	mockUsecase := NewTrackUsecase(mockClient)

	track := models.Track{
		ID:      1,
		Title:   "title",
		AlbumID: 1,
	}

	nothing := new(proto_track.Nothing)

	mockClient.
		EXPECT().
		DeleteTrack(context.Background(), &proto_track.TrackID{ID: track.ID}).
		Return(nothing, nil)

	err := mockUsecase.DeleteTrack(track.ID)
	assert.Equal(t, err, nil)
}

func TestTrackUsecase_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_track.NewMockTrackServiceClient(ctrl)
	mockUsecase := NewTrackUsecase(mockClient)

	track := models.Track{
		ID:      1,
		Title:   "title",
		AlbumID: 1,
	}

	argTrack := models.Track{
		ID:      1,
		Title:   "title",
		AlbumID: 1,
	}

	mockClient.
		EXPECT().
		GetByID(context.Background(), &proto_track.GetByIdMessage{
			TrackId: argTrack.ID,
			UserId:  0,
		}).
		Return(grpc_track.TrackToTrackGRPC(&argTrack), nil)

	newTrack, err := mockUsecase.GetByID(argTrack.ID, 0)
	assert.Equal(t, err, nil)
	assert.Equal(t, newTrack, track)
}

func TestTrackUsecase_GetByArtistId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_track.NewMockTrackServiceClient(ctrl)
	mockUsecase := NewTrackUsecase(mockClient)

	track := models.Track{
		ID:       1,
		Title:    "title",
		AlbumID:  1,
		ArtistID: 1,
	}

	userId := uint64(1)

	grpcTracks := &proto_track.Tracks{}
	grpcTracks.Tracks = append(grpcTracks.Tracks, grpc_track.TrackToTrackGRPC(&track))

	mockClient.
		EXPECT().
		GetByArtistId(context.Background(), &proto_track.GetByArtistIdMessage{
			ArtistID: track.ArtistID,
			UserID:   userId,
		}).
		Return(grpcTracks, nil)

	_, err := mockUsecase.GetByArtistId(track.ArtistID, userId)
	assert.Equal(t, err, nil)
}

func TestTrackUsecase_GetByParams(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_track.NewMockTrackServiceClient(ctrl)
	mockUsecase := NewTrackUsecase(mockClient)

	track := models.Track{
		ID:       1,
		Title:    "title",
		AlbumID:  1,
		ArtistID: 1,
	}

	count := uint64(5)
	from := uint64(1)
	userId := uint64(1)

	grpcTracks := &proto_track.Tracks{}
	grpcTracks.Tracks = append(grpcTracks.Tracks, grpc_track.TrackToTrackGRPC(&track))

	mockClient.
		EXPECT().
		GetByParams(context.Background(), &proto_track.GetByParamsMessage{
			Count:  count,
			From:   from,
			UserID: userId,
		}).
		Return(grpcTracks, nil)

	_, err := mockUsecase.GetByParams(count, from, userId)
	assert.Equal(t, err, nil)
}

func TestTrackUsecase_GetTopByParams(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_track.NewMockTrackServiceClient(ctrl)
	mockUsecase := NewTrackUsecase(mockClient)

	track := models.Track{
		ID:       1,
		Title:    "title",
		AlbumID:  1,
		ArtistID: 1,
	}

	count := uint64(5)
	from := uint64(1)
	userId := uint64(1)

	grpcTracks := &proto_track.Tracks{}
	grpcTracks.Tracks = append(grpcTracks.Tracks, grpc_track.TrackToTrackGRPC(&track))

	mockClient.
		EXPECT().
		GetTopByParams(context.Background(), &proto_track.GetTopByParamsMessage{
			Count:  count,
			From:   from,
			UserID: userId,
		}).
		Return(grpcTracks, nil)

	_, err := mockUsecase.GetTopByParams(count, from, userId)
	assert.Equal(t, err, nil)
}

func TestTrackUsecase_UpdateTrack(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_track.NewMockTrackServiceClient(ctrl)
	mockUsecase := NewTrackUsecase(mockClient)

	track := models.Track{
		ID:      1,
		Title:   "title",
		AlbumID: 1,
	}

	nothing := new(proto_track.Nothing)

	mockClient.
		EXPECT().
		UpdateTrack(context.Background(), grpc_track.TrackToTrackGRPC(&track)).
		Return(nothing, nil)

	argTrack := models.Track{
		Title:   "title",
		AlbumID: 1,
	}

	mockClient.
		EXPECT().
		GetByID(context.Background(), &proto_track.GetByIdMessage{
			TrackId: track.ID,
			UserId:  0,
		}).
		Return(grpc_track.TrackToTrackGRPC(&argTrack), nil)

	err := mockUsecase.UpdateTrack(&track, 0)
	assert.Equal(t, err, nil)
}

func TestTrackUsecase_UpdateTrackAudio(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_track.NewMockTrackServiceClient(ctrl)
	mockUsecase := NewTrackUsecase(mockClient)

	track := models.Track{
		ID:      1,
		Title:   "title",
		AlbumID: 1,
	}

	nothing := new(proto_track.Nothing)

	mockClient.
		EXPECT().
		UpdateTrackAudio(context.Background(), grpc_track.TrackToTrackGRPC(&track)).
		Return(nothing, nil)

	argTrack := models.Track{
		Title:   "title",
		AlbumID: 1,
	}

	mockClient.
		EXPECT().
		GetByID(context.Background(), &proto_track.GetByIdMessage{
			TrackId: track.ID,
			UserId:  0,
		}).
		Return(grpc_track.TrackToTrackGRPC(&argTrack), nil)

	err := mockUsecase.UpdateTrackAudio(&track, 0)
	assert.Equal(t, err, nil)
}

func TestTrackUsecase_GetByAlbumID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_track.NewMockTrackServiceClient(ctrl)
	mockUsecase := NewTrackUsecase(mockClient)

	track := models.Track{
		ID:       1,
		Title:    "title",
		AlbumID:  1,
		ArtistID: 1,
	}

	albumId := uint64(1)

	grpcTracks := &proto_track.Tracks{}
	grpcTracks.Tracks = append(grpcTracks.Tracks, grpc_track.TrackToTrackGRPC(&track))

	mockClient.
		EXPECT().
		GetByAlbumID(context.Background(), &proto_track.GetByAlbumIdMessage{
			AlbumId: albumId,
			UserId:  0,
		}).
		Return(grpcTracks, nil)

	_, err := mockUsecase.GetByAlbumID(albumId, 0)
	assert.Equal(t, err, nil)
}

func TestTrackUsecase_GetFavoritesByUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_track.NewMockTrackServiceClient(ctrl)
	mockUsecase := NewTrackUsecase(mockClient)

	track := models.Track{
		ID:       1,
		Title:    "title",
		AlbumID:  1,
		ArtistID: 1,
	}

	userID := uint64(1)

	grpcTracks := &proto_track.Tracks{}
	grpcTracks.Tracks = append(grpcTracks.Tracks, grpc_track.TrackToTrackGRPC(&track))

	mockClient.
		EXPECT().
		GetFavoritesByUserID(context.Background(), &proto_track.UserID{
			ID: userID,
		}).
		Return(grpcTracks, nil)

	_, err := mockUsecase.GetFavoritesByUserID(userID)
	assert.Equal(t, err, nil)
}

func TestTrackUsecase_AddToFavourites(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_track.NewMockTrackServiceClient(ctrl)
	mockUsecase := NewTrackUsecase(mockClient)

	track := models.Track{
		ID:      1,
		Title:   "title",
		AlbumID: 1,
	}

	argTrack := models.Track{
		Title:   "title",
		AlbumID: 1,
	}

	userID := uint64(1)

	nothing := new(proto_track.Nothing)

	mockClient.
		EXPECT().
		GetByID(context.Background(), &proto_track.GetByIdMessage{
			TrackId: track.ID,
			UserId:  userID,
		}).
		Return(grpc_track.TrackToTrackGRPC(&argTrack), nil)

	mockClient.
		EXPECT().
		AddToFavourites(context.Background(), &proto_track.Favorites{
			UserID:  userID,
			TrackID: track.ID,
		}).
		Return(nothing, nil)

	err := mockUsecase.AddToFavourites(userID, track.ID)
	assert.Equal(t, err, nil)
}

func TestTrackUsecase_DeleteFromFavourites(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_track.NewMockTrackServiceClient(ctrl)
	mockUsecase := NewTrackUsecase(mockClient)

	track := models.Track{
		ID:      1,
		Title:   "title",
		AlbumID: 1,
	}

	argTrack := models.Track{
		Title:   "title",
		AlbumID: 1,
	}

	userID := uint64(1)

	nothing := new(proto_track.Nothing)

	mockClient.
		EXPECT().
		GetByID(context.Background(), &proto_track.GetByIdMessage{
			TrackId: track.ID,
			UserId:  userID,
		}).
		Return(grpc_track.TrackToTrackGRPC(&argTrack), nil)

	mockClient.
		EXPECT().
		DeleteFromFavourites(context.Background(), &proto_track.Favorites{
			UserID:  userID,
			TrackID: track.ID,
		}).
		Return(nothing, nil)

	err := mockUsecase.DeleteFromFavourites(userID, track.ID)
	assert.Equal(t, err, nil)
}

func TestTrackUsecase_GetByPlaylistID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_track.NewMockTrackServiceClient(ctrl)
	mockUsecase := NewTrackUsecase(mockClient)

	track := models.Track{
		ID:       1,
		Title:    "title",
		AlbumID:  1,
		ArtistID: 1,
	}

	playlistID := uint64(1)

	grpcTracks := &proto_track.Tracks{}
	grpcTracks.Tracks = append(grpcTracks.Tracks, grpc_track.TrackToTrackGRPC(&track))

	mockClient.
		EXPECT().
		GetByPlaylistID(context.Background(), &proto_track.GetByPlaylistIdMessage{
			PlaylistId: playlistID,
			UserId:     0,
		}).
		Return(grpcTracks, nil)

	_, err := mockUsecase.GetByPlaylistID(playlistID, 0)
	assert.Equal(t, err, nil)
}

func TestTrackUsecase_LikeTrack(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_track.NewMockTrackServiceClient(ctrl)
	mockUsecase := NewTrackUsecase(mockClient)

	track := models.Track{
		ID:      1,
		Title:   "title",
		AlbumID: 1,
	}

	userID := uint64(1)

	nothing := new(proto_track.Nothing)

	mockClient.
		EXPECT().
		LikeTrack(context.Background(), &proto_track.Likes{
			UserId:  userID,
			TrackId: track.ID,
		}).
		Return(nothing, nil)

	err := mockUsecase.LikeTrack(userID, track.ID)
	assert.Equal(t, err, nil)
}

func TestTrackUsecase_DislikeTrack(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_track.NewMockTrackServiceClient(ctrl)
	mockUsecase := NewTrackUsecase(mockClient)

	track := models.Track{
		ID:      1,
		Title:   "title",
		AlbumID: 1,
	}

	userID := uint64(1)

	nothing := new(proto_track.Nothing)

	mockClient.
		EXPECT().
		DislikeTrack(context.Background(), &proto_track.Likes{
			UserId:  userID,
			TrackId: track.ID,
		}).
		Return(nothing, nil)

	err := mockUsecase.DislikeTrack(userID, track.ID)
	assert.Equal(t, err, nil)
}
