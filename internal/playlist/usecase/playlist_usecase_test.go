package usecase_test

import (
	"database/sql"
	"errors"
	"testing"

	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/consts"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/playlist/mock_playlist"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/playlist/usecase"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/error_response"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
)

func TestPlaylistUsecase_CreatePlaylist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_playlist.NewMockPlaylistRep(ctrl)
	mockUsecase := usecase.NewPlaylistUsecase(mockRepo)

	playlist := &models.Playlist{}

	mockRepo.
		EXPECT().
		Insert(gomock.Eq(playlist)).
		Return(nil)

	err := mockUsecase.CreatePlaylist(playlist)
	assert.Equal(t, err, nil)
}

func TestPlaylistUsecase_CreatePlaylist_Failed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_playlist.NewMockPlaylistRep(ctrl)
	mockUsecase := usecase.NewPlaylistUsecase(mockRepo)

	playlist := &models.Playlist{}

	mockRepo.
		EXPECT().
		Insert(gomock.Eq(playlist)).
		Return(sql.ErrTxDone)

	err := mockUsecase.CreatePlaylist(playlist)
	assert.Equal(t, err, NewErrorResponse(ErrInternal, sql.ErrTxDone))
}

func TestPlaylistUsecase_DeletePlaylist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_playlist.NewMockPlaylistRep(ctrl)
	mockUsecase := usecase.NewPlaylistUsecase(mockRepo)

	id := uint64(5)

	mockRepo.
		EXPECT().
		Delete(gomock.Eq(id)).
		Return(nil)

	err := mockUsecase.DeletePlaylist(id)
	assert.Equal(t, err, nil)
}

func TestPlaylistUsecase_DeletePlaylist_Failed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_playlist.NewMockPlaylistRep(ctrl)
	mockUsecase := usecase.NewPlaylistUsecase(mockRepo)

	id := uint64(5)

	mockRepo.
		EXPECT().
		Delete(gomock.Eq(id)).
		Return(sql.ErrNoRows)

	err := mockUsecase.DeletePlaylist(id)
	assert.Equal(t, err, NewErrorResponse(ErrPlaylistNotExist, sql.ErrNoRows))
}

func TestPlaylistUsecase_DeletePlaylist_Failed_Internal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_playlist.NewMockPlaylistRep(ctrl)
	mockUsecase := usecase.NewPlaylistUsecase(mockRepo)

	id := uint64(5)
	dbErr := errors.New("Some database err")

	mockRepo.
		EXPECT().
		Delete(gomock.Eq(id)).
		Return(dbErr)

	err := mockUsecase.DeletePlaylist(id)
	assert.Equal(t, err, NewErrorResponse(ErrInternal, dbErr))
}

func TestPlaylistUsecase_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_playlist.NewMockPlaylistRep(ctrl)
	mockUsecase := usecase.NewPlaylistUsecase(mockRepo)

	id := uint64(5)

	expectedPlaylist := &models.Playlist{
		ID:     5,
		Title:  "Some title",
		UserID: 0,
		Poster: "Some poster",
	}

	mockRepo.
		EXPECT().
		SelectByID(gomock.Eq(id)).
		Return(expectedPlaylist, nil)

	playlist, err := mockUsecase.GetByID(id)
	assert.Equal(t, err, nil)
	assert.Equal(t, playlist, expectedPlaylist)
}

func TestPlaylistUsecase_GetByID_Failed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_playlist.NewMockPlaylistRep(ctrl)
	mockUsecase := usecase.NewPlaylistUsecase(mockRepo)

	id := uint64(5)

	mockRepo.
		EXPECT().
		SelectByID(gomock.Eq(id)).
		Return(nil, sql.ErrNoRows)

	playlist, err := mockUsecase.GetByID(id)
	assert.Equal(t, err, NewErrorResponse(ErrPlaylistNotExist, sql.ErrNoRows))
	assert.Equal(t, playlist, nil)
}

func TestPlaylistUsecase_GetByID_Failed_Internal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_playlist.NewMockPlaylistRep(ctrl)
	mockUsecase := usecase.NewPlaylistUsecase(mockRepo)

	id := uint64(5)
	dbErr := errors.New("Some database err")

	mockRepo.
		EXPECT().
		SelectByID(gomock.Eq(id)).
		Return(nil, dbErr)

	playlist, err := mockUsecase.GetByID(id)
	assert.Equal(t, err, NewErrorResponse(ErrInternal, dbErr))
	assert.Equal(t, playlist, nil)
}

func TestPlaylistUsecase_GetByUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_playlist.NewMockPlaylistRep(ctrl)
	mockUsecase := usecase.NewPlaylistUsecase(mockRepo)

	user_id := uint64(0)

	expectedPlaylist1 := &models.Playlist{
		ID:     1,
		Title:  "Some title",
		UserID: 0,
		Poster: "Some poster",
	}

	expectedPlaylist2 := &models.Playlist{
		ID:     2,
		Title:  "Some title",
		UserID: 0,
		Poster: "Some poster",
	}

	expectedPlaylists := []*models.Playlist{expectedPlaylist1, expectedPlaylist2}

	mockRepo.
		EXPECT().
		SelectByUserID(gomock.Eq(user_id)).
		Return(expectedPlaylists, nil)

	playlists, err := mockUsecase.GetByUserID(user_id)
	assert.Equal(t, err, nil)
	assert.Equal(t, playlists, expectedPlaylists)
}

func TestPlaylistUsecase_GetByUserID_Failed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_playlist.NewMockPlaylistRep(ctrl)
	mockUsecase := usecase.NewPlaylistUsecase(mockRepo)

	user_id := uint64(5)

	mockRepo.
		EXPECT().
		SelectByUserID(gomock.Eq(user_id)).
		Return(nil, sql.ErrNoRows)

	playlist, err := mockUsecase.GetByUserID(user_id)
	assert.Equal(t, err, nil)
	assert.Equal(t, playlist, nil)
}

func TestPlaylistUsecase_GetByUserID_Failed_Internal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_playlist.NewMockPlaylistRep(ctrl)
	mockUsecase := usecase.NewPlaylistUsecase(mockRepo)

	user_id := uint64(5)
	dbErr := errors.New("Some database err")

	mockRepo.
		EXPECT().
		SelectByUserID(gomock.Eq(user_id)).
		Return(nil, dbErr)

	playlists, err := mockUsecase.GetByUserID(user_id)
	assert.Equal(t, err, NewErrorResponse(ErrInternal, dbErr))
	assert.Equal(t, playlists, nil)
}

func TestPlaylistUsecase_UpdatePlaylist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_playlist.NewMockPlaylistRep(ctrl)
	mockUsecase := usecase.NewPlaylistUsecase(mockRepo)

	playlist := &models.Playlist{}

	mockRepo.
		EXPECT().
		Update(gomock.Eq(playlist)).
		Return(nil)

	err := mockUsecase.UpdatePlaylist(playlist)
	assert.Equal(t, err, nil)
}

func TestPlaylistUsecase_UpdatePlaylist_Failed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_playlist.NewMockPlaylistRep(ctrl)
	mockUsecase := usecase.NewPlaylistUsecase(mockRepo)

	playlist := &models.Playlist{}

	mockRepo.
		EXPECT().
		Update(gomock.Eq(playlist)).
		Return(sql.ErrNoRows)

	err := mockUsecase.UpdatePlaylist(playlist)
	assert.Equal(t, err, NewErrorResponse(ErrPlaylistNotExist, sql.ErrNoRows))
}

func TestPlaylistUsecase_UpdatePlaylist_Internal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_playlist.NewMockPlaylistRep(ctrl)
	mockUsecase := usecase.NewPlaylistUsecase(mockRepo)

	playlist := &models.Playlist{}
	dbErr := errors.New("Some database err")

	mockRepo.
		EXPECT().
		Update(gomock.Eq(playlist)).
		Return(dbErr)

	err := mockUsecase.UpdatePlaylist(playlist)
	assert.Equal(t, err, NewErrorResponse(ErrInternal, dbErr))
}

func TestPlaylistUsecase_AddTrack_Internal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_playlist.NewMockPlaylistRep(ctrl)
	mockUsecase := usecase.NewPlaylistUsecase(mockRepo)

	playlistID := uint64(5)
	trackID := uint64(42)
	dbErr := errors.New("Some database err")

	mockRepo.
		EXPECT().
		InsertTrack(gomock.Eq(trackID), gomock.Eq(playlistID)).
		Return(dbErr)

	err := mockUsecase.AddTrack(trackID, playlistID)
	assert.Equal(t, err, NewErrorResponse(ErrInternal, dbErr))
}

func TestPlaylistUsecase_AddTrack(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_playlist.NewMockPlaylistRep(ctrl)
	mockUsecase := usecase.NewPlaylistUsecase(mockRepo)

	playlistID := uint64(5)
	trackID := uint64(42)

	mockRepo.
		EXPECT().
		InsertTrack(gomock.Eq(trackID), gomock.Eq(playlistID)).
		Return(nil)

	err := mockUsecase.AddTrack(trackID, playlistID)
	assert.Equal(t, err, nil)
}

func TestPlaylistUsecase_DeleteTrack_Internal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_playlist.NewMockPlaylistRep(ctrl)
	mockUsecase := usecase.NewPlaylistUsecase(mockRepo)

	playlistID := uint64(5)
	trackID := uint64(42)
	dbErr := errors.New("Some database err")

	mockRepo.
		EXPECT().
		DeleteTrack(gomock.Eq(trackID), gomock.Eq(playlistID)).
		Return(dbErr)

	err := mockUsecase.DeleteTrack(trackID, playlistID)
	assert.Equal(t, err, NewErrorResponse(ErrInternal, dbErr))
}

func TestPlaylistUsecase_DeleteTrack(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_playlist.NewMockPlaylistRep(ctrl)
	mockUsecase := usecase.NewPlaylistUsecase(mockRepo)

	playlistID := uint64(5)
	trackID := uint64(42)

	mockRepo.
		EXPECT().
		DeleteTrack(gomock.Eq(trackID), gomock.Eq(playlistID)).
		Return(nil)

	err := mockUsecase.DeleteTrack(trackID, playlistID)
	assert.Equal(t, err, nil)
}

func TestPlaylistUsecase_GetPublicByUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_playlist.NewMockPlaylistRep(ctrl)
	mockUsecase := usecase.NewPlaylistUsecase(mockRepo)

	user_id := uint64(0)

	expectedPlaylist1 := &models.Playlist{
		ID:     1,
		Title:  "Some title",
		UserID: 0,
		Poster: "Some poster",
	}

	expectedPlaylist2 := &models.Playlist{
		ID:     2,
		Title:  "Some title",
		UserID: 0,
		Poster: "Some poster",
	}

	expectedPlaylists := []*models.Playlist{expectedPlaylist1, expectedPlaylist2}

	mockRepo.
		EXPECT().
		SelectPublicByUserID(gomock.Eq(user_id)).
		Return(expectedPlaylists, nil)

	playlists, err := mockUsecase.GetPublicByUserID(user_id)
	assert.Equal(t, err, nil)
	assert.Equal(t, playlists, expectedPlaylists)
}

func TestPlaylistUsecase_GetPublicByUserID_Failed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_playlist.NewMockPlaylistRep(ctrl)
	mockUsecase := usecase.NewPlaylistUsecase(mockRepo)

	user_id := uint64(5)

	mockRepo.
		EXPECT().
		SelectPublicByUserID(gomock.Eq(user_id)).
		Return(nil, sql.ErrNoRows)

	playlist, err := mockUsecase.GetPublicByUserID(user_id)
	assert.Equal(t, err, nil)
	assert.Equal(t, playlist, nil)
}

func TestPlaylistUsecase_GetPublicByUserID_Failed_Internal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_playlist.NewMockPlaylistRep(ctrl)
	mockUsecase := usecase.NewPlaylistUsecase(mockRepo)

	user_id := uint64(5)
	dbErr := errors.New("Some database err")

	mockRepo.
		EXPECT().
		SelectPublicByUserID(gomock.Eq(user_id)).
		Return(nil, dbErr)

	playlists, err := mockUsecase.GetPublicByUserID(user_id)
	assert.Equal(t, err, NewErrorResponse(ErrInternal, dbErr))
	assert.Equal(t, playlists, nil)
}
