package gorm

import (
	"github.com/Cardsity/management-api/db/repositories"
	"gorm.io/gorm"
)

// The base every gorm repository contains.
type BaseRepository struct {
	Db *gorm.DB
}

// Performs a simple operation that queries according to the where parameter and writes to first.
func (br *BaseRepository) get(where interface{}, first interface{}) repositories.RepositoryError {
	result := br.Db.Where(where).First(first)
	if result.Error != nil {
		return repositories.NewRepositoryError(result.Error)
	}
	return repositories.RepositoryError{}
}

// Performs a simple operation that deletes according to the supplied value.
func (br *BaseRepository) delete(value interface{}) repositories.RepositoryError {
	result := br.Db.Delete(value)
	if result.Error != nil {
		return repositories.NewRepositoryError(result.Error)
	}
	return repositories.RepositoryError{}
}

// Performs a simple create operation for the supplied value.
func (br *BaseRepository) create(value interface{}) repositories.RepositoryError {
	result := br.Db.Create(value)
	if result.Error != nil {
		return repositories.NewRepositoryError(result.Error)
	}
	return repositories.RepositoryError{}
}
