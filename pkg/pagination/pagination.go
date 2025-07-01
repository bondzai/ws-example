package pagination

import (
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
)

// Pagination holds the basic parameters for pagination.
type Pagination struct {
	Page     int
	PageSize int
}

// Validate ensures that the pagination parameters have valid values.
func (p *Pagination) Validate() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 10
	}
}

// GormScope returns a GORM scope function to apply pagination to a query.
// Usage with GORM (e.g., PostgreSQL):
//
//	db.Scopes(GormScope(pagination)).Find(&results)
func GormScope(p Pagination) func(db *gorm.DB) *gorm.DB {
	p.Validate()
	offset := (p.Page - 1) * p.PageSize

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(offset).Limit(p.PageSize)
	}
}

// MongoOptions returns MongoDB find options with pagination (skip and limit) applied.
// Usage with MongoDB:
//
//	collection.Find(ctx, filter, MongoOptions(pagination))
func MongoOptions(p Pagination) *options.FindOptions {
	p.Validate()
	skip := int64((p.Page - 1) * p.PageSize)
	limit := int64(p.PageSize)

	return options.Find().SetSkip(skip).SetLimit(limit)
}

// Slice paginates a plain slice.
// This generic function works for slices of any type.
// Usage with plain arrays/slices:
//
//	paginatedData := Slice(yourSlice, pagination)
func Slice[T any](data []T, p Pagination) []T {
	p.Validate()
	start := (p.Page - 1) * p.PageSize
	if start >= len(data) {
		return []T{}
	}

	end := start + p.PageSize
	if end > len(data) {
		end = len(data)
	}

	return data[start:end]
}
