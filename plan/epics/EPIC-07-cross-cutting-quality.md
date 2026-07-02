# EPIC-07 — Cross-Cutting Quality

> **Milestone:** all · **Sprints:** every sprint · **Points:** not separately scheduled (embedded in every story's Definition of Done)
> The quality concerns that apply to *every* story, tracked so they never become invisible or optional.

---

## Goal

Ensure that "done" always means tested, idiomatic, documented, and reviewed — regardless of which feature is being built. This epic is not a backlog of stories to schedule; it is the standing quality contract enforced by `../standards/definition-of-done.md` and `../specs/testing-strategy.md` inside every other story.

## Why it is an epic

Making quality a visible epic prevents the classic failure mode where tests, benchmarks, and docs are treated as "extra" and cut under pressure. Here they are first-class, named, and cross-referenced from every sprint's Definition of Done section.

## Standing concerns

Authoritative list: [`../product-backlog.md`](../product-backlog.md#epic-07--cross-cutting-quality-ongoing).

| ID | Concern | Applies to | Reference |
|---|---|---|---|
| XQ-01 | Unit tests (parse/transform/format/error paths). | every story | `../standards/definition-of-done.md` |
| XQ-02 | Fixture-backed collector tests. | collector stories | `../specs/testing-strategy.md` |
| XQ-03 | Linux integration tests (build-tagged). | collector stories | `../specs/testing-strategy.md` |
| XQ-04 | Golden-file output tests. | command stories | `../specs/testing-strategy.md` |
| XQ-05 | Benchmarks on hot paths. | parse/render stories | `../specs/testing-strategy.md` |
| XQ-06 | `gofmt`/`goimports`/`go vet`/`govulncheck` clean. | every story | `../standards/definition-of-done.md` |
| XQ-07 | Docs + CHANGELOG updated. | user-visible stories | `../standards/definition-of-done.md` |
| XQ-08 | Conventional Commits + squash-merge + rebase. | every PR | `../standards/branch-strategy.md` |

## How it is enforced

1. **Definition of Ready** blocks a story from starting until fixtures are identified and acceptance criteria are testable.
2. **Definition of Done** blocks a merge until the applicable XQ items are satisfied.
3. **CI** (`../.github/workflows/ci.yml`, after FND-10) runs fmt/vet/race/integration/coverage/bench/vuln on every PR — the mechanical half of the contract.
4. **Code review** (`../standards/code-review.md`) verifies the judgment half: are the *right* things tested, is the code idiomatic, does it satisfy the spec.

## Metrics

Tracked continuously in `../metrics-and-reporting.md` under quality metrics: unit+race suite, integration suite, golden diffs, benchmark trend, `govulncheck`, coverage. A regression in any is treated as a defect, not a warning.

## Definition of Done for the epic

This epic is never "closed" — it is satisfied continuously. It fails if any shipped story reached `main` without its applicable XQ items, which would be surfaced as an escaped defect (R-03 / quality metrics) and addressed in retrospective.
