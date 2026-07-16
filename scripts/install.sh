#!/bin/sh
# install.sh installs the latest SysKit release for the current Linux CPU.
# Optional overrides: SYSKIT_VERSION=v1.0.0 SYSKIT_INSTALL_PREFIX=$HOME/.local
set -eu

repository="Mersad-Moghaddam/syskit"
prefix="${SYSKIT_INSTALL_PREFIX:-/usr/local}"

require() {
	if ! command -v "$1" >/dev/null 2>&1; then
		echo "syskit installer: $1 is required" >&2
		exit 1
	fi
}

require curl
require tar
require sha256sum
require install
require mktemp

if [ "$(uname -s)" != "Linux" ]; then
	echo "syskit installer: SysKit is supported only on Linux" >&2
	exit 1
fi

case "$(uname -m)" in
	x86_64) arch=amd64 ;;
	aarch64 | arm64) arch=arm64 ;;
	*)
		echo "syskit installer: unsupported architecture: $(uname -m)" >&2
		exit 1
		;;
esac

version="${SYSKIT_VERSION:-}"
if [ -z "$version" ]; then
	metadata="$(curl -fsSL "https://api.github.com/repos/$repository/releases/latest")"
	version="$(printf '%s\n' "$metadata" | sed -n 's/^[[:space:]]*"tag_name"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' | head -n 1)"
fi

case "$version" in
	v[0-9]*) ;;
	*)
		echo "syskit installer: could not determine a valid release version" >&2
		exit 1
		;;
esac

archive="syskit_${version#v}_linux_${arch}.tar.gz"
base_url="https://github.com/$repository/releases/download/$version"
temp_dir="$(mktemp -d "${TMPDIR:-/tmp}/syskit-install.XXXXXX")"
trap 'rm -rf "$temp_dir"' EXIT HUP INT TERM

curl -fsSL -o "$temp_dir/SHA256SUMS" "$base_url/SHA256SUMS"
curl -fsSL -o "$temp_dir/$archive" "$base_url/$archive"
(
	cd "$temp_dir"
	sha256sum -c SHA256SUMS --ignore-missing
)
tar -xzf "$temp_dir/$archive" -C "$temp_dir"

if [ "${SYSKIT_INSTALL_PREFIX+x}" = x ]; then
	mkdir -p "$prefix/bin" "$prefix/share/man/man1"
	privilege=""
elif [ "$(id -u)" -eq 0 ]; then
	privilege=""
elif command -v sudo >/dev/null 2>&1; then
	privilege="sudo"
else
	echo "syskit installer: cannot write to $prefix; rerun as root or set SYSKIT_INSTALL_PREFIX" >&2
	exit 1
fi

# shellcheck disable=SC2086 # privilege is intentionally an optional command.
$privilege install -Dm 0755 "$temp_dir/syskit_${version#v}_linux_${arch}" "$prefix/bin/syskit"
# shellcheck disable=SC2086 # privilege is intentionally an optional command.
$privilege install -Dm 0644 "$temp_dir/syskit.1" "$prefix/share/man/man1/syskit.1"

echo "SysKit $version installed at $prefix/bin/syskit"
echo "Run: syskit version"
