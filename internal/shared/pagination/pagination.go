package pagination

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Pagination struct {
	Limit  int    `json:"limit" validate:"gte=1,lte=20"`
	Offset int    `json:"offset" validate:"gte=0"`
	Sort   string `json:"sort" validate:"oneof=asc desc"`

	Tags   []string `json:"tags" validate:"max=5"`
	Search string   `json:"search" validate:"max=10"`

	From string `json:"from"`
	To   string `json:"to"`
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

	tags := qs.Get("tags")
	if tags != "" {
		pagination.Tags = strings.Split(tags, ",")
	}

	search := qs.Get("search")
	if search != "" {
		pagination.Search = search
	}

	from := qs.Get("from")
	if from != "" {
		pagination.From = parseTime(from)
	}

	to := qs.Get("to")
	if from != "" {
		pagination.To = parseTime(to)
	}

	return pagination, nil
}

func parseTime(value string) string {
	t, err := time.Parse(time.DateTime, value)
	if err != nil {
		return ""
	}

	return t.Format(time.DateTime)
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
