# Code Review

> What every pull request is checked against, and how authors and reviewers work together.

---

## Purpose

Code review protects the quality bar defined in `../specs/constitution.md`. Every change to `main` is reviewed. Review is not a formality — it is where correctness, tests, and constitution alignment are verified before code becomes permanent.

This standard defines what reviewers check, what authors owe reviewers, and the norms for giving and receiving feedback. It complements `definition-of-done.md`, which lists what must be *true* before a PR is opened.

---

## What Reviewers Check

Reviewers evaluate a PR across these dimensions:

| Dimension | Questions |
|---|---|
| **Correctness** | Does it do what the spec says? Are edge cases, empty inputs, and error paths handled? |
| **Tests** | Are there unit tests? Table-driven where appropriate? Do they cover failure modes, not just the happy path? |
| **Constitution alignment** | Linux-first, native-APIs-first, modular, idiomatic Go? No new panics in library code? (`../specs/constitution.md`) |
| **Performance** | Any needless allocations or I/O in a hot path? Are benchmarks present for hot paths? |
| **Style** | gofmt/goimports clean, naming per `naming-conventions.md`, small focused packages? |
| **Dependencies** | Any new dependency? If so, is there an ADR? (`dependency-policy.md`) |
| **Docs** | Are specs, help text, and CHANGELOG updated to match the change? |

---

## Reviewer Checklist

- [ ] PR description explains the *what* and *why* and links the relevant spec/issue.
- [ ] Change is scoped to one logical concern.
- [ ] CI is green (build, `go vet`, gofmt, tests, govulncheck).
- [ ] Tests exist and actually exercise the new behavior, including error paths.
- [ ] No new external dependency without an accompanying ADR.
- [ ] No mutable global state introduced (`coding-conventions.md`).
- [ ] Errors are wrapped with context; no swallowed errors; no library panics.
- [ ] Public surface is minimal and named correctly.
- [ ] CHANGELOG and affected docs are updated.
- [ ] The squash-merge title is a valid Conventional Commit.

---

## Author Responsibilities

- **Keep PRs small.** Aim for a reviewable unit — ideally under ~400 changed lines. Split large work into stacked or sequential PRs.
- **Self-review first.** Read your own diff before requesting review. Remove debug code, resolve unresolved task markers, and check the checklist above.
- **Write a real description.** State the problem, the approach, and anything you want reviewers to scrutinize. Link the spec.
- **Make CI green before requesting review.** Do not ask a human to find what a linter would.
- **Respond, don't defend.** Treat comments as questions about the code, and reply in the thread or with a follow-up commit.

---

## Approval Requirements

A PR may be squash-merged only when all of the following hold:

- **At least one maintainer approval.**
- **CI is green** on the latest commit.
- **Branch is rebased** on current `main` (`branch-strategy.md`).
- All review threads are resolved or explicitly deferred with a tracked follow-up.

---

## Response-Time Expectations

| Actor | Expectation |
|---|---|
| Reviewer | First response within **2 business days** of a review request. |
| Author | Address feedback or reply within **2 business days** of receiving it. |
| Everyone | Flag urgent changes (build broken, security) explicitly for same-day attention. |

Stale PRs (no activity for two weeks) may be closed and reopened when work resumes.

---

## Feedback Norms

**Giving feedback:**

- Comment on the code, not the person. "This can drop the extra allocation," not "you allocated too much."
- Distinguish severity. Prefix non-blocking suggestions with `nit:`; reserve blocking comments for correctness, tests, and constitution violations.
- Explain the *why* and, where useful, suggest the fix.
- Approve when the code is good enough to ship, not only when it is exactly how you would have written it.

**Receiving feedback:**

- Assume good intent; the reviewer is protecting the codebase, not attacking you.
- If you disagree, say so with reasoning — reviews are a discussion, not a verdict.
- Prefer resolving a thread with a commit over a long argument.

---

*Review is how the whole team, not just the author, becomes responsible for every line on `main`.*
