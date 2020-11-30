package usecase

import (
	"database/sql"

	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/consts"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/error_response"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/track"
	"github.com/jinzhu/copier"
)

type TrackUsecase struct {
	trackRep track.TrackRep
}

func NewTrackUsecase(trackRep track.TrackRep) *TrackUsecase {
	return &TrackUsecase{
		trackRep: trackRep,
	}
}

func (aUc *TrackUsecase) CreateTrack(track *models.Track) *ErrorResponse {
	if err := aUc.trackRep.Insert(track); err != nil {
		return NewErrorResponse(ErrInternal, err)
	}

	newTrack, err := aUc.trackRep.SelectByID(track.ID)

	if err != nil {
		return NewErrorResponse(ErrInternal, err)
	}

	_ = copier.Copy(&track, &newTrack)

	return nil
}

func (aUc *TrackUsecase) DeleteTrack(id uint64) *ErrorResponse {
	err := aUc.trackRep.Delete(id)

	if err == sql.ErrNoRows {
		return NewErrorResponse(ErrTrackNotExist, err)
	}

	if err != nil {
		return NewErrorResponse(ErrInternal, err)
	}

	return nil
}

func (aUc *TrackUsecase) GetByID(id uint64) (*models.Track, *ErrorResponse) {
	track, err := aUc.trackRep.SelectByID(id)

	if err == sql.ErrNoRows {
		return nil, NewErrorResponse(ErrTrackNotExist, err)
	}

	if err != nil {
		return nil, NewErrorResponse(ErrInternal, err)
	}

	return track, nil
}

func (aUc *TrackUsecase) GetByArtistId(artistId uint64, userId uint64) ([]*models.Track, *ErrorResponse) {
	tracks, err := aUc.trackRep.SelectByArtistId(artistId, userId)
	if err == sql.ErrNoRows {
		return nil, NewErrorResponse(ErrArtistNotExist, err)
	}
	if err != nil {
		return nil, NewErrorResponse(ErrInternal, err)
	}

	return tracks, nil
}

func (aUc *TrackUsecase) GetByParams(count uint64, from uint64, userId uint64) ([]*models.Track, *ErrorResponse) {
	tracks, err := aUc.trackRep.SelectByParams(count, from, userId)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, NewErrorResponse(ErrInternal, err)
	}

	return tracks, nil
}

func (aUc *TrackUsecase) UpdateTrack(track *models.Track) *ErrorResponse {
	err := aUc.trackRep.Update(track)

	if err == sql.ErrNoRows {
		return NewErrorResponse(ErrTrackNotExist, err)
	}

	if err != nil {
		return NewErrorResponse(ErrInternal, err)
	}

	newTrack, err := aUc.trackRep.SelectByID(track.ID)

	if err != nil {
		return NewErrorResponse(ErrInternal, err)
	}

	_ = copier.Copy(&track, &newTrack)

	return nil
}

func (aUc *TrackUsecase) UpdateTrackAudio(track *models.Track) *ErrorResponse {
	err := aUc.trackRep.UpdateAudio(track)

	if err == sql.ErrNoRows {
		return NewErrorResponse(ErrTrackNotExist, err)
	}

	if err != nil {
		return NewErrorResponse(ErrInternal, err)
	}

	newTrack, err := aUc.trackRep.SelectByID(track.ID)

	if err != nil {
		return NewErrorResponse(ErrInternal, err)
	}

	_ = copier.Copy(&track, &newTrack)

	return nil
}

func (aUc *TrackUsecase) GetByAlbumID(albumID uint64) ([]*models.Track, *ErrorResponse) {
	tracks, err := aUc.trackRep.SelectByAlbumID(albumID)

	if err == sql.ErrNoRows {
		return nil, NewErrorResponse(ErrArtistNotExist, err)
	}

	if err != nil {
		return nil, NewErrorResponse(ErrInternal, err)
	}

	return tracks, nil
}

func (aUc *TrackUsecase) GetFavoritesByUserID(userID uint64) ([]*models.Track, *ErrorResponse) {
	tracks, err := aUc.trackRep.SelectFavoritesByUserID(userID)

	if err == sql.ErrNoRows {
		return nil, NewErrorResponse(ErrNoFavoritesTracks, err)
	}

	if err != nil {
		return nil, NewErrorResponse(ErrInternal, err)
	}

	return tracks, nil
}

func (aUc *TrackUsecase) AddToFavourites(userID uint64, trackID uint64) *ErrorResponse {
	_, err := aUc.trackRep.SelectByID(trackID)

	if err == sql.ErrNoRows {
		return NewErrorResponse(ErrTrackNotExist, err)
	}

	if err := aUc.trackRep.InsertTrackToUser(userID, trackID); err != nil {
		return NewErrorResponse(ErrInternal, err)
	}

	return nil
}

func (aUc *TrackUsecase) DeleteFromFavourites(userID uint64, trackID uint64) *ErrorResponse {
	_, err := aUc.trackRep.SelectByID(trackID)

	if err == sql.ErrNoRows {
		return NewErrorResponse(ErrTrackNotExist, err)
	}

	if err := aUc.trackRep.DeleteTrackFromUsersTracks(userID, trackID); err != nil {
		return NewErrorResponse(ErrInternal, err)
	}

	return nil
}

func (aUc *TrackUsecase) GetByPlaylistID(playlistID uint64) ([]*models.Track, *ErrorResponse) {
	tracks, err := aUc.trackRep.SelectByPlaylistID(playlistID)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, NewErrorResponse(ErrInternal, err)
	}

	return tracks, nil
}
