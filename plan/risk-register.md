# Risk Register

> Known risks to delivering SysKit v1.0, with likelihood, impact, mitigation, and owner. Reviewed at every sprint retrospective; new risks enter here as they are discovered.

---

## Scoring

- **Likelihood / Impact:** Low (L) · Medium (M) · High (H).
- **Exposure** = Likelihood × Impact, used to order the register (highest first).
- **Owner:** the role accountable for the mitigation, not necessarily who executes it.

---

## Active risks

| ID | Risk | Likelihood | Impact | Exposure | Mitigation | Owner |
|---|---|---|---|---|---|---|
| R-01 | **Netlink integration harder than estimated** (NET-01, 13 pts) blows Sprint 5. | H | H | **H** | Time-box a spike in Sprint 4 refinement; keep NET-02/PRT-01 deferrable; capture protocol quirks in a learning note. | Product Owner |
| R-02 | **Plugin protocol (PLG-01) churns** the whole extensibility epic. | M | H | **H** | ADR-007 already commits to out-of-process; prototype the wire protocol as a spike before Sprint 12; freeze protocol before PLG-02. | Maintainer |
| R-03 | **Fixture drift** — captured `/proc`/`/sys` fixtures diverge from real kernels, masking bugs. | M | H | **H** | Integration tests on real Linux CI (mandatory in DoD); `scripts/capture-fixtures.sh` records provenance; multiple fixture sets per collector. | Whole team |
| R-04 | **Contributor availability drops** (open-source volatility) → velocity collapse. | H | M | **M** | Forecast by rolling velocity, not hope; keep stories small and independent; maintain 1.5 sprints of Ready work. | Scrum Master |
| R-05 | **TUI concurrency races** in refresh pipeline (RT-01, DSH-01, TOP-01). | M | H | **M** | All tests under `-race`; design refresh with channels + single owner of state; review against testing-strategy. | Developers |
| R-06 | **Scope creep** from post-1.0 "Future Considerations" leaking into sprints. | M | M | **M** | Non-goals defended by Product Owner; future items stay unscheduled in backlog; ADR required to promote. | Product Owner |
| R-07 | **CI cannot exercise all kernel variants** (cgroup v1 vs v2, VM-hidden topology). | M | M | **M** | Curate diverse fixtures; document unsupported/edge kernels; assert on invariants not exact values. | Developers |
| R-08 | **cgroup v1/v2 divergence** (CG-01) doubles container work. | M | M | **M** | Isolate the difference inside the platform layer; fixtures for both; single normalized model upward. | Developers |
| R-09 | **Golden-file churn** — cosmetic output tweaks break many golden tests, review fatigue. | M | L | **L** | Stabilize output contract early (FND-06); regenerate goldens intentionally with reviewed diffs. | Developers |
| R-10 | **Dependency risk** — a new library introduces a vulnerability or unmaintained code. | L | M | **L** | `govulncheck` in CI; every new dep needs an ADR (`../standards/dependency-policy.md`); prefer stdlib. | Maintainer |
| R-11 | **Estimation immaturity** early — first sprints mis-sized. | H | L | **L** | Treat Sprint 0–1 velocity as warm-up; re-forecast from Sprint 3; reference stories for calibration. | Scrum Master |
| R-12 | **Packaging surprises** (deb/rpm/AUR) surface late in Sprint 14. | M | M | **M** | Dry-run one package format during v0.3 slack; don't leave all packaging to the final sprint. | Maintainer |

---

## Spikes (risk-reduction time-boxes)

Some risks are best retired by a **spike**: a short, time-boxed investigation whose output is knowledge (a note or a throwaway prototype), not shippable code.

| Spike | Retires | When | Timebox | Output |
|---|---|---|---|---|
| SPK-NET | R-01 | Sprint 4 refinement | 1 day | Netlink message-parse prototype + learning note; re-point NET-01 if needed. |
| SPK-PLG | R-02 | Sprint 11 refinement | 1 day | Out-of-process protocol sketch; confirm ADR-007 assumptions. |
| SPK-PKG | R-12 | Sprint 9 slack | ½ day | One working `.deb` build to de-risk PKG-01. |

Spikes are added to the sprint backlog like any story but are pointed for the *investigation*, and always produce a written result.

---

## Retired risks

_(none yet — move a risk here with the retiring sprint and evidence once its exposure drops to negligible.)_

---

## How this register is used

1. **Retrospective:** re-score active risks; retire or add.
2. **Planning:** if a committed story touches a High-exposure risk, confirm its mitigation is scheduled (often a spike).
3. **Cross-link:** high-exposure risks are referenced from the sprint files that carry them.
