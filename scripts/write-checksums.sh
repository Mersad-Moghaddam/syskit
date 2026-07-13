#!/usr/bin/env bash
set -euo pipefail

artifact_dir="${1:-dist}"
shopt -s nullglob
artifacts=("$artifact_dir"/*.tar.gz "$artifact_dir"/*.deb "$artifact_dir"/*.rpm)
if (( ${#artifacts[@]} == 0 )); then
  echo "no release artifacts found in $artifact_dir" >&2
  exit 1
fi

for i in "${!artifacts[@]}"; do
  artifacts[$i]="${artifacts[$i]##*/}"
done
(
  cd "$artifact_dir"
  printf '%s\n' "${artifacts[@]}" | sort | while IFS= read -r artifact; do
    sha256sum "$artifact"
  done
) > "$artifact_dir/SHA256SUMS"
chmod 0644 "$artifact_dir/SHA256SUMS"
