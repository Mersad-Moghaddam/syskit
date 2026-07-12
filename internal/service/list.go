package service

import (
	"fmt"
	"sort"
	"strings"
)

// Filter is the parsed `field=value` predicate shared by list services.
type Filter struct{ Field, Value string }

// ParseFilter accepts one equality predicate. Commands may repeat --filter;
// callers apply every predicate as an AND operation.
func ParseFilter(raw string) (Filter, error) {
	field, value, ok := strings.Cut(raw, "=")
	field, value = strings.TrimSpace(field), strings.TrimSpace(value)
	if !ok || field == "" || value == "" {
		return Filter{}, fmt.Errorf("filter %q must use field=value", raw)
	}
	return Filter{Field: field, Value: value}, nil
}

// FilterItems returns the items for which every predicate's accessor value
// equals its requested value. An unknown field is a usage-level error.
func FilterItems[T any](items []T, filters []Filter, fields map[string]func(T) string) ([]T, error) {
	for _, filter := range filters {
		if _, ok := fields[filter.Field]; !ok {
			return nil, fmt.Errorf("unknown filter field %q", filter.Field)
		}
	}
	result := make([]T, 0, len(items))
	for _, item := range items {
		matched := true
		for _, filter := range filters {
			if fields[filter.Field](item) != filter.Value {
				matched = false
				break
			}
		}
		if matched {
			result = append(result, item)
		}
	}
	return result, nil
}

// SortItems sorts a copy of items by a declared field. reverse inverts the
// supplied less function without embedding any command-specific policy here.
func SortItems[T any](items []T, field string, less map[string]func(T, T) bool, reverse bool) ([]T, error) {
	compare, ok := less[field]
	if !ok {
		return nil, fmt.Errorf("unknown sort field %q", field)
	}
	result := append([]T(nil), items...)
	sort.SliceStable(result, func(i, j int) bool {
		if reverse {
			return compare(result[j], result[i])
		}
		return compare(result[i], result[j])
	})
	return result, nil
}

// LimitItems returns at most limit items. Zero means unlimited; negative is invalid.
func LimitItems[T any](items []T, limit int) ([]T, error) {
	if limit < 0 {
		return nil, fmt.Errorf("limit must not be negative")
	}
	if limit == 0 || limit >= len(items) {
		return items, nil
	}
	return items[:limit], nil
}
