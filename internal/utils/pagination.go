package utils

import (
	"math"

	"gorm.io/gorm"
)

type Pagination struct {
	Total      int `json:"total"`
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalPages int `json:"total_pages"`
}

func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *Pagination) GetLimit() int {
	if p.Limit == 0 {
		p.Limit = 10
	}
	return p.Limit
}

func (p *Pagination) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}

func Paginate(pagination Pagination) func(db *gorm.DB) *gorm.DB {
	// var totalRows int64
	// query.Count(&totalRows)

	// pagination.Total = int(totalRows)
	// totalPages := int(math.Ceil(float64(totalRows) / float64(pagination.Limit)))
	// pagination.TotalPages = totalPages

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit())
	}
}

func GetPaginationInfo(Pagination *Pagination, query *gorm.DB) error {
	var totalRows int64
	result := query.Count(&totalRows)
	if result.Error != nil {
		return result.Error
	}

	Pagination.Total = int(totalRows)
	totalPages := int(math.Ceil(float64(totalRows) / float64(Pagination.Limit)))
	Pagination.TotalPages = totalPages
	return nil
}
