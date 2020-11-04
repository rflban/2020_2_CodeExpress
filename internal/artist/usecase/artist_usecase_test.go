package usecase_test

import (
	"database/sql"
	"testing"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/artist/usecase"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/consts"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/error_response"
	"github.com/go-playground/assert/v2"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/artist/mock_artist"

	"github.com/golang/mock/gomock"
)

func TestArtistUsecase_CreateArtist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_artist.NewMockArtistRep(ctrl)
	mockUsecase := usecase.NewArtistUsecase(mockRepo)

	artist := &models.Artist{
		Name: "Imagine Dragons",
	}

	mockRepo.
		EXPECT().
		Insert(gomock.Eq(artist)).
		Return(nil)

	mockRepo.
		EXPECT().
		SelectByName(gomock.Eq(artist.Name)).
		Return(nil, sql.ErrNoRows)

	err := mockUsecase.CreateArtist(artist)
	assert.Equal(t, err, nil)
}

func TestArtistUsecase_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_artist.NewMockArtistRep(ctrl)
	mockUsecase := usecase.NewArtistUsecase(mockRepo)

	expectedArtist := &models.Artist{
		ID:   1,
		Name: "Imagine Dragons",
	}

	mockRepo.
		EXPECT().
		SelectByID(gomock.Eq(expectedArtist.ID)).
		Return(expectedArtist, nil)

	artist, err := mockUsecase.GetByID(expectedArtist.ID)
	assert.Equal(t, err, nil)
	assert.Equal(t, expectedArtist, artist)
}

func TestArtistUsecase_GetByID_Failed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_artist.NewMockArtistRep(ctrl)
	mockUsecase := usecase.NewArtistUsecase(mockRepo)

	expectedArtist := &models.Artist{
		ID:   1,
		Name: "Imagine Dragons",
	}

	mockRepo.
		EXPECT().
		SelectByID(gomock.Eq(expectedArtist.ID)).
		Return(nil, sql.ErrNoRows)

	artist, err := mockUsecase.GetByID(expectedArtist.ID)
	assert.Equal(t, err, NewErrorResponse(ErrArtistNotExist, sql.ErrNoRows))
	assert.Equal(t, nil, artist)
}

func TestArtistUsecase_DeleteArtist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_artist.NewMockArtistRep(ctrl)
	mockUsecase := usecase.NewArtistUsecase(mockRepo)

	expectedArtist := &models.Artist{
		ID:   1,
		Name: "Imagine Dragons",
	}

	mockRepo.
		EXPECT().
		Delete(gomock.Eq(expectedArtist.ID)).
		Return(nil)

	err := mockUsecase.DeleteArtist(expectedArtist.ID)
	assert.Equal(t, err, nil)
}

func TestArtistUsecase_DeleteArtist_Failed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_artist.NewMockArtistRep(ctrl)
	mockUsecase := usecase.NewArtistUsecase(mockRepo)

	expectedArtist := &models.Artist{
		ID:   1,
		Name: "Imagine Dragons",
	}

	mockRepo.
		EXPECT().
		Delete(gomock.Eq(expectedArtist.ID)).
		Return(sql.ErrNoRows)

	err := mockUsecase.DeleteArtist(expectedArtist.ID)
	assert.Equal(t, err, NewErrorResponse(ErrArtistNotExist, sql.ErrNoRows))
}

func TestArtistUsecase_UpdateArtistName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_artist.NewMockArtistRep(ctrl)
	mockUsecase := usecase.NewArtistUsecase(mockRepo)

	expectedArtist := &models.Artist{
		ID:   1,
		Name: "Imagine Dragons",
	}

	mockRepo.
		EXPECT().
		SelectByName(gomock.Eq(expectedArtist.Name)).
		Return(nil, sql.ErrNoRows)

	mockRepo.
		EXPECT().
		UpdateName(gomock.Eq(expectedArtist)).
		Return(nil)

	err := mockUsecase.UpdateArtistName(expectedArtist)
	assert.Equal(t, err, nil)
}

func TestArtistUsecase_UpdateArtistName_FailedSelect(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_artist.NewMockArtistRep(ctrl)
	mockUsecase := usecase.NewArtistUsecase(mockRepo)

	expectedArtist := &models.Artist{
		ID:   1,
		Name: "Imagine Dragons",
	}

	mockRepo.
		EXPECT().
		SelectByName(gomock.Eq(expectedArtist.Name)).
		Return(expectedArtist, nil)

	err := mockUsecase.UpdateArtistName(expectedArtist)
	assert.Equal(t, err, NewErrorResponse(ErrNameAlreadyExist, nil))
}

func TestArtistUsecase_UpdateArtistName_FailedUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_artist.NewMockArtistRep(ctrl)
	mockUsecase := usecase.NewArtistUsecase(mockRepo)

	expectedArtist := &models.Artist{
		ID:   1,
		Name: "Imagine Dragons",
	}

	mockRepo.
		EXPECT().
		SelectByName(gomock.Eq(expectedArtist.Name)).
		Return(nil, sql.ErrNoRows)

	mockRepo.
		EXPECT().
		UpdateName(gomock.Eq(expectedArtist)).
		Return(sql.ErrNoRows)

	err := mockUsecase.UpdateArtistName(expectedArtist)
	assert.Equal(t, err, NewErrorResponse(ErrArtistNotExist, sql.ErrNoRows))
}

func TestArtistUsecase_UpdateArtistPoster(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_artist.NewMockArtistRep(ctrl)
	mockUsecase := usecase.NewArtistUsecase(mockRepo)

	expectedArtist := &models.Artist{
		ID:     1,
		Name:   "Imagine Dragons",
		Poster: "some poster",
	}

	mockRepo.
		EXPECT().
		UpdatePoster(gomock.Eq(expectedArtist)).
		Return(nil)

	err := mockUsecase.UpdateArtistPoster(expectedArtist)
	assert.Equal(t, err, nil)
}

func TestArtistUsecase_UpdateArtistPoster_Failed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_artist.NewMockArtistRep(ctrl)
	mockUsecase := usecase.NewArtistUsecase(mockRepo)

	expectedArtist := &models.Artist{
		ID:     1,
		Name:   "Imagine Dragons",
		Poster: "some poster",
	}

	mockRepo.
		EXPECT().
		UpdatePoster(gomock.Eq(expectedArtist)).
		Return(sql.ErrNoRows)

	err := mockUsecase.UpdateArtistPoster(expectedArtist)
	assert.Equal(t, err, NewErrorResponse(ErrArtistNotExist, sql.ErrNoRows))
}

func TestArtistUsecase_GetByParams(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_artist.NewMockArtistRep(ctrl)
	mockUsecase := usecase.NewArtistUsecase(mockRepo)

	expectedArtist1 := &models.Artist{
		ID:     1,
		Name:   "Imagine Dragons",
		Poster: "some poster",
	}

	expectedArtist2 := &models.Artist{
		ID:     2,
		Name:   "Imagine Dragons two",
		Poster: "some poster",
	}

	expectedArtists := []*models.Artist{expectedArtist1, expectedArtist2}

	count := uint64(2)
	from := uint64(0)

	mockRepo.
		EXPECT().
		SelectByParam(gomock.Eq(count), gomock.Eq(from)).
		Return(expectedArtists, nil)

	artists, err := mockUsecase.GetByParams(count, from)
	assert.Equal(t, err, nil)
	assert.Equal(t, artists, expectedArtists)
}

func TestArtistUsecase_GetByParams_Failed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_artist.NewMockArtistRep(ctrl)
	mockUsecase := usecase.NewArtistUsecase(mockRepo)

	count := uint64(2)
	from := uint64(100)

	mockRepo.
		EXPECT().
		SelectByParam(gomock.Eq(count), gomock.Eq(from)).
		Return(nil, sql.ErrNoRows)

	artists, err := mockUsecase.GetByParams(count, from)
	assert.Equal(t, err, nil)
	assert.Equal(t, artists, nil)
}
