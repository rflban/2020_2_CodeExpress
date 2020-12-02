package grpc_admin

import (
	"context"

	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/admin/proto_admin"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/artist"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
)

type AdminGRPCUsecase struct {
	artistRep artist.ArtistRep
}

func NewAdminGRPCUsecase(artistRep artist.ArtistRep) *AdminGRPCUsecase {
	return &AdminGRPCUsecase{
		artistRep: artistRep,
	}
}

func ArtistToGRPCArtist(artist *models.Artist) *proto_admin.Artist {
	return &proto_admin.Artist{
		ID:          artist.ID,
		Name:        artist.Name,
		Poster:      artist.Poster,
		Avatar:      artist.Avatar,
		Description: artist.Description,
	}
}

func GRPCArtistToArtist(artist *proto_admin.Artist) *models.Artist {
	return &models.Artist{
		ID:          artist.ID,
		Name:        artist.Name,
		Poster:      artist.Poster,
		Avatar:      artist.Avatar,
		Description: artist.Description,
	}
}

func (aGu *AdminGRPCUsecase) CreateArtist(ctx context.Context, newArtist *proto_admin.Artist) (*proto_admin.Artist, error) {
	artist := &models.Artist{
		Name:        newArtist.Name,
		Description: newArtist.Description,
	}

	err := aGu.artistRep.Insert(artist)

	if err != nil {
		return &proto_admin.Artist{}, err
	}

	return &proto_admin.Artist{
		ID:          artist.ID,
		Name:        artist.Name,
		Poster:      artist.Poster,
		Avatar:      artist.Avatar,
		Description: artist.Description,
	}, nil
}
