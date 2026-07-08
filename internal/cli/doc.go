// Package cli is the top layer: it wires the Cobra root command and its
// subcommands (internal/cli/command), loads configuration (defaults←file←env←
// flags), selects the renderer, configures the slog logger, presents errors,
// and maps propagated sentinel errors to process exit codes.
//
// It is the only layer that assigns exit codes and the only one that may run
// the TUI. Per ADR-004 cli sits above every other SysKit layer and may import
// them; nothing imports cli. Business logic and kernel I/O never live here.
package cli
