#!/usr/bin/env bash
set -euo pipefail

version="${1:-}"
arch="${2:-amd64}"
output_dir="${3:-dist}"
if [[ -z "$version" ]]; then echo "usage: $0 <version> [amd64|arm64] [output-dir]" >&2; exit 2; fi
case "$version" in v[0-9]*) ;; *) echo "version must start with v" >&2; exit 2;; esac
case "$arch" in amd64|arm64) ;; *) echo "architecture must be amd64 or arm64" >&2; exit 2;; esac
export SOURCE_DATE_EPOCH="${SOURCE_DATE_EPOCH:-$(git log -1 --format=%ct)}"

stage="$(mktemp -d)"
trap 'rm -rf "$stage"' EXIT
install -d "$stage/DEBIAN" "$stage/usr/bin" "$stage/usr/share/doc/syskit" \
  "$stage/usr/share/man/man1"
CGO_ENABLED=0 GOOS=linux GOARCH="$arch" go build -trimpath \
  -ldflags "-s -w -X github.com/Mersad-Moghaddam/syskit/internal/cli.version=$version" \
  -o "$stage/usr/bin/syskit" ./cmd/syskit
chmod 0755 "$stage/usr/bin/syskit"
install -m 0644 LICENSE "$stage/usr/share/doc/syskit/copyright"
install -m 0644 docs/man/syskit.1 "$stage/usr/share/man/man1/syskit.1"
gzip -n -9 "$stage/usr/share/man/man1/syskit.1"
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
