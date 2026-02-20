package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lgxju/gogretago/internal/domain/entities"
)

func parsePagination(c *gin.Context) entities.PaginationParams {
	params := entities.DefaultPagination()
	if p, err := strconv.Atoi(c.DefaultQuery("page", "1")); err == nil && p > 0 {
		params.Page = p
	}
	if l, err := strconv.Atoi(c.DefaultQuery("limit", "20")); err == nil && l > 0 && l <= 100 {
		params.Limit = l
	}
	return params
}
