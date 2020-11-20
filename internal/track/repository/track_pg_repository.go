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
	query := "insert into tracks(album_id, title) values($1, $2) returning id, index"

	err := ar.dbConn.QueryRow(query, track.AlbumID, track.Title).Scan(&track.ID, &track.Index)

	if err != nil {
		return err
	}

	return nil
}

func (ar *TrackRep) Update(track *models.Track) error {
	query := "update tracks set album_id = $1, title = $2, index = $3 where id = $4 returning id"

	err := ar.dbConn.QueryRow(query, track.AlbumID, track.Title, track.Index, track.ID).Scan(&track.ID)

	if err != nil {
		return err
	}

	return nil
}

func (ar *TrackRep) UpdateAudio(track *models.Track) error {
	query := "update tracks set audio = $1, duration = $2 where id = $3 returning id"

	err := ar.dbConn.QueryRow(query, track.Audio, track.Duration, track.ID).Scan(&track.ID)

	if err != nil {
		return err
	}

	return nil
}

func (ar *TrackRep) Delete(id uint64) error {
	query := "delete from tracks where id = $1 returning id"

	err := ar.dbConn.QueryRow(query, id).Scan(&id)

	if err != nil {
		return err
	}

	return nil
}

func (ar *TrackRep) SelectByID(id uint64) (*models.Track, error) {
	query := `
	select 
	t.id, 
	t.album_id, 
	t.title, 
	t.duration, 
	t.index, 
	t.audio,
	a.poster,
	ar.name,
	a.artist_id
	from tracks as t join albums a on t.album_id = a.id join artists ar on a.artist_id = ar.id where t.id = $1`

	track := &models.Track{}

	err := ar.dbConn.QueryRow(query, id).
		Scan(&track.ID,
			&track.AlbumID,
			&track.Title,
			&track.Duration,
			&track.Index,
			&track.Audio,
			&track.AlbumPoster,
			&track.Artist,
			&track.ArtistID)

	if err != nil {
		return nil, err
	}

	return track, nil
}

func (ar *TrackRep) SelectByArtistID(artistID uint64) ([]*models.Track, error) {
	query := `
	select 
	t.id, 
	t.album_id, 
	t.title, 
	t.duration, 
	t.index, 
	t.audio,
	a.poster,
	ar.name,
	a.artist_id
	from tracks as t join albums a on t.album_id = a.id join artists ar on a.artist_id = ar.id where a.artist_id = $1`

	tracks := []*models.Track{}

	rows, err := ar.dbConn.Query(query, artistID)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		track := &models.Track{}
		err := rows.
			Scan(&track.ID,
				&track.AlbumID,
				&track.Title,
				&track.Duration,
				&track.Index,
				&track.Audio,
				&track.AlbumPoster,
				&track.Artist,
				&track.ArtistID)

		if err != nil {
			return nil, err
		}

		tracks = append(tracks, track)
	}

	return tracks, nil
}

func (ar *TrackRep) SelectByAlbumID(albumID uint64) ([]*models.Track, error) {
	query := `
	select 
	t.id, 
	t.album_id, 
	t.title, 
	t.duration, 
	t.index, 
	t.audio,
	a.poster,
	ar.name,
	a.artist_id
	from tracks as t join albums a on t.album_id = a.id join artists ar on a.artist_id = ar.id where a.id = $1`

	tracks := []*models.Track{}

	rows, err := ar.dbConn.Query(query, albumID)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		track := &models.Track{}
		err := rows.
			Scan(&track.ID,
				&track.AlbumID,
				&track.Title,
				&track.Duration,
				&track.Index,
				&track.Audio,
				&track.AlbumPoster,
				&track.Artist,
				&track.ArtistID)

		if err != nil {
			return nil, err
		}

		tracks = append(tracks, track)
	}

	return tracks, nil
}

func (ar *TrackRep) SelectByParam(count uint64, from uint64) ([]*models.Track, error) {
	query := `
	select 
	t.id, 
	t.album_id, 
	t.title, 
	t.duration, 
	t.index, 
	t.audio,
	a.poster,
	ar.name,
	a.artist_id
	from tracks as t join albums a on t.album_id = a.id join artists ar on a.artist_id = ar.id order by t.title limit $1 offset $2`

	tracks := []*models.Track{}

	rows, err := ar.dbConn.Query(query, count, from)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		track := &models.Track{}
		err := rows.Scan(&track.ID,
			&track.AlbumID,
			&track.Title,
			&track.Duration,
			&track.Index,
			&track.Audio,
			&track.AlbumPoster,
			&track.Artist,
			&track.ArtistID)
		if err != nil {
			return nil, err
		}

		tracks = append(tracks, track)
	}

	return tracks, nil
}

func (ar *TrackRep) SelectFavoritesByUserID(userID uint64) ([]*models.Track, error) {
	query := `
	select 
	t.id, 
	t.album_id, 
	t.title, 
	t.duration,
	t.index, 
	t.audio, 
	a.poster, 
	ar.name, 
	a.artist_id 
	from user_track as u join tracks as t on u.track_id = t.id join albums a on t.album_id = a.id join artists ar on a.artist_id = ar.id where u.user_id = $1`

	tracks := []*models.Track{}

	rows, err := ar.dbConn.Query(query, userID)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		track := &models.Track{}
		err := rows.
			Scan(&track.ID,
				&track.AlbumID,
				&track.Title,
				&track.Duration,
				&track.Index,
				&track.Audio,
				&track.AlbumPoster,
				&track.Artist,
				&track.ArtistID)

		if err != nil {
			return nil, err
		}

		tracks = append(tracks, track)
	}

	return tracks, nil
}

func (ar *TrackRep) InsertTrackToUser(userID, trackID uint64) error {
	query := "insert into user_track(user_id, track_id) values($1, $2) returning track_id"

	err := ar.dbConn.QueryRow(query, userID, trackID).Scan(&trackID)

	if err != nil {
		return err
	}

	return nil
}

func (ar *TrackRep) DeleteTrackFromUsersTracks(userID, trackID uint64) error {
	query := "delete from user_track where user_id = $1 and track_id = $2 returning track_id"

	err := ar.dbConn.QueryRow(query, userID, trackID).Scan(&trackID)

	if err != nil {
		return err
	}

	return nil
}

func (ar *TrackRep) SelectByPlaylistID(playlistID uint64) ([]*models.Track, error) {
	query := `
	select 
	t.id, 
	t.album_id, 
	t.title, 
	t.duration,
	t.index, 
	t.audio, 
	a.poster, 
	ar.name, 
	a.artist_id 
	from track_playlist as tp 
	join tracks as t on tp.track_id = t.id 
	join albums a on t.album_id = a.id 
	join artists ar on a.artist_id = ar.id where tp.playlist_id = $1`

	tracks := []*models.Track{}

	rows, err := ar.dbConn.Query(query, playlistID)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		track := &models.Track{}
		err := rows.
			Scan(&track.ID,
				&track.AlbumID,
				&track.Title,
				&track.Duration,
				&track.Index,
				&track.Audio,
				&track.AlbumPoster,
				&track.Artist,
				&track.ArtistID)

		if err != nil {
			return nil, err
		}

		tracks = append(tracks, track)
	}

	return tracks, nil
}
