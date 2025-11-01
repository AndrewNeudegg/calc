# Copilot Brief: Calc MVP

## Project Context
- Build a local-only, dependency-free, terminal notepad calculator in Go that mirrors Soulver style workflows.
- Each line is free-form text that may mix prose with values, units, dates, or variables; the right side shows the interpreted result, refreshed immediately on edits.
- The interface is keyboard-driven; parsing must be deterministic (no AI calls, no network).

## Primary Objectives
- Parse line-oriented input into tokens for arithmetic, units, dates/times, and named variables.
- Evaluate expressions in topological order so downstream lines update when dependencies change; detect and flag circular references.
- Render both the interpreted expression and formatted result, surfacing clear inline errors without crashing.
- Persist session-level settings (precision, date format, currency, etc.) and apply them consistently.

## Functional Requirements
- **Arithmetic & Variables**: Support operator precedence, parentheses, and named assignments reuseable across lines.
- **Unit Conversion**: Convert compatible SI units and user-defined compound units; reject incompatible conversions cleanly.
- **Custom Units & Rates**: Allow defining custom units/rates, preventing circular definitions.
- **Currencies**: Perform conversions via a bundled static table; block mixing currencies without explicit conversion when disabled.
- **Percentages**: Handle `X% of Y`, `increase/decrease X by Y%`, and keep `%` in outputs.
- **Date & Time Arithmetic**: Handle `today/tomorrow/yesterday`, natural weekdays, offsets, and rollovers.
- **Time Zones**: Convert between friendly place names and canonical IANA zones with DST awareness.
- **Fuzzy Phrases**: Map deterministic phrases (`half of X`, `twice X`, etc.) to arithmetic operations.
- **Totals & Functions**: Implement `sum`, `average`, `mean`, and natural-language equivalents like `total of`.
- **Command Mode**: Recognise `:` prefixed commands for save/open, settings, timezone lists, and rate loading.
- **Formatting Controls**: Honour global settings for numeric precision, currency, and date formatting.
- **Localisation**: Provide full localisation coverage for numeric, currency, date, and timezone presentation, including British and international formats, with illustrative examples captured in tests.
- **Error Handling**: Provide inline, non-fatal errors for syntax issues, unknown units, incompatible conversions, circular references, and invalid commands.

## Technical Constraints & Guidance
- Restrict dependencies to the Go standard library; absolutely no third-party modules, system calls to package managers, or hidden transitive dependencies.
- Maintain deterministic parsing—prefer explicit grammars/patterns over heuristic AI approaches.
- Keep the update loop typable: edits should re-use the parsed tree when possible and recompute only impacted lines.
- Use straightforward data structures (e.g., DAG for dependencies, environments for variables/units) to keep logic transparent.
- Ensure the command layer runs synchronously and surfaces confirmation messages or usage hints.
- Write all identifiers, comments, user-facing strings, and documentation in British English.

## Testing Expectations
- Derive table-driven tests from every example and acceptance test in the specification.
- Include coverage for parsing, evaluation, formatting, settings persistence, and error messaging.
- Validate timezone, currency, localisation variants, and unit conversions against the canonical values given in the spec.
- Maintain a comprehensive automated test suite that sustains >90% code coverage and passes `go test ./... -race`; include localisation-focused cases that exercise multiple locales and formatting rules.

## Deliverables
- Core evaluation engine with command-mode handling and display-ready formatting hooks.
- Comprehensive unit tests rooted in the specification’s canonical scenarios.
- Inline documentation for tricky parsing logic and subsystem boundaries.
