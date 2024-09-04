package meta

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WithPaginator struct {
	TotalCount int `gorm:"->"`
}

type Paginator interface {
	GetPaginator() WithPaginator
}

type ValidSortColumns []string

// Paginate is a middleware that paginates the query result based on the given configuration.
func Paginate(config PaginateRequest) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if !config.HasPaginate {
			return db
		}

		return db.
			Clauses(clause.Select{
				Columns: []clause.Column{
					{
						Name:  "count(*) OVER()",
						Alias: "total_count",
						Raw:   true,
					},
					{
						Name:  "*",
						Table: clause.CurrentTable,
						Raw:   true,
					},
				},
			}).
			Offset(config.GetOffset()).
			Limit(config.GetLimit())
	}
}

// Sort is a middleware that sorts the query result based on the given configuration.
func Sort(config map[string]SorterEnum, valid ValidSortColumns) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for k, v := range config {
			if valid == nil {
				db = db.Order(k + " " + string(v))
				continue
			}
			for _, a := range valid {
				if k == a {
					db = db.Order(k + " " + string(v))
				}
			}
		}
		return db
	}
}
