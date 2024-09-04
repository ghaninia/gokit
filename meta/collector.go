package meta

import (
	"github.com/gin-gonic/gin"
	"reflect"
)

type collection[T any] struct {
	DomainData []T        `json:"domainData"`
	Paginator  Pagination `json:"paginator"`
}

type Collect[T any] interface {
	GetDomainData() []T
	ToDomain(data []T) Collect[T]
	GetMeta(ctx *gin.Context) Meta
}

// NewCollection creates a new collection.
func NewCollection[T any](data any) Collect[T] {
	return collection[T]{
		DomainData: make([]T, 0),
		Paginator: Pagination{
			TotalCount: func() int {
				if data != nil {
					switch reflect.TypeOf(data).Kind() {
					case reflect.Slice:
						slice := reflect.ValueOf(data)
						if slice.Len() > 0 {
							return slice.Index(0).
								FieldByName("TotalCount").
								Interface().(int)
						}
					}
				}
				return 0
			}(),
		},
	}
}

// GetDomainData returns the domain data.
func (c collection[T]) GetDomainData() []T {
	return c.DomainData
}

// ToDomain sets the domain data.
func (c collection[T]) ToDomain(data []T) Collect[T] {
	c.DomainData = data
	return c
}

func (c collection[T]) GetMeta(ctx *gin.Context) Meta {
	return Meta{
		Pagination: Pagination{
			TotalCount: c.Paginator.TotalCount,
			PerPage:    getParameterStringToInt(ctx, LimitField, 10),
			Page:       getParameterStringToInt(ctx, PageField, 1),
			PageCount: func() int {
				if c.Paginator.TotalCount > 0 {
					return c.Paginator.TotalCount / getParameterStringToInt(ctx, LimitField, 10)
				}
				return 0
			}(),
		},
	}
}

type Meta struct {
	Pagination Pagination `json:"pagination"`
}

type Pagination struct {
	Page       int `json:"page"`
	PerPage    int `json:"perPage"`
	PageCount  int `json:"pageCount"`
	TotalCount int `json:"totalCount"`
}
