#!/usr/bin/env bash
set -euo pipefail

version="${1:-}"
if [[ -z "$version" ]]; then echo "usage: $0 <version> [output-dir]" >&2; exit 2; fi
output_dir="${2:-dist}"
case "$version" in v[0-9]*) ;; *) echo "version must start with v" >&2; exit 2;; esac
export SOURCE_DATE_EPOCH="${SOURCE_DATE_EPOCH:-$(git log -1 --format=%ct)}"
rm -rf "$output_dir"
mkdir -p "$output_dir"
for arch in amd64 arm64; do
  name="syskit_${version#v}_linux_${arch}"
  CGO_ENABLED=0 GOOS=linux GOARCH="$arch" go build -trimpath -ldflags "-s -w -X github.com/Mersad-Moghaddam/syskit/internal/cli.version=$version" -o "$output_dir/$name" ./cmd/syskit
  chmod 0755 "$output_dir/$name"
  install -m 0644 LICENSE "$output_dir/LICENSE"
  install -m 0644 docs/man/syskit.1 "$output_dir/syskit.1"
  tar -C "$output_dir" --sort=name --mtime="@$SOURCE_DATE_EPOCH" \
    --owner=0 --group=0 --numeric-owner -cf - LICENSE syskit.1 "$name" |
    gzip -n > "$output_dir/$name.tar.gz"
  chmod 0644 "$output_dir/$name.tar.gz"
  rm "$output_dir/$name"
  rm "$output_dir/LICENSE"
  rm "$output_dir/syskit.1"
done
scripts/write-checksums.sh "$output_dir"
