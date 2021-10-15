package dto

import (
	"github.com/c0llinn/ebook-store/internal/catalog/model"
	"time"
)

type SearchBooks struct {
	Title       string `form:"title"`
	Description string `form:"description"`
	AuthorName  string `form:"authorName"`
	Page        int    `form:"page"`
	PerPage     int    `form:"perPage"`
}

func (s *SearchBooks) ToDomain() model.BookQuery {
	if s.Page == 0 {
		s.Page = 1
	}

	if s.PerPage == 0 {
		s.PerPage = 10
	}

	return model.BookQuery{
		Title:       s.Title,
		Description: s.Description,
		AuthorName:  s.AuthorName,
		Limit:       s.PerPage,
		Offset:      (s.Page - 1) * s.PerPage,
	}
}

type CreateBook struct {
	Title       string    `form:"title" binding:"required,max=100"`
	Description string    `form:"description" binding:"required"`
	AuthorName  string    `form:"authorName" binding:"required,max=100"`
	Price       int       `form:"price" binding:"required,gt=0"`
	ReleaseDate time.Time `form:"releaseDate" binding:"required"`
}

func (c CreateBook) ToDomain(id string) model.Book {
	return model.Book{
		ID:          id,
		Title:       c.Title,
		Description: c.Description,
		AuthorName:  c.AuthorName,
		Price:       c.Price,
		ReleaseDate: c.ReleaseDate,
	}
}

type UpdateBook struct {
	Title       *string `form:"title" binding:"omitempty,max=100"`
	Description *string `form:"description" binding:"omitempty"`
	AuthorName  *string `form:"authorName" binding:"omitempty,max=100"`
}

func (u UpdateBook) ToDomain(existing model.Book) model.Book {
	updated := existing

	if u.Title != nil && *u.Title != "" {
		updated.Title = *u.Title
	}

	if u.Description != nil && *u.Description != "" {
		updated.Description = *u.Description
	}

	if u.AuthorName != nil && *u.AuthorName != "" {
		updated.AuthorName = *u.AuthorName
	}

	return updated
}
