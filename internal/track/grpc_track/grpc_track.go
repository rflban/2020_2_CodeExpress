package grpc_track

import (
	"context"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/track"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/track/proto_track"
)

type TrackGRPCUsecase struct {
	trackRep track.TrackRep
}

func NewTrackGRPCUsecase(trackRep track.TrackRep) *TrackGRPCUsecase {
	return &TrackGRPCUsecase{
		trackRep: trackRep,
	}
}

func TrackGRPCToTrack(track *proto_track.Track) *models.Track {
	return &models.Track{
		ID:          track.ID,
		Title:       track.Title,
		Duration:    int(track.Duration),
		AlbumPoster: track.AlbumPoster,
		AlbumID:     track.AlbumID,
		Index:       uint8(track.Index),
		Audio:       track.Audio,
		Artist:      track.Artist,
		ArtistID:    track.ArtistID,
		IsFavorite:  track.IsFavorite,
		IsLiked:     track.IsLiked,
	}
}

func TrackToTrackGRPC(track *models.Track) *proto_track.Track {
	return &proto_track.Track{
		ID:          track.ID,
		Title:       track.Title,
		Duration:    int64(track.Duration),
		AlbumPoster: track.AlbumPoster,
		AlbumID:     track.AlbumID,
		Index:       uint32(track.Index),
		Audio:       track.Audio,
		Artist:      track.Artist,
		ArtistID:    track.ArtistID,
		IsFavorite:  track.IsFavorite,
		IsLiked:     track.IsLiked,
	}
}

func (tGu *TrackGRPCUsecase) CreateTrack(ctx context.Context, grpcTrack *proto_track.Track) (*proto_track.Track, error) {
	track := TrackGRPCToTrack(grpcTrack)

	if err := tGu.trackRep.Insert(track); err != nil {
		return grpcTrack, err
	}

	grpcTrack.ID = track.ID

	return grpcTrack, nil
}

func (tGu *TrackGRPCUsecase) DeleteTrack(ctx context.Context, trackID *proto_track.TrackID) (*proto_track.Nothing, error) {
	nothing := new(proto_track.Nothing)

	err := tGu.trackRep.Delete(trackID.ID)

	if err != nil {
		return nothing, err
	}

	return nothing, nil
}

func (tGu *TrackGRPCUsecase) GetByArtistId(ctx context.Context, mes *proto_track.GetByArtistIdMessage) (*proto_track.Tracks, error) {
	tracks, err := tGu.trackRep.SelectByArtistId(mes.ArtistID, mes.UserID)

	if err != nil {
		return new(proto_track.Tracks), err
	}

	grpcTracks := &proto_track.Tracks{}

	for _, track := range tracks {
		grpcTracks.Tracks = append(grpcTracks.Tracks, TrackToTrackGRPC(track))
	}

	return grpcTracks, nil
}

func (tGu *TrackGRPCUsecase) GetByAlbumID(ctx context.Context, mes *proto_track.GetByAlbumIdMessage) (*proto_track.Tracks, error) {
	tracks, err := tGu.trackRep.SelectByAlbumID(mes.AlbumId, mes.UserId)

	if err != nil {
		return new(proto_track.Tracks), err
	}

	grpcTracks := &proto_track.Tracks{}

	for _, track := range tracks {
		grpcTracks.Tracks = append(grpcTracks.Tracks, TrackToTrackGRPC(track))
	}

	return grpcTracks, nil
}

func (tGu *TrackGRPCUsecase) GetByID(ctx context.Context, mes *proto_track.GetByIdMessage) (*proto_track.Track, error) {
	track, err := tGu.trackRep.SelectByID(mes.TrackId, mes.UserId)

	if err != nil {
		return new(proto_track.Track), err
	}

	return TrackToTrackGRPC(track), nil
}

func (tGu *TrackGRPCUsecase) GetByParams(ctx context.Context, mes *proto_track.GetByParamsMessage) (*proto_track.Tracks, error) {
	tracks, err := tGu.trackRep.SelectByParams(mes.Count, mes.From, mes.UserID)

	if err != nil {
		return new(proto_track.Tracks), err
	}

	grpcTracks := &proto_track.Tracks{}

	for _, track := range tracks {
		grpcTracks.Tracks = append(grpcTracks.Tracks, TrackToTrackGRPC(track))
	}

	return grpcTracks, nil
}

func (tGu *TrackGRPCUsecase) GetTopByParams(ctx context.Context, mes *proto_track.GetTopByParamsMessage) (*proto_track.Tracks, error) {
	tracks, err := tGu.trackRep.SelectTopByParams(mes.Count, mes.From, mes.UserID)

	if err != nil {
		return new(proto_track.Tracks), err
	}

	grpcTracks := &proto_track.Tracks{}

	for _, track := range tracks {
		grpcTracks.Tracks = append(grpcTracks.Tracks, TrackToTrackGRPC(track))
	}

	return grpcTracks, nil
}

func (tGu *TrackGRPCUsecase) GetFavoritesByUserID(ctx context.Context, userID *proto_track.UserID) (*proto_track.Tracks, error) {
	tracks, err := tGu.trackRep.SelectFavoritesByUserID(userID.ID)

	if err != nil {
		return new(proto_track.Tracks), err
	}

	grpcTracks := &proto_track.Tracks{}

	for _, track := range tracks {
		grpcTracks.Tracks = append(grpcTracks.Tracks, TrackToTrackGRPC(track))
	}

	return grpcTracks, nil
}

func (tGu *TrackGRPCUsecase) UpdateTrack(ctx context.Context, track *proto_track.Track) (*proto_track.Nothing, error) {
	nothing := new(proto_track.Nothing)
	err := tGu.trackRep.Update(TrackGRPCToTrack(track))

	if err != nil {
		return nothing, err
	}

	return nothing, nil
}

func (tGu *TrackGRPCUsecase) UpdateTrackAudio(ctx context.Context, track *proto_track.Track) (*proto_track.Nothing, error) {
	nothing := new(proto_track.Nothing)
	err := tGu.trackRep.UpdateAudio(TrackGRPCToTrack(track))

	if err != nil {
		return nothing, err
	}

	return nothing, nil
}

func (tGu *TrackGRPCUsecase) AddToFavourites(ctx context.Context, mes *proto_track.Favorites) (*proto_track.Nothing, error) {
	nothing := new(proto_track.Nothing)
	err := tGu.trackRep.InsertTrackToUser(mes.UserID, mes.TrackID)

	if err != nil {
		return nothing, err
	}

	return nothing, nil
}

func (tGu *TrackGRPCUsecase) DeleteFromFavourites(ctx context.Context, mes *proto_track.Favorites) (*proto_track.Nothing, error) {
	nothing := new(proto_track.Nothing)
	err := tGu.trackRep.DeleteTrackFromUsersTracks(mes.UserID, mes.TrackID)

	if err != nil {
		return nothing, err
	}

	return nothing, nil
}

func (tGu *TrackGRPCUsecase) GetByPlaylistID(ctx context.Context, mes *proto_track.GetByPlaylistIdMessage) (*proto_track.Tracks, error) {
	tracks, err := tGu.trackRep.SelectByPlaylistID(mes.PlaylistId, mes.UserId)

	if err != nil {
		return new(proto_track.Tracks), err
	}

	grpcTracks := &proto_track.Tracks{}

	for _, track := range tracks {
		grpcTracks.Tracks = append(grpcTracks.Tracks, TrackToTrackGRPC(track))
	}

	return grpcTracks, nil
}

func (tGu *TrackGRPCUsecase) LikeTrack(ctx context.Context, mes *proto_track.Likes) (*proto_track.Nothing, error) {
	nothing := new(proto_track.Nothing)
	if err := tGu.trackRep.LikeTrack(mes.UserId, mes.TrackId); err != nil {
		return nothing, err
	}

	return nothing, nil
}

func (tGu *TrackGRPCUsecase) DislikeTrack(ctx context.Context, mes *proto_track.Likes) (*proto_track.Nothing, error) {
	nothing := new(proto_track.Nothing)
	if err := tGu.trackRep.DislikeTrack(mes.UserId, mes.TrackId); err != nil {
		return nothing, err
	}

	return nothing, nil
}
