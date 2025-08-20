package filters

import (
	"gorm.io/gorm"
	"strings"
)

type Pagination struct {
	Offset  *int    `json:"offset,omitempty"`
	Limit   *int    `json:"limit,omitempty"`
	OrderBy *string `json:"order_by,omitempty"`
	SortDir *string `json:"sort_dir,omitempty"` // "asc", "desc"
}

const (
	minLimit = 50
	maxLimit = 1000
)

func (p *Pagination) GetLimit() int {
	if p.Limit == nil || *p.Limit <= 0 {
		return minLimit
	}
	if *p.Limit > maxLimit {
		return maxLimit
	}
	return *p.Limit
}

func (p *Pagination) GetOffset() int {
	if p.Offset == nil || *p.Offset < 0 {
		return 0
	}
	return *p.Offset
}

func PaginateQuery(query *gorm.DB, pagination *Pagination, allowedOrderFields []string) *gorm.DB {
	if pagination == nil {
		return query
	}

	limit := pagination.GetLimit()
	offset := pagination.GetOffset()
	query.Limit(limit).Offset(offset)

	if pagination.OrderBy != nil {
		orderField := strings.TrimSpace(*pagination.OrderBy)
		if isValidOrderField(orderField, allowedOrderFields) {
			orderClause := orderField

			if pagination.SortDir != nil {
				sortDir := strings.ToUpper(strings.TrimSpace(*pagination.SortDir))
				if sortDir == "ASC" || sortDir == "DESC" {
					orderClause += " " + sortDir
				}
			}

			query = query.Order(orderClause)
		}
	}

	return query
}

func isValidOrderField(field string, allowedFields []string) bool {
	if field == "" {
		return false
	}

	for _, allowed := range allowedFields {
		if strings.EqualFold(field, allowed) {
			return true
		}
	}
	return false
}
