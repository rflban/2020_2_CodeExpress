package usecase

import (
	"database/sql"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/playlist"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/error_response"

	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/consts"
)

type PlaylistUsecase struct {
	playlistRep playlist.PlaylistRep
}

func NewPlaylistUsecase(playlistRep playlist.PlaylistRep) playlist.PlaylistUsecase {
	return &PlaylistUsecase{
		playlistRep: playlistRep,
	}
}

func (pu *PlaylistUsecase) CreatePlaylist(playlist *models.Playlist) *ErrorResponse {
	if err := pu.playlistRep.Insert(playlist); err != nil {
		return NewErrorResponse(ErrInternal, err)
	}
	return nil
}

func (pu *PlaylistUsecase) UpdatePlaylist(playlist *models.Playlist) *ErrorResponse {
	err := pu.playlistRep.Update(playlist)

	if err == sql.ErrNoRows {
		return NewErrorResponse(ErrPlaylistNotExist, err)
	}

	if err != nil {
		return NewErrorResponse(ErrInternal, err)
	}

	return nil
}

func (pu *PlaylistUsecase) DeletePlaylist(id uint64) *ErrorResponse {
	err := pu.playlistRep.Delete(id)

	if err == sql.ErrNoRows {
		return NewErrorResponse(ErrPlaylistNotExist, err)
	}

	if err != nil {
		return NewErrorResponse(ErrInternal, err)
	}

	return nil
}

func (pu *PlaylistUsecase) GetByID(id uint64) (*models.Playlist, *ErrorResponse) {
	playlist, err := pu.playlistRep.SelectByID(id)

	if err == sql.ErrNoRows {
		return nil, NewErrorResponse(ErrPlaylistNotExist, err)
	}

	if err != nil {
		return nil, NewErrorResponse(ErrInternal, err)
	}

	return playlist, nil
}

func (pu *PlaylistUsecase) GetByUserID(userID uint64) ([]*models.Playlist, *ErrorResponse) {
	playlists, err := pu.playlistRep.SelectByUserID(userID)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, NewErrorResponse(ErrInternal, err)
	}

	return playlists, nil
}

func (pu *PlaylistUsecase) AddTrack(trackID uint64, playlistID uint64) *ErrorResponse {
	err := pu.playlistRep.InsertTrack(trackID, playlistID)

	if err == sql.ErrNoRows {
		return NewErrorResponse(ErrPlaylistNotExist, err)
	}

	if err != nil {
		return NewErrorResponse(ErrInternal, err)
	}

	return nil
}

func (pu *PlaylistUsecase) DeleteTrack(trackID uint64, playlistID uint64) *ErrorResponse {
	err := pu.playlistRep.DeleteTrack(trackID, playlistID)

	if err == sql.ErrNoRows {
		return NewErrorResponse(ErrPlaylistNotExist, err)
	}

	if err != nil {
		return NewErrorResponse(ErrInternal, err)
	}

	return nil
}

func (pu *PlaylistUsecase) GetPublicByUserID(userID uint64) ([]*models.Playlist, *ErrorResponse) {
	playlists, err := pu.playlistRep.SelectPublicByUserID(userID)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, NewErrorResponse(ErrInternal, err)
	}

	return playlists, nil
}
