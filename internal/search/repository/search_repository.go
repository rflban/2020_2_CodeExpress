package repository

import (
	"database/sql"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/search"
)

type SearchRep struct {
	dbConn *sql.DB
}

func NewSearchRep(dbConn *sql.DB) search.SearchRep {
	return &SearchRep{
		dbConn: dbConn,
	}
}

func (sr *SearchRep) SelectAlbums(query string, offset uint64, limit uint64) ([]*models.Album, error) {
	rows, err := sr.dbConn.Query(`SELECT 
    	albums.id, 
       	albums.title, 
       	albums.artist_id, 
       	artists.name, 
       	albums.poster FROM albums 
       	    JOIN artists ON albums.artist_id = artists.id 
		WHERE albums.title ILIKE '%' || $1 || '%' COLLATE "C" 
		ORDER BY albums.title, artists.name 
		LIMIT $2 
		OFFSET $3;`,
		query, limit, offset)
	if err != nil {
		return nil, err
	}

	var albums []*models.Album
	for rows.Next() {
		album := &models.Album{}
		if err := rows.Scan(&album.ID, &album.Title, &album.ArtistID, &album.ArtistName, &album.Poster); err != nil {
			return nil, err
		}

		albums = append(albums, album)
	}

	return albums, nil
}

func (sr *SearchRep) SelectArtists(query string, offset uint64, limit uint64) ([]*models.Artist, error) {
	rows, err := sr.dbConn.Query(`SELECT 
       	artists.id, 
       	artists.name, 
      	artists.description, 
       	artists.poster, 
       	artists.avatar FROM artists 
		WHERE artists.name ILIKE '%' || $1 || '%' COLLATE "C" 
		ORDER BY artists.name 
		LIMIT $2 
		OFFSET $3;`,
		query, limit, offset)
	if err != nil {
		return nil, err
	}

	var artists []*models.Artist
	for rows.Next() {
		artist := &models.Artist{}
		if err := rows.Scan(&artist.ID, &artist.Name, &artist.Description, &artist.Poster, &artist.Avatar); err != nil {
			return nil, err
		}

		artists = append(artists, artist)
	}

	return artists, nil
}

func (sr *SearchRep) SelectTracks(query string, offset uint64, limit uint64, userId uint64) ([]*models.Track, error) {
	rows, err := sr.dbConn.Query(`SELECT 
       	tracks.id, 
       	tracks.album_id, 
       	albums.poster, 
        albums.artist_id, 
        artists.name, 
        tracks.title, 
        tracks.duration, 
        tracks.index, 
        tracks.audio,
       	user_track_like.track_id FROM tracks 
            JOIN albums ON tracks.album_id = albums.id 
            JOIN artists ON albums.artist_id = artists.id 
			LEFT JOIN user_track_like ON tracks.id = user_track_like.track_id AND user_track_like.user_id = $4 
		WHERE tracks.title ILIKE '%' || $1 || '%' COLLATE "C" 
		ORDER BY tracks.title, artists.name 
		LIMIT $2 
		OFFSET $3;`,
		query, limit, offset, userId)
	if err != nil {
		return nil, err
	}

	var tracks []*models.Track
	for rows.Next() {
		track := &models.Track{}
		var isLiked sql.NullInt64
		if err := rows.Scan(&track.ID, &track.AlbumID, &track.AlbumPoster, &track.ArtistID, &track.Artist, &track.Title,
			&track.Duration, &track.Index, &track.Audio, &isLiked); err != nil {
			return nil, err
		}

		if isLiked.Valid {
			track.IsLiked = true
		}

		tracks = append(tracks, track)
	}

	return tracks, nil
}

func (sr *SearchRep) SelectUsers(query string, offset uint64, limit uint64) ([]*models.User, error) {
	rows, err := sr.dbConn.Query(`SELECT
       	users.id,
       	users.name,
       	users.avatar FROM users
        WHERE users.name ILIKE '%' || $1 || '%' COLLATE "C" 
		ORDER BY users.name 
		LIMIT $2 
		OFFSET $3;`,
		query, limit, offset)
	if err != nil {
		return nil, err
	}

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		if err := rows.Scan(&user.ID, &user.Name, &user.Avatar); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}
