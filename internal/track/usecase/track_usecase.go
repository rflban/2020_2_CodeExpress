package usecase

import (
	"database/sql"
	"fmt"

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

func (aUc *TrackUsecase) CreateTrack(track *models.Track) *ErrorResponse {
	grpcTrack, err := aUc.trackGRPC.CreateTrack(context.Background(), grpc_track.TrackToTrackGRPC(track))

	if err != nil {
		fmt.Println("HERE", err)
		return NewErrorResponse(ErrInternal, err)
	}

	newTrack, err := aUc.trackGRPC.GetByID(context.Background(), &proto_track.TrackID{
		ID: grpcTrack.ID,
	})

	if err != nil {
		fmt.Println("HERE1", err)
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

func (aUc *TrackUsecase) GetByID(id uint64) (*models.Track, *ErrorResponse) {
	track, err := aUc.trackGRPC.GetByID(context.Background(), &proto_track.TrackID{
		ID: id,
	})

	if err == sql.ErrNoRows {
		return nil, NewErrorResponse(ErrTrackNotExist, err)
	}

	if err != nil {
		return nil, NewErrorResponse(ErrInternal, err)
	}

	return grpc_track.TrackGRPCToTrack(track), nil
}

func (aUc *TrackUsecase) GetByArtistId(artistId uint64, userId uint64) ([]*models.Track, *ErrorResponse) {
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

func (aUc *TrackUsecase) GetByParams(count uint64, from uint64, userId uint64) ([]*models.Track, *ErrorResponse) {
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

func (aUc *TrackUsecase) UpdateTrack(track *models.Track) *ErrorResponse {
	_, err := aUc.trackGRPC.UpdateTrack(context.Background(), grpc_track.TrackToTrackGRPC(track))

	if err == sql.ErrNoRows {
		return NewErrorResponse(ErrTrackNotExist, err)
	}

	if err != nil {
		return NewErrorResponse(ErrInternal, err)
	}

	newTrack, err := aUc.trackGRPC.GetByID(context.Background(), &proto_track.TrackID{
		ID: track.ID,
	})

	if err != nil {
		return NewErrorResponse(ErrInternal, err)
	}

	track = grpc_track.TrackGRPCToTrack(newTrack)

	return nil
}

func (aUc *TrackUsecase) UpdateTrackAudio(track *models.Track) *ErrorResponse {
	_, err := aUc.trackGRPC.UpdateTrackAudio(context.Background(), grpc_track.TrackToTrackGRPC(track))

	if err == sql.ErrNoRows {
		return NewErrorResponse(ErrTrackNotExist, err)
	}

	if err != nil {
		return NewErrorResponse(ErrInternal, err)
	}

	newTrack, err := aUc.trackGRPC.GetByID(context.Background(), &proto_track.TrackID{
		ID: track.ID,
	})

	if err != nil {
		return NewErrorResponse(ErrInternal, err)
	}

	track = grpc_track.TrackGRPCToTrack(newTrack)

	return nil
}

func (aUc *TrackUsecase) GetByAlbumID(albumID uint64) ([]*models.Track, *ErrorResponse) {
	grpcTracks, err := aUc.trackGRPC.GetByAlbumID(context.Background(), &proto_track.AlbumID{
		ID: albumID,
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

func (aUc *TrackUsecase) AddToFavourites(userID uint64, trackID uint64) *ErrorResponse {
	_, err := aUc.trackGRPC.GetByID(context.Background(), &proto_track.TrackID{
		ID: trackID,
	})

	if err == sql.ErrNoRows {
		return NewErrorResponse(ErrTrackNotExist, err)
	}

	_, err = aUc.trackGRPC.AddToFavourites(context.Background(), &proto_track.Favorites{
		UserID:  userID,
		TrackID: trackID,
	})

	if err != nil {
		return NewErrorResponse(ErrInternal, err)
	}

	return nil
}

func (aUc *TrackUsecase) DeleteFromFavourites(userID uint64, trackID uint64) *ErrorResponse {
	_, err := aUc.trackGRPC.GetByID(context.Background(), &proto_track.TrackID{
		ID: trackID,
	})

	if err == sql.ErrNoRows {
		return NewErrorResponse(ErrTrackNotExist, err)
	}

	_, err = aUc.trackGRPC.DeleteFromFavourites(context.Background(), &proto_track.Favorites{
		UserID:  userID,
		TrackID: trackID,
	})

	if err != nil {
		return NewErrorResponse(ErrInternal, err)
	}

	return nil
}

func (aUc *TrackUsecase) GetByPlaylistID(playlistID uint64) ([]*models.Track, *ErrorResponse) {
	grpcTracks, err := aUc.trackGRPC.GetByPlaylistID(context.Background(), &proto_track.PlaylistID{
		ID: playlistID,
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
