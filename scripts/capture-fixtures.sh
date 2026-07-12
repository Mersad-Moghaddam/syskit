#!/usr/bin/env bash
#
# capture-fixtures.sh — capture a SysKit test fixture set from the current host.
#
# Reads a fixed list of /proc and /sys pseudo-files and copies them into a
# target directory, mirroring their absolute paths, then writes a SOURCE file
# recording provenance (kernel, distro, arch, container status, date).
#
# This script is READ-ONLY with respect to the system: it only reads /proc and
# /sys and writes into the target directory. It never modifies system state.
#
# Usage:
#   capture-fixtures.sh [--force] <target-dir>
#   capture-fixtures.sh --help
#
# Example:
#   capture-fixtures.sh testdata/fixtures/my-8-core-host
#
set -euo pipefail

# The pseudo-files captured for a fixture set. Paths are relative to / and are
# reproduced under the target directory. Missing files are skipped with a
# warning — not every host exposes every interface, and a fixture set that
# records which files were absent is itself useful.
readonly SOURCES=(
	proc/uptime
	proc/loadavg
	proc/stat
	proc/meminfo
	proc/cpuinfo
	proc/self/mountinfo
	proc/self/status
	proc/sys/kernel/hostname
	proc/sys/kernel/osrelease
	proc/sys/kernel/version
	sys/devices/system/cpu/present
	sys/devices/system/cpu/online
)

usage() {
	sed -n '2,20p' "$0" | sed 's/^# \{0,1\}//'
}

main() {
	local force=0
	local target=""

	while [ "$#" -gt 0 ]; do
		case "$1" in
		-h | --help)
			usage
			return 0
			;;
		-f | --force)
			force=1
			shift
			;;
		-*)
			echo "capture-fixtures.sh: unknown option: $1" >&2
			echo "try: capture-fixtures.sh --help" >&2
			return 2
			;;
		*)
			if [ -n "$target" ]; then
				echo "capture-fixtures.sh: unexpected extra argument: $1" >&2
				return 2
			fi
			target="$1"
			shift
			;;
		esac
	done

	if [ -z "$target" ]; then
		echo "capture-fixtures.sh: missing <target-dir>" >&2
		echo "try: capture-fixtures.sh --help" >&2
		return 2
	fi

	# Refuse to write into a non-empty directory unless --force is given, so an
	# existing fixture set is never silently clobbered.
	if [ -d "$target" ] && [ -n "$(ls -A "$target" 2>/dev/null)" ] && [ "$force" -eq 0 ]; then
		echo "capture-fixtures.sh: target '$target' exists and is not empty; pass --force to overwrite" >&2
		return 1
	fi

	mkdir -p "$target"

	local captured=0 missing=0
	for rel in "${SOURCES[@]}"; do
		local src="/$rel"
		local dst="$target/$rel"
		if [ -r "$src" ]; then
			mkdir -p "$(dirname "$dst")"
			# cat (not cp): /proc and /sys pseudo-files report size 0, and a
			# byte-accurate copy must read to EOF rather than trust the size.
			if cat "$src" >"$dst" 2>/dev/null; then
				captured=$((captured + 1))
			else
				echo "capture-fixtures.sh: warning: could not read $src (skipped)" >&2
				rm -f "$dst"
				missing=$((missing + 1))
			fi
		else
			echo "capture-fixtures.sh: warning: $src absent or unreadable (skipped)" >&2
			missing=$((missing + 1))
		fi
	done

	write_source "$target" "$captured" "$missing"

	echo "capture-fixtures.sh: captured $captured file(s), skipped $missing, into $target"
}

# write_source records provenance so a reviewer can tell whether the fixture
# reflects a real kernel and reproduce the capture.
write_source() {
	local target="$1" captured="$2" missing="$3"

	local kernel distro arch container date_utc
	kernel="$(uname -r 2>/dev/null || echo unknown)"
	arch="$(uname -m 2>/dev/null || echo unknown)"
	date_utc="$(date -u '+%Y-%m-%d %H:%M:%S UTC' 2>/dev/null || echo unknown)"

	distro="unknown"
	if [ -r /etc/os-release ]; then
		# shellcheck disable=SC1091
		distro="$(. /etc/os-release 2>/dev/null && echo "${NAME:-unknown} ${VERSION:-}")"
	fi

	container="no"
	if [ -f /.dockerenv ]; then
		container="yes (docker: /.dockerenv present)"
	elif grep -qaE '(docker|containerd|kubepods|lxc)' /proc/1/cgroup 2>/dev/null; then
		container="yes (detected via /proc/1/cgroup)"
	fi

	cat >"$target/SOURCE" <<EOF
kernel: $kernel
distro: $distro
arch: $arch
container: $container
date: $date_utc
captured_files: $captured
skipped_files: $missing

Captured by scripts/capture-fixtures.sh (read-only /proc and /sys copy).
See testdata/README.md for the fixture layout and provenance conventions.
EOF
}

main "$@"
