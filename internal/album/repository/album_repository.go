package repository

import (
	"database/sql"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/album"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
)

type AlbumRep struct {
	dbConn *sql.DB
}

func NewAlbumRep(dbConn *sql.DB) album.AlbumRep {
	return &AlbumRep{
		dbConn: dbConn,
	}
}

func (ar *AlbumRep) Insert(album *models.Album) error {
	query := "insert into albums(artist_id, title) values($1, $2) returning id"

	err := ar.dbConn.QueryRow(query, album.ArtistID, album.Title).Scan(&album.ID)

	if err != nil {
		return err
	}

	return nil
}

func (ar *AlbumRep) Update(album *models.Album) error {
	query := "update albums set title = $1, artist_id = $2 where id = $3 returning id"

	err := ar.dbConn.QueryRow(query, album.Title, album.ArtistID, album.ID).Scan(&album.ID)

	if err != nil {
		return err
	}

	return nil
}

func (ar *AlbumRep) UpdatePoster(album *models.Album) error {
	query := "update albums set poster = $1 where id = $2 returning id"

	err := ar.dbConn.QueryRow(query, album.Poster, album.ID).Scan(&album.ID)

	if err != nil {
		return err
	}

	return nil
}

func (ar *AlbumRep) Delete(id uint64) error {
	query := "delete from albums where id = $1 returning id"

	err := ar.dbConn.QueryRow(query, id).Scan(&id)

	if err != nil {
		return err
	}

	return nil
}

func (ar *AlbumRep) SelectByID(id uint64) (*models.Album, error) {
	query := `select 
	al.id, 
	al.artist_id, 
	al.title, 
	al.poster, 
	a.name 
	from albums al join artists a on al.artist_id = a.id where al.id = $1`

	album := &models.Album{}

	err := ar.dbConn.QueryRow(query, id).
		Scan(&album.ID,
			&album.ArtistID,
			&album.Title,
			&album.Poster,
			&album.ArtistName)

	if err != nil {
		return nil, err
	}

	return album, nil
}

func (ar *AlbumRep) SelectByArtistID(artistID uint64) ([]*models.Album, error) {
	query := `select 
	al.id, 
	al.artist_id, 
	al.title, 
	al.poster, 
	a.name 
	from albums al join artists a on al.artist_id = a.id where al.artist_id = $1
	ORDER BY al.title`

	albums := []*models.Album{}

	rows, err := ar.dbConn.Query(query, artistID)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		album := &models.Album{}
		err := rows.
			Scan(&album.ID,
				&album.ArtistID,
				&album.Title,
				&album.Poster,
				&album.ArtistName)

		if err != nil {
			return nil, err
		}

		albums = append(albums, album)
	}

	return albums, nil
}

func (ar *AlbumRep) SelectByParam(count uint64, from uint64) ([]*models.Album, error) {
	query := `select 
	al.id, 
	al.artist_id, 
	al.title, 
	al.poster, 
	a.name 
	from albums al join artists a on al.artist_id = a.id
	ORDER BY a.name, al.title
	limit $1 offset $2`

	albums := []*models.Album{}

	rows, err := ar.dbConn.Query(query, count, from)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		album := &models.Album{}
		err := rows.
			Scan(&album.ID,
				&album.ArtistID,
				&album.Title,
				&album.Poster,
				&album.ArtistName)

		if err != nil {
			return nil, err
		}

		albums = append(albums, album)
	}

	return albums, nil
}
