# Definition of Ready

> The checklist a task must satisfy before implementation begins.

Definition of Ready protects contributors from starting work while important decisions are still missing. A task that is not ready should remain in planning, even if the implementation seems straightforward.

## Ready Checklist

A feature or implementation task is Ready only when:

- [ ] The relevant feature spec exists and is reviewed.
- [ ] User story and motivation are clear.
- [ ] Expected CLI and flags are documented.
- [ ] Expected output is documented for table and structured formats.
- [ ] Linux data sources are identified.
- [ ] Edge cases are listed.
- [ ] Acceptance criteria are testable.
- [ ] Required fixtures are identified.
- [ ] Dependencies and new libraries are reviewed.
- [ ] Security, permission, and partial-data behavior are understood.
- [ ] Documentation impact is known.

## Not Ready Examples

- "Implement network support" without deciding whether Netlink is required.
- "Show memory usage" without defining how cache and available memory are represented.
- "Add JSON output" without defining field names and units.
- "Create process tree" without handling disappearing processes.

## Ready For Planning

Some tasks are not ready for implementation but are ready for design work. In that case, open a design proposal issue and identify the missing decisions.

## Relationship To Done

Ready is the entry gate. Done is the exit gate. A task can be ready and still fail later review if the implementation does not satisfy the spec, tests, or documentation requirements.
