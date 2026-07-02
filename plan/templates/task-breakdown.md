# Template: Task Breakdown

> Copy this when decomposing a committed story into tasks at Sprint Planning. Tasks are the developer-level checklist; each should be ≲ 1 day and independently reviewable.

---

```markdown
## <ID> — <story title>  (<points> pts)

Decompose along the architecture layers so each task is small and reviewable
(dependency direction: CLI → Command → Service → Collector → Platform → kernel).

- [ ] platform: <read source(s) via SysFS> (+ capture fixtures)
- [ ] collector: <parse/normalize into typed struct> (+ unit tests, error paths)
- [ ] service: <aggregate / derive metrics / filter / sort> (+ unit tests)
- [ ] command: <flags, validation, wiring>
- [ ] render: <table columns + structured fields> (+ golden tests)
- [ ] integration: <build-tagged Linux test on real /proc or /sys>
- [ ] benchmark: <hot path Benchmark + b.ReportAllocs()>  (if perf-sensitive)
- [ ] docs: <command help text + relevant ../docs update>
- [ ] changelog: <user-facing entry under the right heading>
```

---

## Guidance

- **Order by dependency.** Platform/collector tasks unblock service/render tasks; do them first.
- **One task, one concern.** If a task mixes parsing and rendering, split it.
- **Tests are tasks, not afterthoughts.** They are inside the story's estimate (`../estimation-and-velocity.md`) and inside the Definition of Done.
- **Not every task applies.** A pure formatter story has no collector/platform tasks; a docs-only story has no benchmark. Delete what doesn't fit.
- **Surface unknowns.** A task nobody can size is a hidden unknown — raise it at planning or spike it (`../risk-register.md`).

## Example (filled)

```markdown
## MEM-01 — syskit memory  (8 pts)

- [ ] platform: read /proc/meminfo, /proc/vmstat via SysFS (+ fixtures: standard, legacy no-MemAvailable, high-pressure)
- [ ] collector: parse meminfo fields, handle missing MemAvailable → ErrFieldMissing (+ table-driven unit tests)
- [ ] service: derive used/available/pressure, keep raw counters separate from derived
- [ ] command: `syskit memory` wiring, --format flag
- [ ] render: table (total/used/free/available/swap) + JSON fields with units (+ golden)
- [ ] integration: //go:build linux && integration — assert non-zero total on real host
- [ ] benchmark: BenchmarkParseMeminfo with ReportAllocs
- [ ] docs: command help + getting-started example
- [ ] changelog: "add `memory` command"
```
