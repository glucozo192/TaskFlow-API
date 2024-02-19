package models

import (
	"math"

	"github.com/glu-project/idl/pb"
	"github.com/glu-project/internal/user/constants"
)

// type Paging struct {
// 	OrderType string
// 	Offset    int32
// 	Limit     int32
// 	OrderBy   string
// 	Search    string
// }

// func (p *Paging) Default() {
// 	if p.Limit == 0 {
// 		p.Limit = 10
// 	}
// 	if p.Limit > constants.MAX_Retried_Record {
// 		p.Limit = 100
// 	}
// 	if p.OrderType == pb.OrderType_OrderType_NONE.String() {
// 		p.OrderType = pb.OrderType_DESC.String()
// 	}
// 	if p.OrderBy == "" {
// 		p.OrderBy = "created_at"
// 	}
// 	p.OrderBy = string_utils.ToSnakeCase(p.OrderBy)
// }

// func (p *Paging) DefaultV2() {
// 	if p.OrderType == pb.OrderType_OrderType_NONE.String() {
// 		p.OrderType = pb.OrderType_DESC.String()
// 	}
// 	if p.OrderBy == "" {
// 		p.OrderBy = "created_at"
// 	}
// 	p.OrderBy = string_utils.ToSnakeCase(p.OrderBy)
// }

type Paging interface {
	GetLimit() int32
	GetOffset() int32
	GetOrderBy() string
	GetOrderType() string
	GetPage() int32
	GetPageSize() int32
	GetQuery() string
	CalTotalPages(total int32) int32
}

type paging struct {
	page      int32
	pageSize  int32
	orderBy   string
	orderType string
	q         string
}

func NewPagingWithDefault(page, pageSize int32, orderBy, orderType, search string) Paging {
	p := &paging{
		page:      page,
		pageSize:  pageSize,
		orderBy:   orderBy,
		orderType: orderType,
		q:         search,
	}

	if page == 0 {
		p.page = 1
	}
	if pageSize == 0 {
		p.pageSize = 10
	}
	if pageSize > constants.MAX_Retried_Record {
		p.pageSize = constants.MAX_Retried_Record
	}
	if orderType == pb.OrderType_OrderType_NONE.String() {
		p.orderType = pb.OrderType_DESC.String()
	}
	if orderBy == "" {
		p.orderBy = "created_at"
	}
	return p
}

func (p *paging) GetLimit() int32 {
	return p.pageSize * p.page
}
func (p *paging) GetOffset() int32 {
	return (p.page - 1) * p.pageSize
}

func (p *paging) GetOrderBy() string {
	return p.orderBy
}
func (p *paging) GetOrderType() string {
	return p.orderType
}
func (p *paging) GetPage() int32 {
	return p.page
}
func (p *paging) GetPageSize() int32 {
	return p.pageSize
}

func (p *paging) GetQuery() string {
	return p.q
}

func (p *paging) CalTotalPages(total int32) int32 {
	return int32(math.Ceil(float64(total) / float64(p.pageSize)))
}
