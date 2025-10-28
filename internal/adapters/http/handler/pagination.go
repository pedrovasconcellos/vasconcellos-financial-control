package handler

import (
	"errors"
	"strconv"
)

const (
	defaultLimit = 100
	maxLimit     = 200
)

func parsePagination(limitRaw string, offsetRaw string) (int64, int64, error) {
	limit := int64(defaultLimit)
	offset := int64(0)

	if limitRaw != "" {
		value, err := strconv.ParseInt(limitRaw, 10, 64)
		if err != nil || value <= 0 {
			return 0, 0, errors.New("invalid limit parameter")
		}
		if value > maxLimit {
			value = maxLimit
		}
		limit = value
	}

	if offsetRaw != "" {
		value, err := strconv.ParseInt(offsetRaw, 10, 64)
		if err != nil || value < 0 {
			return 0, 0, errors.New("invalid offset parameter")
		}
		offset = value
	}

	return limit, offset, nil
}
