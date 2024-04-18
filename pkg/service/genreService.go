package service

import "ozinshe/pkg/entity"

type GenreSerivce interface {
	GetAllGenres() ([]entity.Genre, error)
	GetMovieMainsByGenre(int) ([]entity.MovieMain, error)
}

func (s *Service) GetAllGenres() ([]entity.Genre, error) {
	return s.repo.GetAllGenres()
}

func (s *Service) GetMovieMainsByGenre(genreId int) ([]entity.MovieMain, error) {
	return s.repo.GetMovieMainsByGenre(genreId)
}
