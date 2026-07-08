// Package collector and its subpackages parse one domain's raw kernel bytes
// (obtained through the injected platform.SysFS) into typed model structs.
//
// Each domain (cpu, memory, disk, process, network, ports, fs) is an
// independent collector that knows nothing about the others. Collectors do not
// log, render, read configuration, or shell out (ADR-003); they return typed
// data and domain sentinel errors. Per ADR-004 a collector may import platform
// and model, but never service, command, cli, or render.
package collector
