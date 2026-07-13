#!/usr/bin/env bash
set -euo pipefail

version="${1:-}"
if [[ -z "$version" ]]; then echo "usage: $0 <version> [output-dir]" >&2; exit 2; fi
output_dir="${2:-dist}"
case "$version" in v[0-9]*) ;; *) echo "version must start with v" >&2; exit 2;; esac
rm -rf "$output_dir"
mkdir -p "$output_dir"
for arch in amd64 arm64; do
  name="syskit_${version#v}_linux_${arch}"
  CGO_ENABLED=0 GOOS=linux GOARCH="$arch" go build -trimpath -ldflags "-s -w -X github.com/Mersad-Moghaddam/syskit/internal/cli.version=$version" -o "$output_dir/$name" ./cmd/syskit
  tar -C "$output_dir" -czf "$output_dir/$name.tar.gz" "$name"
  rm "$output_dir/$name"
done
(cd "$output_dir" && sha256sum *.tar.gz > SHA256SUMS)
