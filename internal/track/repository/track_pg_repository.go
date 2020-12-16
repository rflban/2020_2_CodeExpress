package repository

import (
	"database/sql"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/track"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
)

type TrackRep struct {
	dbConn *sql.DB
}

func NewTrackRep(dbConn *sql.DB) track.TrackRep {
	return &TrackRep{
		dbConn: dbConn,
	}
}

func (ar *TrackRep) Insert(track *models.Track) error {
	if err := ar.dbConn.QueryRow("INSERT INTO tracks (album_id, title) VALUES ($1, $2) RETURNING id, index;",
		track.AlbumID, track.Title).Scan(&track.ID, &track.Index); err != nil {
		return err
	}
	return nil
}

func (ar *TrackRep) Update(track *models.Track) error {
	if err := ar.dbConn.QueryRow("UPDATE tracks SET album_id = $1, title = $2, index = $3 WHERE id = $4 RETURNING id;",
		track.AlbumID, track.Title, track.Index, track.ID).Scan(&track.ID); err != nil {
		return err
	}
	return nil
}

func (ar *TrackRep) UpdateAudio(track *models.Track) error {
	if err := ar.dbConn.QueryRow("UPDATE tracks SET audio = $1, duration = $2 WHERE id = $3 RETURNING id;",
		track.Audio, track.Duration, track.ID).Scan(&track.ID); err != nil {
		return err
	}
	return nil
}

func (ar *TrackRep) Delete(id uint64) error {
	if err := ar.dbConn.QueryRow("DELETE FROM tracks WHERE id = $1 RETURNING id;", id).Scan(&id); err != nil {
		return err
	}
	return nil
}

func (ar *TrackRep) SelectByID(id, userId uint64) (*models.Track, error) {
	track := &models.Track{}
	var isLiked sql.NullBool
	if err := ar.dbConn.QueryRow(`SELECT 
	tracks.id, 
	tracks.album_id, 
	tracks.title, 
	tracks.duration, 
	tracks.index, 
	tracks.audio, 
	albums.poster, 
	artists.name, 
	artists.id, 
    user_track_like.is_like FROM tracks
	    JOIN albums ON tracks.album_id = albums.id
	    JOIN artists ON albums.artist_id = artists.id
		LEFT JOIN user_track_like ON tracks.id = user_track_like.track_id AND user_track_like.user_id = $2
	WHERE tracks.id = $1`, id, userId).Scan(&track.ID, &track.AlbumID, &track.Title, &track.Duration, &track.Index,
		&track.Audio, &track.AlbumPoster, &track.Artist, &track.ArtistID, &isLiked); err != nil {
		return nil, err
	}

	if isLiked.Valid {
		track.IsLiked = isLiked.Bool
	}

	return track, nil
}

func (ar *TrackRep) SelectByArtistId(artistId, userId uint64) ([]*models.Track, error) {
	rows, err := ar.dbConn.Query(`SELECT 
	tracks.id, 
	tracks.album_id, 
	tracks.title, 
	tracks.duration, 
	tracks.index, 
	tracks.audio, 
	albums.poster, 
	artists.id, 
	artists.name, 
	user_track.user_id, 
	user_track_like.is_like FROM tracks 
		JOIN albums ON tracks.album_id = albums.id 
		JOIN artists ON albums.artist_id = artists.id 
		LEFT JOIN user_track ON tracks.id = user_track.track_id AND user_track.user_id = $2 
		LEFT JOIN user_track_like ON tracks.id = user_track_like.track_id AND user_track_like.user_id = $2 
	WHERE artists.id = $1`, artistId, userId)
	if err != nil {
		return nil, err
	}

	tracks := []*models.Track{}
	for rows.Next() {
		track := &models.Track{}
		var userFavouriteId sql.NullInt64
		var isLiked sql.NullBool
		if err := rows.Scan(&track.ID, &track.AlbumID, &track.Title, &track.Duration, &track.Index, &track.Audio,
			&track.AlbumPoster, &track.ArtistID, &track.Artist, &userFavouriteId, &isLiked); err != nil {
			return nil, err
		}

		if userFavouriteId.Valid {
			track.IsFavorite = true
		}
		if isLiked.Valid {
			track.IsLiked = isLiked.Bool
		}

		tracks = append(tracks, track)
	}

	return tracks, nil
}

func (ar *TrackRep) SelectByAlbumID(albumID, userId uint64) ([]*models.Track, error) {
	rows, err := ar.dbConn.Query(`SELECT 
	tracks.id, 
	tracks.album_id, 
	tracks.title, 
	tracks.duration, 
	tracks.index, 
	tracks.audio, 
	albums.poster, 
	artists.name, 
	albums.artist_id, 
    user_track_like.is_like FROM tracks 
	    JOIN albums ON tracks.album_id = albums.id 
	    JOIN artists ON albums.artist_id = artists.id 
		LEFT JOIN user_track_like ON tracks.id = user_track_like.track_id AND user_track_like.user_id = $2 
	WHERE albums.id = $1`, albumID, userId)
	if err != nil {
		return nil, err
	}

	tracks := []*models.Track{}
	for rows.Next() {
		track := &models.Track{}
		var isLiked sql.NullBool
		if err := rows.Scan(&track.ID, &track.AlbumID, &track.Title, &track.Duration, &track.Index, &track.Audio,
			&track.AlbumPoster, &track.Artist, &track.ArtistID, &isLiked); err != nil {
			return nil, err
		}

		if isLiked.Valid {
			track.IsLiked = isLiked.Bool
		}

		tracks = append(tracks, track)
	}

	return tracks, nil
}

func (ar *TrackRep) SelectByParams(count, from, userId uint64) ([]*models.Track, error) {
	rows, err := ar.dbConn.Query(`SELECT 
	tracks.id, 
	tracks.album_id, 
	tracks.title, 
	tracks.duration, 
	tracks.index, 
	tracks.audio, 
	albums.poster, 
	artists.id, 
	artists.name, 
	user_track.user_id, 
    user_track_like.is_like FROM tracks 
		JOIN albums ON tracks.album_id = albums.id 
		JOIN artists ON albums.artist_id = artists.id 
		LEFT JOIN user_track ON tracks.id = user_track.track_id AND user_track.user_id = $3
		LEFT JOIN user_track_like ON tracks.id = user_track_like.track_id AND user_track_like.user_id = $3
	ORDER BY tracks.title LIMIT $1 OFFSET $2`,
		count, from, userId)
	if err != nil {
		return nil, err
	}

	tracks := []*models.Track{}
	for rows.Next() {
		track := &models.Track{}
		var userFavouriteId sql.NullInt64
		var isLiked sql.NullBool
		if err := rows.Scan(&track.ID, &track.AlbumID, &track.Title, &track.Duration, &track.Index, &track.Audio,
			&track.AlbumPoster, &track.ArtistID, &track.Artist, &userFavouriteId, &isLiked); err != nil {
			return nil, err
		}

		if userFavouriteId.Valid {
			track.IsFavorite = true
		}
		if isLiked.Valid {
			track.IsLiked = isLiked.Bool
		}

		tracks = append(tracks, track)
	}

	return tracks, nil
}

func (ar *TrackRep) SelectFavoritesByUserID(userID uint64) ([]*models.Track, error) {
	rows, err := ar.dbConn.Query(`SELECT 
	tracks.id, 
	tracks.album_id, 
	tracks.title, 
	tracks.duration, 
	tracks.index, 
	tracks.audio, 
	albums.poster, 
	artists.id, 
	artists.name, 
    user_track_like.is_like FROM user_track 
		JOIN tracks ON user_track.track_id = tracks.id 
		JOIN albums ON tracks.album_id = albums.id 
		JOIN artists ON albums.artist_id = artists.id 
		LEFT JOIN user_track_like ON tracks.id = user_track_like.track_id AND user_track_like.user_id = $1 
	WHERE user_track.user_id = $1`, userID)
	if err != nil {
		return nil, err
	}

	tracks := []*models.Track{}
	for rows.Next() {
		track := &models.Track{}
		var isLiked sql.NullBool
		if err := rows.Scan(&track.ID, &track.AlbumID, &track.Title, &track.Duration, &track.Index, &track.Audio,
			&track.AlbumPoster, &track.ArtistID, &track.Artist, &isLiked); err != nil {
			return nil, err
		}

		track.IsFavorite = true
		if isLiked.Valid {
			track.IsLiked = isLiked.Bool
		}

		tracks = append(tracks, track)
	}

	return tracks, nil
}

func (ar *TrackRep) InsertTrackToUser(userID, trackID uint64) error {
	if err := ar.dbConn.QueryRow("INSERT INTO user_track (user_id, track_id) VALUES ($1, $2) RETURNING track_id;", //TODO: здесь достаточно dbConn.Exec без RETURNING?...
		userID, trackID).Scan(&trackID); err != nil {
		return err
	}
	return nil
}

func (ar *TrackRep) DeleteTrackFromUsersTracks(userID, trackID uint64) error {
	if err := ar.dbConn.QueryRow("DELETE FROM user_track WHERE user_id = $1 AND track_id = $2 RETURNING track_id;", //TODO: здесь достаточно dbConn.Exec без RETURNING?...
		userID, trackID).Scan(&trackID); err != nil {
		return err
	}
	return nil
}

func (ar *TrackRep) SelectByPlaylistID(playlistID, userId uint64) ([]*models.Track, error) {
	rows, err := ar.dbConn.Query(`SELECT 
	tracks.id, 
	tracks.album_id, 
	tracks.title, 
	tracks.duration, 
	tracks.index, 
	tracks.audio, 
	albums.poster, 
	artists.name, 
	albums.artist_id, 
    user_track_like.is_like FROM track_playlist 
		JOIN tracks ON track_playlist.track_id = tracks.id 
		JOIN albums ON tracks.album_id = albums.id 
		JOIN artists ON albums.artist_id = artists.id 
		LEFT JOIN user_track_like ON tracks.id = user_track_like.track_id AND user_track_like.user_id = $2
	WHERE track_playlist.playlist_id = $1`, playlistID, userId)
	if err != nil {
		return nil, err
	}

	tracks := []*models.Track{}
	for rows.Next() {
		track := &models.Track{}
		var isLiked sql.NullBool
		if err := rows.Scan(&track.ID, &track.AlbumID, &track.Title, &track.Duration, &track.Index, &track.Audio,
			&track.AlbumPoster, &track.Artist, &track.ArtistID, &isLiked); err != nil {
			return nil, err
		}

		if isLiked.Valid {
			track.IsLiked = isLiked.Bool
		}

		tracks = append(tracks, track)
	}

	return tracks, nil
}

func (ar *TrackRep) LikeOrDislikeTrack(userId, trackId uint64, like bool) error {
	if _, err := ar.dbConn.Exec(`INSERT INTO user_track_like (user_id, track_id, is_like)
	VALUES ($1, $2, $3)
	ON CONFLICT (user_id, track_id) DO UPDATE SET is_like = excluded.is_like;`,
		userId, trackId, like); err != nil {
		return err
	}
	return nil
}
