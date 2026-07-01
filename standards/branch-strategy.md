# Branch Strategy

> How branches are named, used, and merged to keep `main` stable and history linear.

---

## Purpose

SysKit uses a trunk-based workflow with short-lived feature branches. The goal is a clean, linear, releasable history on `main` at all times. This standard defines branch naming, lifecycle, and merge policy.

It works together with `commit-conventions.md` (what commits look like) and `code-review.md` (how changes are approved before merge).

---

## Core Rules

- **`main` is always stable and releasable.** Every commit on `main` passes CI, is reviewed, and could be tagged and shipped.
- **No direct commits to `main`.** All changes arrive through a reviewed pull request.
- **Branches are short-lived.** A feature branch exists for hours or days, not weeks. Long-lived branches accumulate merge debt and drift from `main`.
- **History is linear.** No merge commits from feature branches; we rebase and squash.

---

## Branch Naming

Feature branches use a `type/kebab-description` form. The type matches the primary commit type (see `commit-conventions.md`).

| Prefix | Use for | Example |
|---|---|---|
| `feat/` | New feature or command | `feat/cpu-frequency` |
| `fix/` | Bug fix | `fix/memory-available-calc` |
| `docs/` | Documentation change | `docs/standards-set` |
| `chore/` | Maintenance, deps, tooling | `chore/bump-go-1.23` |

Rules:

- Description is lowercase, kebab-case, and concise: `feat/process-tree`, not `feat/Process_Tree_Visualization`.
- Include an issue number when it aids tracking: `fix/58-slab-double-count`.
- One logical change per branch. Split unrelated work into separate branches.

---

## Keeping History Linear

Rebase your branch onto the latest `main` before opening or updating a PR. Do **not** merge `main` into your branch.

```bash
git fetch origin
git rebase origin/main
# resolve conflicts, then:
git push --force-with-lease
```

Use `--force-with-lease`, never a bare `--force`, so you cannot clobber someone else's push.

---

## Merge Policy

- **Squash-merge only.** A feature branch collapses into a single commit on `main`. This keeps `main` history one-change-per-commit and lets contributors commit freely (including `wip`) on their branch.
- The squash commit's **title and body must follow Conventional Commits** — this is the message that feeds the CHANGELOG and version selection.
- Delete the branch immediately after merge.
- Requirements before merge (see `code-review.md`): CI green, at least one maintainer approval, branch rebased on current `main`.

---

## Tags and Releases

- Releases are marked with annotated, SemVer tags on `main`: `v0.1.0`, `v0.2.0` (see `versioning.md`).
- Tags are created only on `main`, only after CI is green.

```bash
git tag -a v0.1.0 -m "SysKit v0.1.0 — Foundation"
git push origin v0.1.0
```

- No release branches while the project is small; `main` is the release source. A `release/x.y` branch may be introduced later only to backport fixes to an older line, and that decision will be recorded in `../decisions/`.

---

## Branch Lifecycle

```text
                 rebase onto main
                 (as needed)
   origin/main ──●──────●──────────────●───────────●──▶  (always stable)
                  \                    ▲            ▲
                   \  feat/cpu-freq    │ squash-merge│ tag v0.1.0
                    ●──●──●──●──────────┘            │
                    create  work  PR + review + CI  │
                    branch  commits  green          │
                                                    │
                    (branch deleted after merge) ───┘
```

1. **Branch** from current `origin/main` with a `type/description` name.
2. **Commit** freely; intermediate messages need not be conventional.
3. **Rebase** onto `main` to stay current.
4. **Open PR**; pass review and CI.
5. **Squash-merge** with a Conventional Commit message.
6. **Delete** the branch.
7. **Tag** on `main` when a milestone is ready.

---

*Short branches and linear history are not bureaucracy — they are what make `main` trustworthy.*
