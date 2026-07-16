#!/bin/sh
# install.sh installs the latest SysKit release for the current Linux CPU.
# Optional overrides: SYSKIT_VERSION=v1.0.0 SYSKIT_INSTALL_PREFIX=$HOME/.local
set -eu

repository="Mersad-Moghaddam/syskit"
prefix="${SYSKIT_INSTALL_PREFIX:-/usr/local}"
curl_family="${SYSKIT_CURL_FAMILY:-}"

info() {
	printf '%s\n' "==> $*" >&2
}

success() {
	printf '%s\n' " ✓  $*" >&2
}

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

case "$curl_family" in
	"") ;;
	4 | 6) ;;
	*)
		echo "syskit installer: SYSKIT_CURL_FAMILY must be 4 or 6" >&2
		exit 1
		;;
esac

fetch() {
	case "$curl_family" in
		4) curl -4 "$@" ;;
		6) curl -6 "$@" ;;
		*) curl "$@" ;;
	esac
}

download() {
	description="$1"
	destination="$2"
	url="$3"
	info "$description"
	if [ -t 2 ]; then
		fetch -fL --connect-timeout 15 --progress-bar -o "$destination" "$url"
	else
		fetch -fL --connect-timeout 15 -o "$destination" "$url"
	fi
	success "$description"
}

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
success "Linux $(uname -m) detected; selecting $arch release"

version="${SYSKIT_VERSION:-}"
if [ -z "$version" ]; then
	# Use github.com rather than api.github.com: some server firewalls allow
	# release downloads but block the GitHub API. GitHub redirects this URL to
	# the latest release tag, which curl exposes as url_effective.
	info "Resolving the latest SysKit release"
	release_url="$(fetch -fsSIL --connect-timeout 15 --max-time 30 -o /dev/null -w '%{url_effective}' "https://github.com/$repository/releases/latest")"
	version="${release_url##*/}"
fi

case "$version" in
	v[0-9]*) ;;
	*)
		echo "syskit installer: could not determine a valid release version" >&2
		exit 1
		;;
esac
success "Using SysKit $version"

archive="syskit_${version#v}_linux_${arch}.tar.gz"
base_url="https://github.com/$repository/releases/download/$version"
temp_dir="$(mktemp -d "${TMPDIR:-/tmp}/syskit-install.XXXXXX")"
trap 'rm -rf "$temp_dir"' EXIT HUP INT TERM
info "Preparing temporary workspace: $temp_dir"

download "Downloading release checksums" "$temp_dir/SHA256SUMS" "$base_url/SHA256SUMS"
download "Downloading $archive" "$temp_dir/$archive" "$base_url/$archive"
info "Verifying SHA-256 checksum"
(
	cd "$temp_dir"
	sha256sum -c SHA256SUMS --ignore-missing
)
success "Release archive checksum verified"
info "Extracting release archive"
tar -xzf "$temp_dir/$archive" -C "$temp_dir"
success "Release archive extracted"

if [ "${SYSKIT_INSTALL_PREFIX+x}" = x ]; then
	info "Creating user-selected install directories under $prefix"
	mkdir -p "$prefix/bin" "$prefix/share/man/man1"
	privilege=""
elif [ "$(id -u)" -eq 0 ]; then
	info "Installing as root under $prefix"
	privilege=""
elif command -v sudo >/dev/null 2>&1; then
	info "Installing under $prefix with sudo"
	privilege="sudo"
else
	echo "syskit installer: cannot write to $prefix; rerun as root or set SYSKIT_INSTALL_PREFIX" >&2
	exit 1
fi

info "Installing syskit binary"
# shellcheck disable=SC2086 # privilege is intentionally an optional command.
$privilege install -Dm 0755 "$temp_dir/syskit_${version#v}_linux_${arch}" "$prefix/bin/syskit"
info "Installing syskit(1) manual page"
# shellcheck disable=SC2086 # privilege is intentionally an optional command.
$privilege install -Dm 0644 "$temp_dir/syskit.1" "$prefix/share/man/man1/syskit.1"

success "SysKit $version installed at $prefix/bin/syskit"
printf '%s\n' 'Run: syskit version' >&2
