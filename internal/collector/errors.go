package collector

import "errors"

// Domain sentinel errors that every collector reuses so the service and CLI
// layers can classify failures with errors.Is instead of string matching
// (specs/error-handling.md "Sentinel Errors"; specs/collectors.md "Error
// Classification"). Collectors wrap these with %w, adding operation context at
// each step (e.g. fmt.Errorf("parsing proc/loadavg: %w", ErrParse)).
//
// Error-classification model (specs/collectors.md):
//
//   - Malformed kernel data          -> ErrParse
//   - Required field absent          -> ErrFieldMissing
//   - Missing OPTIONAL data          -> NOT an error. The field is represented
//     in the model as unavailable (a zero value, empty slice, or an explicit
//     "present" flag) and Collect still returns a snapshot with nil error. A
//     collector must never fail a whole snapshot because one optional interface
//     is absent — that would defeat partial-data reporting.
//   - Permission denied / not found / unsupported kernel capability -> these
//     originate in the platform layer as platform.ErrPermission,
//     platform.ErrNotFound, and platform.ErrUnsupported. Collectors pass them
//     through unchanged (wrapped with %w for context), so a caller's
//     errors.Is(err, platform.ErrPermission) still matches. Collectors do not
//     redefine these conditions with their own sentinels.
//
// The service layer, not the collector, decides whether a returned error is
// fatal or can be surfaced as partial data.
var (
	// ErrParse indicates the kernel interface was readable but its contents
	// could not be parsed into the expected shape (e.g. a numeric field that is
	// not a number, or a structurally corrupt line).
	ErrParse = errors.New("malformed kernel data")

	// ErrFieldMissing indicates a REQUIRED field was absent from otherwise
	// well-formed kernel data. Optional missing data is represented as
	// unavailable in the model and does not use this sentinel.
	ErrFieldMissing = errors.New("required field missing")
)
