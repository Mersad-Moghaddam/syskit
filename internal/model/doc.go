// Package model holds the shared, typed domain structs that flow up from the
// collectors through the services to the renderers.
//
// Structs carry json (and later yaml) tags with snake_case field names and
// explicit units. Raw kernel counters are preserved alongside — never
// overwritten by — any derived values, so a single snapshot can feed both the
// raw JSON output and service-computed rates. This package has no dependencies
// on any other SysKit layer.
package model
