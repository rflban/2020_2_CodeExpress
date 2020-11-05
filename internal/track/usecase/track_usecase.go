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

	copier.Copy(&track, &newTrack)

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

func (aUc *TrackUsecase) GetByArtistID(artistID uint64) ([]*models.Track, *ErrorResponse) {
	tracks, err := aUc.trackRep.SelectByArtistID(artistID)

	if err == sql.ErrNoRows {
		return nil, NewErrorResponse(ErrArtistNotExist, err)
	}

	if err != nil {
		return nil, NewErrorResponse(ErrInternal, err)
	}

	return tracks, nil
}

func (aUc *TrackUsecase) GetByParams(count uint64, from uint64) ([]*models.Track, *ErrorResponse) {
	tracks, err := aUc.trackRep.SelectByParam(count, from)

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

	copier.Copy(&track, &newTrack)

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

	copier.Copy(&track, &newTrack)

	return nil
}
