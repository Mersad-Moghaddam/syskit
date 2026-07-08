// Command syskit is the entry point for the SysKit CLI — a Linux-only,
// read-only system-inspection toolkit that reads native kernel interfaces
// (/proc, /sys, Netlink, cgroups) directly and renders the result as table,
// JSON, YAML, or an interactive terminal dashboard.
//
// main only wires and exits: it builds and runs the Cobra command tree in
// internal/cli and exits with the code that layer returns. All CLI logic lives
// in internal/cli.
package main

import (
	"os"

	"github.com/Mersad-Moghaddam/syskit/internal/cli"
)

func main() {
	os.Exit(cli.Main())
}
