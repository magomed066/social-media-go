package pagination

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

type Pagination struct {
	Limit  int    `json:"limit" validate:"gte=1,lte=20"`
	Offset int    `json:"offset" validate:"gte=0"`
	Sort   string `json:"sort" validate:"oneof=asc desc"`
}

func (pagination Pagination) Parse(r *http.Request) (Pagination, error) {
	qs := r.URL.Query()

	limit, err := checkAndParseValue(qs, "limit")
	if err != nil {
		return pagination, err
	}

	pagination.Limit = limit

	offset, err := checkAndParseValue(qs, "offset")
	if err != nil {
		return pagination, err
	}

	pagination.Offset = offset

	sort := qs.Get("sort")
	if sort != "" {
		pagination.Sort = sort
	}

	return pagination, nil
}

func checkAndParseValue(qs url.Values, key string) (int, error) {
	value := qs.Get(key)
	if value != "" {
		parsedValue, err := strconv.Atoi(value)
		if err != nil {
			return 0, err
		}

		return parsedValue, nil

	}

	return 0, errors.New("Failed to get value from url")
}
