// Package service holds the business logic for each domain: it drives the
// collectors, aggregates their snapshots, filters and sorts, and computes
// derived values (rates, deltas, percentages) that require more than a single
// point-in-time read.
//
// Services keep collectors stateless — a rate metric is computed here from two
// collector snapshots plus a time delta. Per ADR-004 a service may import
// collector, platform, and model, but never command, cli, or render.
package service
