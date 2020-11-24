package usecase_test

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/album/mock_album"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/album/usecase"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/consts"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/error_response"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
)

func TestAlbumUsecase_CreateAlbum(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_album.NewMockAlbumRep(ctrl)
	mockUsecase := usecase.NewAlbumUsecase(mockRepo)

	album := &models.Album{}

	mockRepo.
		EXPECT().
		Insert(gomock.Eq(album)).
		Return(nil)

	err := mockUsecase.CreateAlbum(album)
	assert.Equal(t, err, nil)
}

func TestAlbumUsecase_CreateAlbum_Failed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_album.NewMockAlbumRep(ctrl)
	mockUsecase := usecase.NewAlbumUsecase(mockRepo)

	album := &models.Album{}

	mockRepo.
		EXPECT().
		Insert(gomock.Eq(album)).
		Return(sql.ErrTxDone)

	err := mockUsecase.CreateAlbum(album)
	assert.Equal(t, err, NewErrorResponse(ErrInternal, sql.ErrTxDone))
}

func TestAlbumUsecase_DeleteAlbum(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_album.NewMockAlbumRep(ctrl)
	mockUsecase := usecase.NewAlbumUsecase(mockRepo)

	id := uint64(5)

	mockRepo.
		EXPECT().
		Delete(gomock.Eq(id)).
		Return(nil)

	err := mockUsecase.DeleteAlbum(id)
	assert.Equal(t, err, nil)
}

func TestAlbumUsecase_DeleteAlbum_Failed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_album.NewMockAlbumRep(ctrl)
	mockUsecase := usecase.NewAlbumUsecase(mockRepo)

	id := uint64(5)

	mockRepo.
		EXPECT().
		Delete(gomock.Eq(id)).
		Return(sql.ErrNoRows)

	err := mockUsecase.DeleteAlbum(id)
	assert.Equal(t, err, NewErrorResponse(ErrAlbumNotExist, sql.ErrNoRows))
}

func TestAlbumUsecase_DeleteAlbum_Failed_Internal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_album.NewMockAlbumRep(ctrl)
	mockUsecase := usecase.NewAlbumUsecase(mockRepo)

	id := uint64(5)
	dbErr := errors.New("Some database err\n")

	mockRepo.
		EXPECT().
		Delete(gomock.Eq(id)).
		Return(dbErr)

	err := mockUsecase.DeleteAlbum(id)
	assert.Equal(t, err, NewErrorResponse(ErrInternal, dbErr))
}

func TestAlbumUsecase_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_album.NewMockAlbumRep(ctrl)
	mockUsecase := usecase.NewAlbumUsecase(mockRepo)

	id := uint64(5)

	expectedAlbum := &models.Album{
		ID:         5,
		Title:      "Some title",
		ArtistID:   0,
		ArtistName: "Some artist name",
		Poster:     "Some poster",
	}

	mockRepo.
		EXPECT().
		SelectByID(gomock.Eq(id)).
		Return(expectedAlbum, nil)

	album, err := mockUsecase.GetByID(id)
	assert.Equal(t, err, nil)
	assert.Equal(t, album, expectedAlbum)
}

func TestAlbumUsecase_GetByID_Failed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_album.NewMockAlbumRep(ctrl)
	mockUsecase := usecase.NewAlbumUsecase(mockRepo)

	id := uint64(5)

	mockRepo.
		EXPECT().
		SelectByID(gomock.Eq(id)).
		Return(nil, sql.ErrNoRows)

	album, err := mockUsecase.GetByID(id)
	assert.Equal(t, err, NewErrorResponse(ErrAlbumNotExist, sql.ErrNoRows))
	assert.Equal(t, album, nil)
}

func TestAlbumUsecase_GetByID_Failed_Internal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_album.NewMockAlbumRep(ctrl)
	mockUsecase := usecase.NewAlbumUsecase(mockRepo)

	id := uint64(5)
	dbErr := errors.New("Some database err\n")

	mockRepo.
		EXPECT().
		SelectByID(gomock.Eq(id)).
		Return(nil, dbErr)

	album, err := mockUsecase.GetByID(id)
	assert.Equal(t, err, NewErrorResponse(ErrInternal, dbErr))
	assert.Equal(t, album, nil)
}

func TestAlbumUsecase_GetByArtistID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_album.NewMockAlbumRep(ctrl)
	mockUsecase := usecase.NewAlbumUsecase(mockRepo)

	artistId := uint64(0)

	expectedAlbum1 := &models.Album{
		ID:         1,
		Title:      "Some title",
		ArtistID:   0,
		ArtistName: "Some artist name",
		Poster:     "Some poster",
	}

	expectedAlbum2 := &models.Album{
		ID:         2,
		Title:      "Some title",
		ArtistID:   0,
		ArtistName: "Some artist name",
		Poster:     "Some poster",
	}

	expectedAlbums := []*models.Album{expectedAlbum1, expectedAlbum2}

	mockRepo.
		EXPECT().
		SelectByArtistID(gomock.Eq(artistId)).
		Return(expectedAlbums, nil)

	album, err := mockUsecase.GetByArtistID(artistId)
	assert.Equal(t, err, nil)
	assert.Equal(t, album, expectedAlbums)
}

func TestAlbumUsecase_GetByArtistID_Failed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_album.NewMockAlbumRep(ctrl)
	mockUsecase := usecase.NewAlbumUsecase(mockRepo)

	artistId := uint64(5)

	mockRepo.
		EXPECT().
		SelectByArtistID(gomock.Eq(artistId)).
		Return(nil, sql.ErrNoRows)

	album, err := mockUsecase.GetByArtistID(artistId)
	assert.Equal(t, err, nil)
	assert.Equal(t, album, nil)
}

func TestAlbumUsecase_GetByArtistID_Failed_Internal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_album.NewMockAlbumRep(ctrl)
	mockUsecase := usecase.NewAlbumUsecase(mockRepo)

	artistId := uint64(5)
	dbErr := errors.New("Some database err\n")

	mockRepo.
		EXPECT().
		SelectByArtistID(gomock.Eq(artistId)).
		Return(nil, dbErr)

	album, err := mockUsecase.GetByArtistID(artistId)
	assert.Equal(t, err, NewErrorResponse(ErrInternal, dbErr))
	assert.Equal(t, album, nil)
}

func TestAlbumUsecase_UpdateAlbum(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_album.NewMockAlbumRep(ctrl)
	mockUsecase := usecase.NewAlbumUsecase(mockRepo)

	album := &models.Album{}

	mockRepo.
		EXPECT().
		Update(gomock.Eq(album)).
		Return(nil)

	err := mockUsecase.UpdateAlbum(album)
	assert.Equal(t, err, nil)
}

func TestAlbumUsecase_UpdateAlbum_Failed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_album.NewMockAlbumRep(ctrl)
	mockUsecase := usecase.NewAlbumUsecase(mockRepo)

	album := &models.Album{}

	mockRepo.
		EXPECT().
		Update(gomock.Eq(album)).
		Return(sql.ErrNoRows)

	err := mockUsecase.UpdateAlbum(album)
	assert.Equal(t, err, NewErrorResponse(ErrAlbumNotExist, sql.ErrNoRows))
}

func TestAlbumUsecase_UpdateAlbum_Internal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_album.NewMockAlbumRep(ctrl)
	mockUsecase := usecase.NewAlbumUsecase(mockRepo)

	album := &models.Album{}
	dbErr := errors.New("Some database err\n")

	mockRepo.
		EXPECT().
		Update(gomock.Eq(album)).
		Return(dbErr)

	err := mockUsecase.UpdateAlbum(album)
	assert.Equal(t, err, NewErrorResponse(ErrInternal, dbErr))
}

func TestAlbumUsecase_UpdateAlbumPoster(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_album.NewMockAlbumRep(ctrl)
	mockUsecase := usecase.NewAlbumUsecase(mockRepo)

	album := &models.Album{}

	mockRepo.
		EXPECT().
		UpdatePoster(gomock.Eq(album)).
		Return(nil)

	err := mockUsecase.UpdateAlbumPoster(album)
	assert.Equal(t, err, nil)
}

func TestAlbumUsecase_UpdateAlbumPoster_Failed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_album.NewMockAlbumRep(ctrl)
	mockUsecase := usecase.NewAlbumUsecase(mockRepo)

	album := &models.Album{}

	mockRepo.
		EXPECT().
		UpdatePoster(gomock.Eq(album)).
		Return(sql.ErrNoRows)

	err := mockUsecase.UpdateAlbumPoster(album)
	assert.Equal(t, err, NewErrorResponse(ErrAlbumNotExist, sql.ErrNoRows))
}

func TestAlbumUsecase_UpdateAlbumPoster_Internal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_album.NewMockAlbumRep(ctrl)
	mockUsecase := usecase.NewAlbumUsecase(mockRepo)

	album := &models.Album{}
	dbErr := errors.New("Some database err\n")

	mockRepo.
		EXPECT().
		UpdatePoster(gomock.Eq(album)).
		Return(dbErr)

	err := mockUsecase.UpdateAlbumPoster(album)
	assert.Equal(t, err, NewErrorResponse(ErrInternal, dbErr))
}

func TestAlbumUsecase_GetByParams(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_album.NewMockAlbumRep(ctrl)
	mockUsecase := usecase.NewAlbumUsecase(mockRepo)

	count, from := uint64(2), uint64(0)

	expectedAlbum1 := &models.Album{
		ID:         1,
		Title:      "Some title",
		ArtistID:   0,
		ArtistName: "Some artist name",
		Poster:     "Some poster",
	}

	expectedAlbum2 := &models.Album{
		ID:         2,
		Title:      "Some title",
		ArtistID:   0,
		ArtistName: "Some artist name",
		Poster:     "Some poster",
	}

	expectedAlbums := []*models.Album{expectedAlbum1, expectedAlbum2}

	mockRepo.
		EXPECT().
		SelectByParam(gomock.Eq(count), gomock.Eq(from)).
		Return(expectedAlbums, nil)

	album, err := mockUsecase.GetByParams(count, from)
	assert.Equal(t, err, nil)
	assert.Equal(t, album, expectedAlbums)
}

func TestAlbumUsecase_GetByParams_Failed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_album.NewMockAlbumRep(ctrl)
	mockUsecase := usecase.NewAlbumUsecase(mockRepo)

	count, from := uint64(2), uint64(0)

	mockRepo.
		EXPECT().
		SelectByParam(gomock.Eq(count), gomock.Eq(from)).
		Return(nil, nil)

	album, err := mockUsecase.GetByParams(count, from)
	assert.Equal(t, err, nil)
	assert.Equal(t, album, nil)
}

func TestAlbumUsecase_GetByParams_Failed_Internal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_album.NewMockAlbumRep(ctrl)
	mockUsecase := usecase.NewAlbumUsecase(mockRepo)

	count, from := uint64(2), uint64(0)
	dbErr := errors.New("Some database err\n")

	mockRepo.
		EXPECT().
		SelectByParam(gomock.Eq(count), gomock.Eq(from)).
		Return(nil, dbErr)

	album, err := mockUsecase.GetByParams(count, from)
	assert.Equal(t, err, NewErrorResponse(ErrInternal, dbErr))
	assert.Equal(t, album, nil)
}
