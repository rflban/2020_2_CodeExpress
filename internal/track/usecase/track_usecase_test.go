package usecase_test

import (
	"database/sql"
	"errors"
	"testing"

	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/consts"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/error_response"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/track/mock_track"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/track/usecase"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
)

func TestTrackUsecase_CreateTrack(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_track.NewMockTrackRep(ctrl)
	mockUsecase := usecase.NewTrackUsecase(mockRepo)

	id := uint64(42)
	index := uint8(50)
	track := &models.Track{}
	newTrack := &models.Track{ID: id, Index: index}

	mockRepo.
		EXPECT().
		Insert(gomock.Eq(track)).
		DoAndReturn(func(track *models.Track) error {
			track.ID = id
			return nil
		})

	mockRepo.
		EXPECT().
		SelectByID(id).
		Return(newTrack, nil)

	err := mockUsecase.CreateTrack(track)
	assert.Equal(t, err, nil)
	assert.Equal(t, track.Index, index)
}

func TestTrackUsecase_CreateTrack_Internal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_track.NewMockTrackRep(ctrl)
	mockUsecase := usecase.NewTrackUsecase(mockRepo)

	track := &models.Track{}
	dbErr := errors.New("Some db error")

	mockRepo.
		EXPECT().
		Insert(gomock.Eq(track)).
		Return(dbErr)

	err := mockUsecase.CreateTrack(track)
	assert.Equal(t, err, NewErrorResponse(ErrInternal, dbErr))
}

func TestTrackUsecase_CreateTrack_Failed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_track.NewMockTrackRep(ctrl)
	mockUsecase := usecase.NewTrackUsecase(mockRepo)

	id := uint64(42)
	track := &models.Track{}
	dbErr := errors.New("Some db error")

	mockRepo.
		EXPECT().
		Insert(gomock.Eq(track)).
		DoAndReturn(func(track *models.Track) error {
			track.ID = id
			return nil
		})

	mockRepo.
		EXPECT().
		SelectByID(id).
		Return(nil, dbErr)

	err := mockUsecase.CreateTrack(track)
	assert.Equal(t, err, NewErrorResponse(ErrInternal, dbErr))
}

func TestArtistUsecase_DeleteArtist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_track.NewMockTrackRep(ctrl)
	mockUsecase := usecase.NewTrackUsecase(mockRepo)

	id := uint64(5)
	mockRepo.
		EXPECT().
		Delete(gomock.Eq(id)).
		Return(nil)

	err := mockUsecase.DeleteTrack(id)
	assert.Equal(t, err, nil)
}

func TestArtistUsecase_DeleteArtist_Failed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_track.NewMockTrackRep(ctrl)
	mockUsecase := usecase.NewTrackUsecase(mockRepo)

	id := uint64(5)

	mockRepo.
		EXPECT().
		Delete(gomock.Eq(id)).
		Return(sql.ErrNoRows)

	err := mockUsecase.DeleteTrack(id)
	assert.Equal(t, err, NewErrorResponse(ErrTrackNotExist, sql.ErrNoRows))
}

func TestArtistUsecase_DeleteArtist_Internal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_track.NewMockTrackRep(ctrl)
	mockUsecase := usecase.NewTrackUsecase(mockRepo)

	id := uint64(5)
	dbErr := errors.New("Some db error")

	mockRepo.
		EXPECT().
		Delete(gomock.Eq(id)).
		Return(dbErr)

	err := mockUsecase.DeleteTrack(id)
	assert.Equal(t, err, NewErrorResponse(ErrInternal, dbErr))
}

func TestArtistUsecase_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_track.NewMockTrackRep(ctrl)
	mockUsecase := usecase.NewTrackUsecase(mockRepo)

	expectedTrack := &models.Track{
		ID: 1,
	}

	mockRepo.
		EXPECT().
		SelectByID(gomock.Eq(expectedTrack.ID)).
		Return(expectedTrack, nil)

	track, err := mockUsecase.GetByID(expectedTrack.ID)
	assert.Equal(t, err, nil)
	assert.Equal(t, track, expectedTrack)
}

func TestArtistUsecase_GetByID_Failed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_track.NewMockTrackRep(ctrl)
	mockUsecase := usecase.NewTrackUsecase(mockRepo)

	expectedTrack := &models.Track{
		ID: 1,
	}

	mockRepo.
		EXPECT().
		SelectByID(gomock.Eq(expectedTrack.ID)).
		Return(nil, sql.ErrNoRows)

	track, err := mockUsecase.GetByID(expectedTrack.ID)
	assert.Equal(t, track, nil)
	assert.Equal(t, err, NewErrorResponse(ErrTrackNotExist, sql.ErrNoRows))
}

func TestArtistUsecase_GetByID_Internal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_track.NewMockTrackRep(ctrl)
	mockUsecase := usecase.NewTrackUsecase(mockRepo)

	dbErr := errors.New("Some db error")
	expectedTrack := &models.Track{
		ID: 1,
	}

	mockRepo.
		EXPECT().
		SelectByID(gomock.Eq(expectedTrack.ID)).
		Return(nil, dbErr)

	track, err := mockUsecase.GetByID(expectedTrack.ID)
	assert.Equal(t, track, nil)
	assert.Equal(t, err, NewErrorResponse(ErrInternal, dbErr))
}

func TestArtistUsecase_GetByArtistID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_track.NewMockTrackRep(ctrl)
	mockUsecase := usecase.NewTrackUsecase(mockRepo)

	expectedTrack := &models.Track{
		ID:       1,
		ArtistID: 5,
	}

	expectedTracks := []*models.Track{expectedTrack}

	mockRepo.
		EXPECT().
		SelectByArtistID(gomock.Eq(expectedTrack.ArtistID)).
		Return(expectedTracks, nil)

	tracks, err := mockUsecase.GetByArtistID(expectedTrack.ArtistID)
	assert.Equal(t, err, nil)
	assert.Equal(t, tracks, expectedTracks)
}

func TestArtistUsecase_GetByArtistID_Failed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_track.NewMockTrackRep(ctrl)
	mockUsecase := usecase.NewTrackUsecase(mockRepo)

	expectedTrack := &models.Track{
		ID: 1,
	}

	mockRepo.
		EXPECT().
		SelectByID(gomock.Eq(expectedTrack.ID)).
		Return(nil, sql.ErrNoRows)

	track, err := mockUsecase.GetByID(expectedTrack.ID)
	assert.Equal(t, track, nil)
	assert.Equal(t, err, NewErrorResponse(ErrTrackNotExist, sql.ErrNoRows))
}

func TestArtistUsecase_GetByArtistID_Internal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_track.NewMockTrackRep(ctrl)
	mockUsecase := usecase.NewTrackUsecase(mockRepo)

	dbErr := errors.New("Some db error")
	expectedTrack := &models.Track{
		ID: 1,
	}

	mockRepo.
		EXPECT().
		SelectByID(gomock.Eq(expectedTrack.ID)).
		Return(nil, dbErr)

	track, err := mockUsecase.GetByID(expectedTrack.ID)
	assert.Equal(t, track, nil)
	assert.Equal(t, err, NewErrorResponse(ErrInternal, dbErr))
}

func TestArtistUsecase_GetByArtistParams(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_track.NewMockTrackRep(ctrl)
	mockUsecase := usecase.NewTrackUsecase(mockRepo)

	count, from := uint64(1), uint64(0)
	expectedTrack := &models.Track{
		ID:       1,
		ArtistID: 5,
	}

	expectedTracks := []*models.Track{expectedTrack}

	mockRepo.
		EXPECT().
		SelectByParam(gomock.Eq(count), gomock.Eq(from)).
		Return(expectedTracks, nil)

	tracks, err := mockUsecase.GetByParams(count, from)
	assert.Equal(t, err, nil)
	assert.Equal(t, tracks, expectedTracks)
}

func TestArtistUsecase_GetByArtistParams_Internal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_track.NewMockTrackRep(ctrl)
	mockUsecase := usecase.NewTrackUsecase(mockRepo)

	dbErr := errors.New("Some db error")
	count, from := uint64(1), uint64(0)

	mockRepo.
		EXPECT().
		SelectByParam(gomock.Eq(count), gomock.Eq(from)).
		Return(nil, dbErr)

	tracks, err := mockUsecase.GetByParams(count, from)
	assert.Equal(t, err, NewErrorResponse(ErrInternal, dbErr))
	assert.Equal(t, tracks, nil)
}

func TestArtistUsecase_GetByArtistParams_Failed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_track.NewMockTrackRep(ctrl)
	mockUsecase := usecase.NewTrackUsecase(mockRepo)

	count, from := uint64(1), uint64(0)

	mockRepo.
		EXPECT().
		SelectByParam(gomock.Eq(count), gomock.Eq(from)).
		Return(nil, sql.ErrNoRows)

	tracks, err := mockUsecase.GetByParams(count, from)
	assert.Equal(t, err, nil)
	assert.Equal(t, tracks, nil)
}

func TestTrackUsecase_UpdateTrack(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_track.NewMockTrackRep(ctrl)
	mockUsecase := usecase.NewTrackUsecase(mockRepo)

	id := uint64(42)
	index := uint8(50)
	track := &models.Track{ID: id}
	newTrack := &models.Track{ID: id, Index: index}

	mockRepo.
		EXPECT().
		Update(gomock.Eq(track)).
		Return(nil)

	mockRepo.
		EXPECT().
		SelectByID(id).
		Return(newTrack, nil)

	err := mockUsecase.UpdateTrack(track)
	assert.Equal(t, err, nil)
}

func TestTrackUsecase_UpdateTrackAudio(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_track.NewMockTrackRep(ctrl)
	mockUsecase := usecase.NewTrackUsecase(mockRepo)

	id := uint64(42)
	index := uint8(50)
	track := &models.Track{ID: id}
	newTrack := &models.Track{ID: id, Index: index}

	mockRepo.
		EXPECT().
		UpdateAudio(gomock.Eq(track)).
		Return(nil)

	mockRepo.
		EXPECT().
		SelectByID(id).
		Return(newTrack, nil)

	err := mockUsecase.UpdateTrackAudio(track)
	assert.Equal(t, err, nil)
}

func TestArtistUsecase_GetByAlbumID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_track.NewMockTrackRep(ctrl)
	mockUsecase := usecase.NewTrackUsecase(mockRepo)

	albumID := uint64(42)
	expectedTrack := &models.Track{
		AlbumID: albumID,
	}

	expectedTracks := []*models.Track{expectedTrack}

	mockRepo.
		EXPECT().
		SelectByAlbumID(albumID).
		Return(expectedTracks, nil)

	tracks, err := mockUsecase.GetByAlbumID(albumID)
	assert.Equal(t, err, nil)
	assert.Equal(t, tracks, expectedTracks)
}

func TestArtistUsecase_GetFavoritesByUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_track.NewMockTrackRep(ctrl)
	mockUsecase := usecase.NewTrackUsecase(mockRepo)

	userID := uint64(42)
	expectedTrack := &models.Track{
		ID: 5,
	}

	expectedTracks := []*models.Track{expectedTrack}

	mockRepo.
		EXPECT().
		SelectFavoritesByUserID(userID).
		Return(expectedTracks, nil)

	tracks, err := mockUsecase.GetFavoritesByUserID(userID)
	assert.Equal(t, err, nil)
	assert.Equal(t, tracks, expectedTracks)
}

func TestArtistUsecase_AddToFavorites(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_track.NewMockTrackRep(ctrl)
	mockUsecase := usecase.NewTrackUsecase(mockRepo)

	userID := uint64(42)
	trackID := uint64(7)

	mockRepo.
		EXPECT().
		SelectByID(gomock.Eq(trackID)).
		Return(nil, nil)

	mockRepo.
		EXPECT().
		InsertTrackToUser(gomock.Eq(userID), gomock.Eq(trackID)).
		Return(nil)

	err := mockUsecase.AddToFavourites(userID, trackID)
	assert.Equal(t, err, nil)
}

func TestArtistUsecase_DeleteFromFavourites(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_track.NewMockTrackRep(ctrl)
	mockUsecase := usecase.NewTrackUsecase(mockRepo)

	userID := uint64(42)
	trackID := uint64(7)

	mockRepo.
		EXPECT().
		SelectByID(gomock.Eq(trackID)).
		Return(nil, nil)

	mockRepo.
		EXPECT().
		DeleteTrackFromUsersTracks(gomock.Eq(userID), gomock.Eq(trackID)).
		Return(nil)

	err := mockUsecase.DeleteFromFavourites(userID, trackID)
	assert.Equal(t, err, nil)
}
