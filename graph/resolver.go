//go:generate go run github.com/99designs/gqlgen
package graph

import "gorm.io/gorm"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	ChatSocket map[int]chan string
}

func NewResolver() *Resolver {
	return &Resolver{ChatSocket: map[int]chan string{}}
}

func Paginate10(r int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page := r

		pageSize := 10

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func Paginate6(r int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page := r

		pageSize := 6

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
