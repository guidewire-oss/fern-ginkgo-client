# safe-reporting-defaults - Requirements

## Introduction

Fern reporting is a side effect of running tests and should never be able to hang
or crash the suite it's reporting on. This closes [guidewire-oss/fern-ginkgo-client#39](https://github.com/guidewire-oss/fern-ginkgo-client/issues/39):
add a default HTTP timeout, stop panicking on marshal failure, and ship a
safe-by-default reporting helper so consumers stop hand-rolling log-and-continue
wrappers.

## Requirements

### Requirement 1: Default HTTP timeout

**User Story** As a consumer of the Fern Ginkgo client, I want `Report()` to fail fast on network problems, so that a Fern outage doesn't hang my whole test suite.

#### Acceptance Criteria

1. WHEN `client.New()` is called without `WithTimeout(...)` THE SYSTEM SHALL configure the internal `http.Client` with a default timeout.
2. WHEN a request in `Report()` exceeds the configured timeout THE SYSTEM SHALL return an error rather than blocking indefinitely.
3. IF a caller supplies `WithTimeout(...)` THEN THE SYSTEM SHALL use that value instead of the default.

### Requirement 2: No panics on marshal failure

**User Story** As a consumer, I want `Report()` to never panic, so that a serialization edge case doesn't crash my test binary.

#### Acceptance Criteria

1. WHEN `json.Marshal(testRun)` fails inside `Report()` THE SYSTEM SHALL return an error instead of panicking.
2. THE SYSTEM SHALL wrap the underlying marshal error (`%w`) so callers can inspect it.

### Requirement 3: Safe-by-default reporting helper

**User Story** As a consumer, I want a ready-made safe reporting entry point, so that I don't have to hand-roll a log-and-continue wrapper.

#### Acceptance Criteria

1. THE SYSTEM SHALL provide a helper (e.g. `client.ReportAfterSuiteSafe`) that wraps `New` + `Report`.
2. WHEN `New` or `Report` returns an error THE SYSTEM SHALL log a warning and return normally, without propagating the error or failing the suite.
3. THE README SHALL demonstrate this pattern with a signature that matches the real `New() (*FernApiClient, error)` API.

### Requirement 4: Non-Functional Requirements

**User Story** As a consumer, I want this change to be backward compatible, so that upgrading doesn't break my existing integration.

#### Acceptance Criteria

1. THE SYSTEM SHALL preserve existing exported signatures of `New()` and `Report()`.
2. WHILE no `WithTimeout(...)` is supplied THE SYSTEM SHALL apply the new default timeout without requiring any consumer code changes.

## Constraints

- No breaking changes to existing exported function signatures.
- No new third-party dependencies.
