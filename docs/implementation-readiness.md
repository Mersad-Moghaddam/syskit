# Implementation Readiness

> Checklist for deciding when SysKit is ready to move from planning into production Go code.

Implementation should begin only after the foundational decisions are stable enough that the first code does not have to invent product behavior, architecture, or process from scratch.

## Required Before First Code

- [ ] Product scope and non-goals are reviewed.
- [ ] Architecture layers and dependency direction are accepted.
- [ ] CLI conventions are accepted.
- [ ] Core feature specs for v0.1 are accepted.
- [ ] Collector architecture is accepted.
- [ ] Rendering architecture is accepted.
- [ ] Error handling, logging, and configuration specs are accepted.
- [ ] Testing strategy is accepted.
- [ ] Dependency policy is accepted.
- [ ] Branch, commit, review, and versioning standards are accepted.
- [ ] Planning-phase CI passes.

## v0.1 Implementation Entry Criteria

The first implementation pull request should be limited to foundation work:

- Create the Go module.
- Establish the approved repository layout.
- Add minimal CLI bootstrapping.
- Add test and lint workflows appropriate for Go code.
- Add no user-facing feature until the foundation compiles and tests.

## Feature Entry Criteria

Before a feature enters implementation:

- Its feature spec exists in `specs/features/`.
- The expected CLI is defined.
- Output examples are defined for table and structured formats.
- Edge cases are listed.
- Acceptance criteria are testable.
- Linux data sources are identified.
- Required fixtures are planned.
- Any new dependency has been reviewed against the dependency policy.

## Feature Exit Criteria

A feature is complete only when:

- The implementation follows the architecture boundaries.
- Unit tests cover parsing, transformation, and errors.
- Integration tests cover live Linux invariants where appropriate.
- Golden output tests cover user-facing output.
- Documentation is updated.
- Changelog notes are prepared when user-visible behavior changes.
- Reviewers agree the feature satisfies its acceptance criteria.

## Explicit Non-Readiness Signals

Do not begin implementation if:

- A feature still has unresolved command naming.
- A data source is unknown or only assumed.
- Output schema fields are unclear.
- Error behavior is not specified.
- Tests require reading the developer's live host state directly.
- The change needs privileged behavior that violates SysKit's read-only scope.

## Transition Pull Request

The transition from planning to implementation should be a dedicated pull request. It should update the CI workflow, README status, and contributing guide so future contributors are not working under planning-phase rules.
