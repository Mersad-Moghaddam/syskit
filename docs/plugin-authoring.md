# Plugin Authoring

SysKit plugins are user-installed executables using JSON protocol v1. A plugin
directory contains `manifest.json` and the declared executable. Discovery and
inspection are read-only; execution occurs only through `plugins run`.

## Manifest

Declare `name`, `version`, `api_version`, `executable`, provided `collectors`,
required `permissions`, human-readable `output_schemas`, `author`, and
`license`. Executable paths are relative to the plugin directory and may not
escape it.

## Protocol

SysKit writes `{"api_version":"v1","action":"collect"}` to stdin. The plugin
must emit exactly one JSON value to stdout, write diagnostics to stderr, and
finish within the selected timeout. Do not render tables or terminal styling;
SysKit owns table, JSON, and YAML rendering.

## Example

```sh
cd examples/plugin-example
go build -o example-plugin .
cd ../..
go run ./cmd/syskit plugins inspect example --plugin-dir ./examples
go run ./cmd/syskit plugins run example --plugin-dir ./examples
```

Plugins execute with the invoking user's privileges. Protocol validation is
not a sandbox; users must review and trust plugin code before running it.
