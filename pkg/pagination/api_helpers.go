package pagination

import (
	"errors"
	"github.com/labstack/echo/v4"
	"strconv"
)

func GetLimitAndOffset(c echo.Context, defaultLimit int32) (limit int32, offset int32, err error) {
	limitStr := c.QueryParam("limit")
	if limitStr == "" {
		limit = defaultLimit
	} else {
		limit64, err := strconv.Atoi(limitStr)
		if err != nil {
			return 0, 0, errors.New("invalid limit")
		}
		limit = int32(limit64)
	}

	offsetStr := c.QueryParam("offset")
	if offsetStr != "" {
		offset64, err := strconv.Atoi(offsetStr)
		if err != nil {
			return 0, 0, errors.New("invalid offset")
		}
		offset = int32(offset64)
	}

	// Check if limit or offset is negative
	if limit < 0 || offset < 0 {
		return 0, 0, errors.New("limit and offset can't be negative")
	}

	return limit, offset, nil
}
