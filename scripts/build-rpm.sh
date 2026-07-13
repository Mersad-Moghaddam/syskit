#!/usr/bin/env bash
set -euo pipefail

version="${1:-}"
arch="${2:-amd64}"
output_dir="${3:-dist}"
if [[ -z "$version" ]]; then echo "usage: $0 <version> [amd64|arm64] [output-dir]" >&2; exit 2; fi
case "$version" in v[0-9]*) ;; *) echo "version must start with v" >&2; exit 2;; esac
case "$arch" in
  amd64) rpm_arch="x86_64" ;;
  arm64) rpm_arch="aarch64" ;;
  *) echo "architecture must be amd64 or arm64" >&2; exit 2 ;;
esac
if ! command -v rpmbuild >/dev/null 2>&1; then
  echo "rpmbuild is required to create RPM packages" >&2
  exit 1
fi

export SOURCE_DATE_EPOCH="${SOURCE_DATE_EPOCH:-$(git log -1 --format=%ct)}"
stage="$(mktemp -d)"
trap 'rm -rf "$stage"' EXIT
topdir="$stage/rpmbuild"
install -d "$topdir/BUILD" "$topdir/BUILDROOT" "$topdir/RPMS" \
  "$topdir/SOURCES" "$topdir/SPECS" "$topdir/SRPMS"

CGO_ENABLED=0 GOOS=linux GOARCH="$arch" go build -trimpath \
  -ldflags "-s -w -X github.com/Mersad-Moghaddam/syskit/internal/cli.version=$version" \
  -o "$topdir/SOURCES/syskit" ./cmd/syskit
install -m 0644 LICENSE "$topdir/SOURCES/LICENSE"

cat > "$topdir/SPECS/syskit.spec" <<EOF
%global debug_package %{nil}
%global _build_id_links none
%global __strip /bin/true
Name: syskit
Version: ${version#v}
Release: 1%{?dist}
Summary: Linux-native system inspection and diagnostics toolkit
License: MIT
URL: https://github.com/Mersad-Moghaddam/syskit
Source0: syskit
Source1: LICENSE

%description
SysKit reads procfs, sysfs, Netlink, and cgroups directly without shelling out.

%prep

%build

%install
install -Dpm0755 %{SOURCE0} %{buildroot}%{_bindir}/syskit
install -Dpm0644 %{SOURCE1} %{buildroot}%{_licensedir}/syskit/LICENSE

%files
%{_bindir}/syskit
%license %{_licensedir}/syskit/LICENSE
EOF

rpmbuild -bb --target "$rpm_arch" --define "_topdir $topdir" \
  --define "_buildhost syskit-build" \
  --define "use_source_date_epoch_as_buildtime 1" \
  --define "clamp_mtime_to_source_date_epoch 1" \
  "$topdir/SPECS/syskit.spec"
package="$(find "$topdir/RPMS" -type f -name '*.rpm' -print -quit)"
if [[ -z "$package" ]]; then
  echo "rpmbuild did not produce an RPM package" >&2
  exit 1
fi
mkdir -p "$output_dir"
install -m 0644 "$package" "$output_dir/syskit-${version#v}-1.$rpm_arch.rpm"
