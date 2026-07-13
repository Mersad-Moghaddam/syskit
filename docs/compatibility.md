# v1 Compatibility Contract

> The public surfaces SysKit commits to preserving after v1.0.0.

## Stability Boundary

SysKit is distributed as a command-line application, not an importable Go
library. Its v1 public contract consists of:

- command paths, documented flags, positional argument shapes, and defaults;
- table/JSON/YAML format selection and the JSON-equivalent field schemas;
- exit-code meanings;
- configuration keys, environment variables, and their precedence;
- plugin manifest fields and the plugin protocol API version.

Internal Go packages, collector implementation details, kernel read order, log
wording, and human-oriented table spacing are not public APIs. Table columns may
evolve compatibly; automation must use JSON or YAML.

The contract is binding from v1.0.0 onward. These files describe the stable
surface and allow CI to detect unintentional drift.

## Machine-Readable Manifests

[v1-cli.json](../contracts/v1-cli.json) is the canonical inventory of commands,
flags, formats, configuration/environment names, exit codes, and accepted plugin
API versions. [v1-schemas.json](../contracts/v1-schemas.json) records every field
reachable from built-in structured command results, its JSON type, and whether
the field is required or may be omitted.

CI constructs the Cobra command tree and reflects the typed output models, then
compares both with these manifests. Removing or renaming a surface, changing a
field type, or making an optional field required therefore fails tests. A
backward-compatible addition also requires an intentional manifest update and a
SemVer/CHANGELOG review.

Regenerate the schema manifest only when a reviewed contract change is intended:

```sh
UPDATE_CONTRACT=1 go test ./internal/cli -run TestV1StructuredOutputContract
git diff -- contracts/v1-schemas.json
```

The CLI manifest is edited deliberately rather than generated so reviewers can
compare the intended command surface with Cobra's actual tree.

## Command-to-Schema Mapping

| Command | Structured result |
|---|---|
| `system` | `model.SystemInfo` |
| `cpu` | `model.CPUInfo` |
| `memory` | `model.MemoryInfo` |
| `disk`, `filesystem` | `model.DiskInfo` |
| `process` | `model.ProcessList` |
| `process tree` | array of `service.ProcessTreeNode` |
| `network` | `model.NetworkInfo` |
| `network interfaces` | array of `model.NetworkInterface` |
| `network routes` | array of `model.Route` |
| `network dns` | array of strings |
| `ports` | `model.PortInfo` |
| `containers` | `model.ContainerList` |
| `containers inspect` | `model.ContainerDetail` |
| `diagnostics` | `model.DiagnosticReport` |
| `plugins list` | array of `plugin.Info` |
| `plugins inspect` | `plugin.Info` |
| `plugins run` | plugin-declared output schema |

Interactive `dashboard`, `top`, and `watch` views do not expose a structured
output schema. `version` emits one version string. Completion subcommands emit
shell source intended for the selected shell.

## Compatible and Breaking Changes

Within the v1 major line, adding a command, flag, optional structured field, or
plugin API version is backward-compatible and requires a MINOR release. Fixing
incorrect values without changing their shape is a PATCH change.

Removing or renaming a command/flag/field, changing a structured field's type or
requiredness, reassigning an exit code, or rejecting an existing plugin API is
breaking. It follows the deprecation process in the
[versioning standard](../standards/versioning.md) and requires a MAJOR release.
