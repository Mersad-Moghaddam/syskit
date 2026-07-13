#!/usr/bin/env bash
set -euo pipefail

version="${1:-}"
release_dir="${2:-dist}"
output_dir="${3:-$release_dir}"
if [[ -z "$version" ]]; then echo "usage: $0 <version> [release-dir] [output-dir]" >&2; exit 2; fi
case "$version" in v[0-9]*) ;; *) echo "version must start with v" >&2; exit 2;; esac

pkgver="${version#v}"
amd64_archive="syskit_${pkgver}_linux_amd64.tar.gz"
arm64_archive="syskit_${pkgver}_linux_arm64.tar.gz"
checksum_file="$release_dir/SHA256SUMS"
if [[ ! -f "$checksum_file" ]]; then
  echo "$checksum_file is required; run scripts/build-release.sh first" >&2
  exit 1
fi
amd64_sum="$(awk -v name="$amd64_archive" '$2 == name {print $1}' "$checksum_file")"
arm64_sum="$(awk -v name="$arm64_archive" '$2 == name {print $1}' "$checksum_file")"
if [[ -z "$amd64_sum" || -z "$arm64_sum" ]]; then
  echo "SHA256SUMS must contain both release archives" >&2
  exit 1
fi

stage="$(mktemp -d)"
trap 'rm -rf "$stage"' EXIT
cat > "$stage/PKGBUILD" <<EOF
pkgname=syskit-bin
pkgver=$pkgver
pkgrel=1
pkgdesc='Linux-native system inspection and diagnostics toolkit'
arch=('x86_64' 'aarch64')
url='https://github.com/Mersad-Moghaddam/syskit'
license=('MIT')
provides=('syskit')
conflicts=('syskit')
options=('!strip')
source_x86_64=("syskit-$pkgver-x86_64.tar.gz::https://github.com/Mersad-Moghaddam/syskit/releases/download/v$pkgver/$amd64_archive")
source_aarch64=("syskit-$pkgver-aarch64.tar.gz::https://github.com/Mersad-Moghaddam/syskit/releases/download/v$pkgver/$arm64_archive")
sha256sums_x86_64=('$amd64_sum')
sha256sums_aarch64=('$arm64_sum')

package() {
  local goarch=amd64
  [[ "\$CARCH" == aarch64 ]] && goarch=arm64
  install -Dm755 "\$srcdir/syskit_${pkgver}_linux_\$goarch" "\$pkgdir/usr/bin/syskit"
  install -Dm644 "\$srcdir/LICENSE" "\$pkgdir/usr/share/licenses/syskit/LICENSE"
  install -Dm644 "\$srcdir/syskit.1" "\$pkgdir/usr/share/man/man1/syskit.1"
}
EOF
cat > "$stage/.SRCINFO" <<EOF
pkgbase = syskit-bin
	pkgdesc = Linux-native system inspection and diagnostics toolkit
	pkgver = $pkgver
	pkgrel = 1
	url = https://github.com/Mersad-Moghaddam/syskit
	arch = x86_64
	arch = aarch64
	license = MIT
	conflicts = syskit
	provides = syskit
	options = !strip
	source_x86_64 = syskit-$pkgver-x86_64.tar.gz::https://github.com/Mersad-Moghaddam/syskit/releases/download/v$pkgver/$amd64_archive
	sha256sums_x86_64 = $amd64_sum
	source_aarch64 = syskit-$pkgver-aarch64.tar.gz::https://github.com/Mersad-Moghaddam/syskit/releases/download/v$pkgver/$arm64_archive
	sha256sums_aarch64 = $arm64_sum

pkgname = syskit-bin
EOF

mkdir -p "$output_dir"
tar -C "$stage" --sort=name --mtime="@${SOURCE_DATE_EPOCH:-$(git log -1 --format=%ct)}" \
  --owner=0 --group=0 --numeric-owner -czf \
  "$output_dir/syskit_${pkgver}_aur.tar.gz" .SRCINFO PKGBUILD
chmod 0644 "$output_dir/syskit_${pkgver}_aur.tar.gz"
