# Cross-package test fixtures

This directory holds **shared, cross-package fixtures** — captured Linux kernel
interface files used by tests that span more than one package. Package-local
fixtures live under that package's own `testdata/` directory instead (for
example `internal/platform/testdata/` and `internal/collector/example/testdata/`);
put a fixture here only when it is genuinely shared.

Go excludes any directory named `testdata` from builds, so nothing here is
compiled into the binary.

## Layout

Fixtures mirror the real filesystem rooted at `/`, so a collector reading
`proc/stat` through a `platform.TestFS` rooted at a fixture set resolves to that
set's `proc/stat` file:

```text
testdata/
├── README.md                 ← this file
└── fixtures/
    └── <host-shape>/         ← one directory per captured host shape
        ├── SOURCE            ← provenance (see below)
        ├── proc/
        │   ├── uptime
        │   ├── loadavg
        │   ├── stat
        │   ├── meminfo
        │   └── self/mountinfo
        └── sys/
            └── devices/system/cpu/present
```

Name each fixture set for the **host shape or kernel behavior it represents**
(`naming-conventions.md`): `8-core-xeon`, `1-core-vm`, `cgroup-v1-host`,
`restricted-perms`, `truncated-meminfo`. A diverse corpus — many cores vs. one,
cgroup v1 vs. v2, an older kernel, malformed/truncated data — is how collectors
prove they handle real-world variation (`specs/collectors.md`, "Fixtures").

## Provenance: the `SOURCE` file

Every fixture set carries a `SOURCE` file recording where the bytes came from,
so a reviewer can tell whether a fixture reflects a real kernel and reproduce a
capture. Record at least:

- **kernel** — `uname -r`
- **distro** — `NAME` / `VERSION` from `/etc/os-release`
- **arch** — `uname -m`
- **container** — whether captured inside a container, and which runtime
- **date** — capture date (UTC)

For a hand-authored (non-captured) fixture, say so explicitly and state what it
is exercising, so it is never mistaken for a real capture.

## Capturing fixtures

Use [`scripts/capture-fixtures.sh`](../scripts/capture-fixtures.sh) to capture a
set from the current host. It is **read-only** with respect to the system — it
only reads `/proc` and `/sys` and writes into the target directory — and it
writes the `SOURCE` provenance file automatically:

```sh
scripts/capture-fixtures.sh testdata/fixtures/my-host
```

The `example-host` set below was authored by hand (synthetic bytes) purely to
document the layout; its `SOURCE` says so.
