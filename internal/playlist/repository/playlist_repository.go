package repository

import (
	"database/sql"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/playlist"
)

type PlaylistRep struct {
	dbConn *sql.DB
}

func NewPlaylistRep(dbConn *sql.DB) playlist.PlaylistRep {
	return &PlaylistRep{
		dbConn: dbConn,
	}
}

func (pr *PlaylistRep) Insert(playlist *models.Playlist) error {
	query := "insert into playlists(user_id, title) values($1, $2) returning id"

	err := pr.dbConn.QueryRow(query, playlist.UserID, playlist.Title).Scan(&playlist.ID)

	if err != nil {
		return err
	}

	return nil
}

func (pr *PlaylistRep) Update(playlist *models.Playlist) error {
	query := "update playlists set title = $1, poster = $2, is_public = $3 where id = $4 returning id"

	err := pr.dbConn.QueryRow(query, playlist.Title, playlist.Poster, playlist.IsPublic, playlist.ID).Scan(&playlist.ID)

	if err != nil {
		return err
	}

	return nil
}

func (pr *PlaylistRep) Delete(id uint64) error {
	query := "delete from playlists where id = $1 returning id"

	err := pr.dbConn.QueryRow(query, id).Scan(&id)

	if err != nil {
		return err
	}

	return nil
}

func (pr *PlaylistRep) SelectByID(id uint64) (*models.Playlist, error) {
	query := "select id, user_id, title, poster, is_public from playlists where id = $1"

	playlist := &models.Playlist{}

	err := pr.dbConn.QueryRow(query, id).
		Scan(&playlist.ID,
			&playlist.UserID,
			&playlist.Title,
			&playlist.Poster,
			&playlist.IsPublic)

	if err != nil {
		return nil, err
	}

	return playlist, nil
}

func (pr *PlaylistRep) SelectByUserID(userID uint64) ([]*models.Playlist, error) {
	query := "select id, user_id, title, poster, is_public from playlists where user_id = $1"

	playlists := []*models.Playlist{}

	rows, err := pr.dbConn.Query(query, userID)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		playlist := &models.Playlist{}

		err := rows.Scan(
			&playlist.ID,
			&playlist.UserID,
			&playlist.Title,
			&playlist.Poster,
			&playlist.IsPublic)

		if err != nil {
			return nil, err
		}

		playlists = append(playlists, playlist)
	}

	return playlists, nil
}

func (pr *PlaylistRep) InsertTrack(trackID uint64, playlistID uint64) error {
	query := "insert into track_playlist(track_id, playlist_id) values($1, $2) returning track_id"

	err := pr.dbConn.QueryRow(query, trackID, playlistID).Scan(&trackID)

	if err != nil {
		return err
	}

	return nil
}

func (pr *PlaylistRep) DeleteTrack(trackID uint64, playlistID uint64) error {
	query := "delete from track_playlist where track_id = $1 and playlist_id = $2 returning track_id"

	err := pr.dbConn.QueryRow(query, trackID, playlistID).Scan(&trackID)

	if err != nil {
		return err
	}

	return nil
}

func (pr *PlaylistRep) SelectPublicByUserID(userID uint64) ([]*models.Playlist, error) {
	query := "select id, user_id, title, poster, is_public from playlists where user_id = $1 and is_public = true"

	playlists := []*models.Playlist{}

	rows, err := pr.dbConn.Query(query, userID)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		playlist := &models.Playlist{}

		err := rows.Scan(
			&playlist.ID,
			&playlist.UserID,
			&playlist.Title,
			&playlist.Poster,
			&playlist.IsPublic)

		if err != nil {
			return nil, err
		}

		playlists = append(playlists, playlist)
	}

	return playlists, nil
}
