package utils

import (
	"net/http"
	"strconv"
	"strings"
)

func ParseQueryParams(r *http.Request) map[string]string {
	query := r.URL.Query()
	params := make(map[string]string, len(query))
	for k, v := range query {
		if len(v) > 0 {
			params[k] = v[0]
		}
	}
	return params
}
func ParseIntWithDefaultAndMax(s string, defaultVal, maxVal int) int {
	if s == "" {
		return defaultVal
	}

	val, err := strconv.Atoi(s)
	if err != nil || val < 1 {
		return defaultVal
	}

	if maxVal > 0 && val > maxVal {
		return maxVal
	}

	return val
}
func ParseCSV(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return []string{}
	}

	parts := strings.Split(raw, ",")
	var result []string
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
