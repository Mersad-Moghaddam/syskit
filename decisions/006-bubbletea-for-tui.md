# 006. Use Bubble Tea for the interactive dashboard TUI

**Status:** Accepted, 2026-07-01

---

## Context

The [v0.3 Real-Time Monitoring milestone](../specs/roadmap.md) introduces an
interactive terminal experience: `syskit dashboard` (a live metrics dashboard),
`syskit watch <command>` (continuous refresh), and `syskit top` (an interactive
process monitor with sorting and filtering). These require a terminal UI (TUI)
framework capable of a real-time refresh pipeline, keyboard navigation, and a
composable widget layout — needs that the one-shot table/JSON/YAML formatters of
v0.1–v0.2 do not cover.

The requirements:

- **A clear, testable update model.** Real-time UIs are notoriously stateful. We
  want state transitions that can be reasoned about and unit-tested
  (constitution principle 5, *Test Everything*).
- **Composability.** The dashboard is a layout of independent widgets (CPU,
  memory, network, processes), mirroring the independent-collector design of
  [ADR 004](./004-layered-architecture.md).
- **Go-native and idiomatic.** A pure-Go library keeps the single static binary
  intact ([ADR 001](./001-use-go.md)) and fits the *Clean Go* principle.
- **A styling companion.** The roadmap already pairs the TUI with a styling
  system for color, alignment, and layout (principle 9, *Consistent CLI
  Experience*).

The framework also must sit cleanly at the top of the layered architecture,
consuming the *same* services and collectors that the non-interactive CLI uses,
so that the dashboard is a second presentation over shared logic rather than a
parallel data path.

---

## Decision

We will use **Bubble Tea** (`github.com/charmbracelet/bubbletea`) for the
interactive dashboard and monitoring TUI, paired with **Lip Gloss**
(`github.com/charmbracelet/lipgloss`) for styling and layout. This is **planned
for the v0.3 dashboard** and is not a dependency of the v0.1–v0.2 command-line
surface.

Bubble Tea implements **The Elm Architecture**: application state is a `Model`,
events (key presses, timer ticks, resize) are processed by a pure `Update`
function that returns the next model and commands, and a `View` function renders
the model to a string. Real-time refresh is expressed as timer-driven messages
fed into `Update`.

Bubble Tea lives entirely in the **CLI Layer**
([ADR 004](./004-layered-architecture.md)). Its models call the existing service
layer for data and hold no data-collection logic. This keeps the TUI a thin
presentation over shared services, exactly as the non-interactive formatters are.

---

## Consequences

### Positive

- The Model/Update/View structure makes state transitions explicit and pure,
  which makes them **unit-testable** — `Update` can be driven with synthetic
  messages and its output asserted, without a real terminal.
- Composable models map naturally onto independent dashboard widgets, echoing the
  independent-collector architecture.
- Pure Go with no cgo, preserving the single static binary.
- Lip Gloss provides declarative styling (borders, colours, alignment, layout),
  supporting a consistent, polished terminal experience.
- The Charm ecosystem (Bubbles component library, Lip Gloss) is actively
  maintained and widely adopted for modern Go TUIs.

### Negative

- Adds dependencies (Bubble Tea, Lip Gloss, and their transitive packages) — the
  heaviest dependency addition in the project, accepted under the documented
  *Minimal Dependencies* exception because building a TUI runtime by hand would
  be unreasonable.
- The Elm architecture has a learning curve; contributors used to imperative UI
  code must adapt to message-driven state updates.
- The message-passing model can feel verbose for very simple screens.

### Neutral

- These dependencies are scoped to v0.3 and to the CLI layer; the v0.1–v0.2 core
  and all lower layers remain free of them.
- We adopt Charm's conventions for models, messages, and commands.

---

## Alternatives Considered

- **`rivo/tview`.** A batteries-included, widget-oriented TUI toolkit with ready
  tables, forms, and grids. Rejected because its imperative, callback-driven
  model is harder to unit-test than Bubble Tea's pure `Update` function and fits
  the *Test Everything* goal less cleanly; its higher-level widgets also offer
  less fine control over the custom dashboard layout we want.
- **`gizak/termui`.** Purpose-built for dashboards with charts and gauges out of
  the box. Rejected because it is less actively maintained, its architecture is
  chart-centric rather than a general application model, and it is a weaker fit
  for interactive, keyboard-driven navigation across `top`-style views.
- **Raw `termbox`/`tcell`.** Maximum control by driving the terminal directly.
  Rejected because it means building the event loop, layout, and rendering
  abstractions ourselves — reimplementing what Bubble Tea already provides and
  contradicting the *Minimal Dependencies* rationale for delegating unreasonable-
  to-build infrastructure. (Bubble Tea itself builds on tcell-class primitives, so
  we get that foundation without maintaining it.)

---

## References

- [Roadmap](../specs/roadmap.md) — v0.3 Real-Time Monitoring: "Bubble Tea
  integration for terminal UI", "Lip Gloss styling system"
- [Constitution](../specs/constitution.md) — principle 5 (Test Everything),
  principle 8 (Minimal Dependencies), principle 9 (Consistent CLI Experience)
- [ADR 001](./001-use-go.md) — pure-Go single static binary
- [ADR 004](./004-layered-architecture.md) — TUI confined to the CLI layer over shared services
- [ADR 005](./005-cobra-for-cli.md) — `dashboard`/`top`/`watch` are Cobra subcommands
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) and
  [The Elm Architecture](https://guide.elm-lang.org/architecture/)
- [Lip Gloss](https://github.com/charmbracelet/lipgloss)
