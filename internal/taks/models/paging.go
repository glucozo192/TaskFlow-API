package models

import (
	"math"

	"github.com/glu-project/internal/taks/constants"
)

// defaultOrderType is the fallback sort direction when none is provided.
// This mirrors pb.OrderType_DESC.String() = "DESC".
const defaultOrderType = "DESC"

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
	if pageSize > constants.MAX_Record {
		p.pageSize = constants.MAX_Record
	}
	if orderType == "" || orderType == "OrderType_NONE" {
		p.orderType = defaultOrderType
	}
	if orderBy == "" {
		p.orderBy = "created_at"
	}
	return p
}

func (p *paging) GetLimit() int32 {
	return p.pageSize
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
