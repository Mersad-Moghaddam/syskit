# Definition of Done

> The checklist a change must satisfy before it can merge to `main`.

---

## Purpose

A feature is not "done" when the code compiles. It is done when it meets the spec, is tested, is idiomatic, is documented, and has passed review. This standard makes that bar explicit and non-negotiable.

Use this list as the final gate before requesting merge. It complements `definition-of-ready.md` (the gate for *starting* work) and `code-review.md` (how the change is verified).

---

## The Checklist

A change is Done only when every item is true:

- [ ] **Spec satisfied & acceptance criteria met** — The change fulfills its spec in `../specs/` and every acceptance criterion is demonstrably met.
- [ ] **Code complete & idiomatic** — No stubs, unresolved task markers, or dead code. Follows `coding-conventions.md` and `naming-conventions.md`.
- [ ] **Unit tests pass** — New behavior has unit tests, including error paths and edge cases; the full suite is green.
- [ ] **Integration tests pass** — Collectors that read real Linux interfaces have integration coverage, and it passes on a real Linux system.
- [ ] **Benchmarks run for hot paths** — Performance-sensitive paths have `BenchmarkXxx` functions, and results show no unexplained regression.
- [ ] **`go vet` & gofmt clean** — `gofmt -l .`, `goimports`, and `go vet ./...` produce no output.
- [ ] **No new advisories** — `govulncheck ./...` is clean; any new dependency has an ADR (`dependency-policy.md`).
- [ ] **Docs updated** — Specs, command help text, and relevant `../docs/` are updated to match the change.
- [ ] **CHANGELOG updated** — A user-facing entry is added under the correct heading (`commit-conventions.md`).
- [ ] **Reviewed & approved** — At least one maintainer has approved (`code-review.md`).
- [ ] **CI green** — All CI checks pass on the latest commit.

---

## Why Each Item

| Item | Why it matters |
|---|---|
| Spec satisfied | We build to specification (Documentation First). Code that drifts from its spec is a defect even if it "works." |
| Idiomatic code | Consistency keeps the codebase reviewable and lowers the cost of every future change. |
| Unit tests | Prove the logic is correct in isolation and lock in behavior against regressions. |
| Integration tests | Kernel interfaces vary; only a real system proves a collector reads them correctly. |
| Benchmarks | A monitoring tool must not distort what it measures; hot paths are tracked over time. |
| vet/fmt clean | Mechanical checks catch mechanical mistakes so reviewers focus on logic. |
| No advisories | Every dependency and its vulnerabilities are our responsibility. |
| Docs updated | Undocumented behavior is unusable behavior; docs are a first-class deliverable. |
| CHANGELOG | Users and the release process depend on an accurate history of changes. |
| Reviewed & approved | The team, not just the author, owns every line on `main`. |
| CI green | The stable-`main` guarantee is only real if it is enforced automatically. |

---

## When It Is *Not* Done

Do not merge — and do not mark a task complete — if any of these are true:

- Tests are failing or missing for new behavior.
- The implementation is partial or gated behind an unfinished branch.
- `go vet`, gofmt, or `govulncheck` report findings.
- The CHANGELOG or affected docs are stale.
- A new dependency lacks an ADR.

If you are blocked from finishing an item, keep the work open and record what is needed, rather than lowering the bar.

---

*Done means a future contributor can build on this change without discovering it was never finished.*
