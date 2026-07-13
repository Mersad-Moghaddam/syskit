#!/usr/bin/env bash
set -euo pipefail

version="${1:-}"
arch="${2:-amd64}"
output_dir="${3:-dist}"
if [[ -z "$version" ]]; then echo "usage: $0 <version> [amd64|arm64] [output-dir]" >&2; exit 2; fi
case "$version" in v[0-9]*) ;; *) echo "version must start with v" >&2; exit 2;; esac
case "$arch" in amd64|arm64) ;; *) echo "architecture must be amd64 or arm64" >&2; exit 2;; esac

stage="$(mktemp -d)"
trap 'rm -rf "$stage"' EXIT
install -d "$stage/DEBIAN" "$stage/usr/bin" "$stage/usr/share/doc/syskit"
CGO_ENABLED=0 GOOS=linux GOARCH="$arch" go build -trimpath \
  -ldflags "-s -w -X github.com/Mersad-Moghaddam/syskit/internal/cli.version=$version" \
  -o "$stage/usr/bin/syskit" ./cmd/syskit
install -m 0644 LICENSE "$stage/usr/share/doc/syskit/copyright"
cat > "$stage/DEBIAN/control" <<EOF
Package: syskit
Version: ${version#v}
Architecture: $arch
Maintainer: SysKit contributors
Section: admin
Priority: optional
Homepage: https://github.com/Mersad-Moghaddam/syskit
Description: Linux-native system inspection and diagnostics toolkit
 SysKit reads procfs, sysfs, Netlink, and cgroups without shelling out.
EOF
mkdir -p "$output_dir"
dpkg-deb --root-owner-group --build "$stage" "$output_dir/syskit_${version#v}_${arch}.deb"
