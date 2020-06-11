package utils

import (
	"net/url"
)

func ParseSinceDescLimit(values url.Values) (since, desc, limit string) {
	since = values.Get("since")

	if desc = values.Get("desc"); desc == "" {
		desc = "false"
	}

	if limit = values.Get("limit"); limit == "" {
		limit = "1"
	}

	return since, desc, limit
}
