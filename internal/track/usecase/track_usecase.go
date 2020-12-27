package usecase

import (
	"database/sql"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/track/grpc_track"
	"golang.org/x/net/context"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/track/proto_track"

	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/consts"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/error_response"
	"github.com/jinzhu/copier"
)

type TrackUsecase struct {
	trackGRPC proto_track.TrackServiceClient
}

func NewTrackUsecase(trackGRPC proto_track.TrackServiceClient) *TrackUsecase {
	return &TrackUsecase{
		trackGRPC: trackGRPC,
	}
}

func (aUc *TrackUsecase) CreateTrack(track *models.Track, userId uint64) *ErrorResponse {
	grpcTrack, err := aUc.trackGRPC.CreateTrack(context.Background(), grpc_track.TrackToTrackGRPC(track))

	if err != nil {
		return NewErrorResponse(ErrInternal, err)
	}

	newTrack, err := aUc.trackGRPC.GetByID(context.Background(), &proto_track.GetByIdMessage{
		TrackId: grpcTrack.ID,
		UserId:  userId,
	})

	if err != nil {
		return NewErrorResponse(ErrInternal, err)
	}

	_ = copier.Copy(&track, &newTrack)

	return nil
}

func (aUc *TrackUsecase) DeleteTrack(id uint64) *ErrorResponse {
	_, err := aUc.trackGRPC.DeleteTrack(context.Background(), &proto_track.TrackID{
		ID: id,
	})

	if err == sql.ErrNoRows {
		return NewErrorResponse(ErrTrackNotExist, err)
	}

	if err != nil {
		return NewErrorResponse(ErrInternal, err)
	}

	return nil
}

func (aUc *TrackUsecase) GetByID(id, userId uint64) (*models.Track, *ErrorResponse) {
	track, err := aUc.trackGRPC.GetByID(context.Background(), &proto_track.GetByIdMessage{
		TrackId: id,
		UserId:  userId,
	})

	if err == sql.ErrNoRows {
		return nil, NewErrorResponse(ErrTrackNotExist, err)
	}

	if err != nil {
		return nil, NewErrorResponse(ErrInternal, err)
	}

	return grpc_track.TrackGRPCToTrack(track), nil
}

func (aUc *TrackUsecase) GetByArtistId(artistId, userId uint64) ([]*models.Track, *ErrorResponse) {
	grpcTracks, err := aUc.trackGRPC.GetByArtistId(context.Background(), &proto_track.GetByArtistIdMessage{
		ArtistID: artistId,
		UserID:   userId,
	})

	if err == sql.ErrNoRows {
		return nil, NewErrorResponse(ErrArtistNotExist, err)
	}
	if err != nil {
		return nil, NewErrorResponse(ErrInternal, err)
	}

	tracks := make([]*models.Track, len(grpcTracks.Tracks))

	for idx, track := range grpcTracks.Tracks {
		tracks[idx] = grpc_track.TrackGRPCToTrack(track)
	}

	return tracks, nil
}

func (aUc *TrackUsecase) GetByParams(count, from, userId uint64) ([]*models.Track, *ErrorResponse) {
	grpcTracks, err := aUc.trackGRPC.GetByParams(context.Background(), &proto_track.GetByParamsMessage{
		Count:  count,
		From:   from,
		UserID: userId,
	})

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, NewErrorResponse(ErrInternal, err)
	}

	tracks := make([]*models.Track, len(grpcTracks.Tracks))

	for idx, track := range grpcTracks.Tracks {
		tracks[idx] = grpc_track.TrackGRPCToTrack(track)
	}

	return tracks, nil
}

func (aUc *TrackUsecase) GetTopByParams(count, from, userId uint64) ([]*models.Track, *ErrorResponse) {
	grpcTracks, err := aUc.trackGRPC.GetTopByParams(context.Background(), &proto_track.GetTopByParamsMessage{
		Count:  count,
		From:   from,
		UserID: userId,
	})

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, NewErrorResponse(ErrInternal, err)
	}

	tracks := make([]*models.Track, len(grpcTracks.Tracks))

	for idx, track := range grpcTracks.Tracks {
		tracks[idx] = grpc_track.TrackGRPCToTrack(track)
	}

	return tracks, nil
}

func (aUc *TrackUsecase) UpdateTrack(track *models.Track, userId uint64) *ErrorResponse {
	_, err := aUc.trackGRPC.UpdateTrack(context.Background(), grpc_track.TrackToTrackGRPC(track))

	if err == sql.ErrNoRows {
		return NewErrorResponse(ErrTrackNotExist, err)
	}

	if err != nil {
		return NewErrorResponse(ErrInternal, err)
	}

	newTrack, err := aUc.trackGRPC.GetByID(context.Background(), &proto_track.GetByIdMessage{
		TrackId: track.ID,
		UserId:  userId,
	})

	if err != nil {
		return NewErrorResponse(ErrInternal, err)
	}
	//nolint:staticcheck
	track = grpc_track.TrackGRPCToTrack(newTrack)

	return nil
}

func (aUc *TrackUsecase) UpdateTrackAudio(track *models.Track, userId uint64) *ErrorResponse {
	_, err := aUc.trackGRPC.UpdateTrackAudio(context.Background(), grpc_track.TrackToTrackGRPC(track))

	if err == sql.ErrNoRows {
		return NewErrorResponse(ErrTrackNotExist, err)
	}

	if err != nil {
		return NewErrorResponse(ErrInternal, err)
	}

	newTrack, err := aUc.trackGRPC.GetByID(context.Background(), &proto_track.GetByIdMessage{
		TrackId: track.ID,
		UserId:  userId,
	})

	if err != nil {
		return NewErrorResponse(ErrInternal, err)
	}
	//nolint:staticcheck
	track = grpc_track.TrackGRPCToTrack(newTrack)

	return nil
}

func (aUc *TrackUsecase) GetByAlbumID(albumId, userId uint64) ([]*models.Track, *ErrorResponse) {
	grpcTracks, err := aUc.trackGRPC.GetByAlbumID(context.Background(), &proto_track.GetByAlbumIdMessage{
		AlbumId: albumId,
		UserId:  userId,
	})

	if err == sql.ErrNoRows {
		return nil, NewErrorResponse(ErrArtistNotExist, err)
	}

	if err != nil {
		return nil, NewErrorResponse(ErrInternal, err)
	}

	tracks := make([]*models.Track, len(grpcTracks.Tracks))

	for idx, track := range grpcTracks.Tracks {
		tracks[idx] = grpc_track.TrackGRPCToTrack(track)
	}

	return tracks, nil
}

func (aUc *TrackUsecase) GetFavoritesByUserID(userID uint64) ([]*models.Track, *ErrorResponse) {
	grpcTracks, err := aUc.trackGRPC.GetFavoritesByUserID(context.Background(), &proto_track.UserID{
		ID: userID,
	})

	if err == sql.ErrNoRows {
		return nil, NewErrorResponse(ErrNoFavoritesTracks, err)
	}

	if err != nil {
		return nil, NewErrorResponse(ErrInternal, err)
	}

	tracks := make([]*models.Track, len(grpcTracks.Tracks))

	for idx, track := range grpcTracks.Tracks {
		tracks[idx] = grpc_track.TrackGRPCToTrack(track)
	}

	return tracks, nil
}

func (aUc *TrackUsecase) AddToFavourites(userId, trackId uint64) *ErrorResponse {
	_, err := aUc.trackGRPC.GetByID(context.Background(), &proto_track.GetByIdMessage{
		TrackId: trackId,
		UserId:  userId,
	})

	if err == sql.ErrNoRows {
		return NewErrorResponse(ErrTrackNotExist, err)
	}

	_, err = aUc.trackGRPC.AddToFavourites(context.Background(), &proto_track.Favorites{
		UserID:  userId,
		TrackID: trackId,
	})

	if err != nil {
		return NewErrorResponse(ErrInternal, err)
	}

	return nil
}

func (aUc *TrackUsecase) DeleteFromFavourites(userId, trackId uint64) *ErrorResponse {
	_, err := aUc.trackGRPC.GetByID(context.Background(), &proto_track.GetByIdMessage{
		TrackId: trackId,
		UserId:  userId,
	})

	if err == sql.ErrNoRows {
		return NewErrorResponse(ErrTrackNotExist, err)
	}

	_, err = aUc.trackGRPC.DeleteFromFavourites(context.Background(), &proto_track.Favorites{
		UserID:  userId,
		TrackID: trackId,
	})

	if err != nil {
		return NewErrorResponse(ErrInternal, err)
	}

	return nil
}

func (aUc *TrackUsecase) GetByPlaylistID(playlistId, userId uint64) ([]*models.Track, *ErrorResponse) {
	grpcTracks, err := aUc.trackGRPC.GetByPlaylistID(context.Background(), &proto_track.GetByPlaylistIdMessage{
		PlaylistId: playlistId,
		UserId:     userId,
	})

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, NewErrorResponse(ErrInternal, err)
	}

	tracks := make([]*models.Track, len(grpcTracks.Tracks))

	for idx, track := range grpcTracks.Tracks {
		tracks[idx] = grpc_track.TrackGRPCToTrack(track)
	}

	return tracks, nil
}

func (aUc *TrackUsecase) LikeTrack(userId uint64, trackId uint64) *ErrorResponse {
	_, err := aUc.trackGRPC.LikeTrack(context.Background(), &proto_track.Likes{
		UserId:  userId,
		TrackId: trackId,
	})
	if err != nil {
		return NewErrorResponse(ErrInternal, err)
	}

	return nil
}

func (aUc *TrackUsecase) DislikeTrack(userId uint64, trackId uint64) *ErrorResponse {
	_, err := aUc.trackGRPC.DislikeTrack(context.Background(), &proto_track.Likes{
		UserId:  userId,
		TrackId: trackId,
	})
	if err != nil {
		return NewErrorResponse(ErrInternal, err)
	}

	return nil
}

func (aUc *TrackUsecase) GetRandomByArtistId(artistId, userId uint64, count uint64) ([]*models.Track, *ErrorResponse) {
	grpcTracks, err := aUc.trackGRPC.GetRandomByArtistID(context.Background(), &proto_track.RandomArtist{
		ArtistId: artistId,
		UserId:   userId,
		Count:    count,
	})

	if err == sql.ErrNoRows {
		return nil, NewErrorResponse(ErrArtistNotExist, err)
	}
	if err != nil {
		return nil, NewErrorResponse(ErrInternal, err)
	}

	tracks := make([]*models.Track, len(grpcTracks.Tracks))

	for idx, track := range grpcTracks.Tracks {
		tracks[idx] = grpc_track.TrackGRPCToTrack(track)
	}

	return tracks[:count], nil
}
