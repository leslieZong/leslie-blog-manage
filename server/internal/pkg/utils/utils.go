package utils

import (
	appErrors "leslie-blog-server/internal/errors"
	"net/http"
	"strconv"
)

func ParseStatus(statusValue string) (*int8, error) {
	if statusValue != "" {
		value, err := strconv.ParseInt(
			statusValue,
			10,
			8,
		)

		if err != nil {
			return nil, appErrors.New(
				http.StatusBadRequest,
				appErrors.ErrInvalidParams,
				"invalid status",
			)

		}

		statusInt := int8(value)

		if statusInt != 0 && statusInt != 1 {
			return nil, appErrors.New(
				http.StatusBadRequest,
				appErrors.ErrInvalidParams,
				"invalid status",
			)
		}
		return &statusInt, nil
	}
	return nil, nil
}
