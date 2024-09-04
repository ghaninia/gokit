package meta

import (
	"github.com/gin-gonic/gin"
	"strings"
)

// --------------
// --------------
// --------------

type SorterEnum string

const (
	SorterEnumAsc  SorterEnum = "ASC"
	SorterEnumDesc SorterEnum = "DESC"
)

type PaginateRequest struct {
	Limit       int
	Page        int
	HasPaginate bool
}

func (p PaginateRequest) GetLimit() int {
	return p.Limit
}

func (p PaginateRequest) GetOffset() int {
	return (p.Page - 1) * p.Limit
}

type DomainWrapper[T any] struct {
	pagination PaginateRequest
	sort       map[string]SorterEnum
	qry        T
}

func (t DomainWrapper[T]) GetPagination() PaginateRequest {
	return t.pagination
}

func (t DomainWrapper[T]) GetSort() map[string]SorterEnum {
	return t.sort
}

func (t DomainWrapper[T]) GetQry() T {
	return t.qry
}

// --------------
// --------------
// --------------

type request[T any] struct {
	ctx *gin.Context
}

func NewRequest[T any](ctx *gin.Context) Abstraction[T] {
	return &request[T]{
		ctx: ctx,
	}
}

type Abstraction[T any] interface {
	Set(qry T) DomainWrapper[T]
}

func (r request[T]) Set(qry T) DomainWrapper[T] {
	return DomainWrapper[T]{
		pagination: PaginateRequest{
			Limit: getParameterStringToInt(r.ctx, LimitField, 10),
			Page:  getParameterStringToInt(r.ctx, PageField, 1),
			HasPaginate: func() bool {
				if value, ok := r.ctx.GetQuery("has_paginate"); ok {
					return strings.ToLower(strings.TrimSpace(value)) == "true"
				}
				return true
			}(),
		},
		sort: func() map[string]SorterEnum {
			sort := make(map[string]SorterEnum)
			for k, v := range r.ctx.QueryMap("sort") {
				sort[k] = SorterEnum(strings.ToUpper(v))
			}
			return sort
		}(),
		qry: qry,
	}
}
