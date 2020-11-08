package repository

import (
	"database/sql"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/artist"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
)

type ArtistRep struct {
	dbConn *sql.DB
}

func NewArtistRep(dbConn *sql.DB) artist.ArtistRep {
	return &ArtistRep{
		dbConn: dbConn,
	}
}

func (ar *ArtistRep) Insert(artist *models.Artist) error {
	query := "insert into artists(name, description) values($1, $2) returning id"

	err := ar.dbConn.QueryRow(query, artist.Name, artist.Description).Scan(&artist.ID)

	if err != nil {
		return err
	}

	return nil
}

func (ar *ArtistRep) Update(artist *models.Artist) error {
	query := "update artists set name = $1, poster = $2, avatar = $3, description = $4 where id = $5 returning id"

	err := ar.dbConn.QueryRow(query,
		artist.Name,
		artist.Poster,
		artist.Avatar,
		artist.Description,
		artist.ID).
		Scan(&artist.ID)

	if err != nil {
		return err
	}

	return nil
}

func (ar *ArtistRep) Delete(id uint64) error {
	query := "delete from artists where id = $1 returning id"

	err := ar.dbConn.QueryRow(query, id).Scan(&id)

	if err != nil {
		return err
	}

	return nil
}

func (ar *ArtistRep) SelectByID(id uint64) (*models.Artist, error) {
	query := "select id, name, poster, avatar, description from artists where id = $1"

	artist := &models.Artist{}

	err := ar.dbConn.QueryRow(query, id).
		Scan(&artist.ID,
			&artist.Name,
			&artist.Poster,
			&artist.Avatar,
			&artist.Description)

	if err != nil {
		return nil, err
	}

	return artist, nil
}

func (ar *ArtistRep) SelectByParam(count uint64, from uint64) ([]*models.Artist, error) {
	query := "select id, name, poster, avatar, description from artists order by name limit $1 offset $2"

	artists := []*models.Artist{}

	rows, err := ar.dbConn.Query(query, count, from)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		artist := &models.Artist{}
		err := rows.Scan(&artist.ID,
			&artist.Name,
			&artist.Poster,
			&artist.Avatar,
			&artist.Description)

		if err != nil {
			return nil, err
		}

		artists = append(artists, artist)
	}

	return artists, nil
}

func (ar *ArtistRep) SelectByName(name string) (*models.Artist, error) {
	query := "select id, name, poster, avatar, description from artists where name = $1"

	artist := &models.Artist{}

	err := ar.dbConn.QueryRow(query, name).
		Scan(&artist.ID,
			&artist.Name,
			&artist.Poster,
			&artist.Avatar,
			&artist.Description)

	if err != nil {
		return nil, err
	}

	return artist, nil
}
