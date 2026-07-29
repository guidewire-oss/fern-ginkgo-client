# safe-reporting-defaults - Design

## Overview

Three small, independent fixes in `pkg/client`, plus a README update. No new
architecture — this is hardening existing code paths.

## Implementation Details

1. **Default timeout** (`pkg/client/fern_api_client.go`, `New()`): set
   `httpClient.Timeout` to a sane default (30s) when constructing the default
   `http.Client`, before applying `ClientOption`s so `WithTimeout(...)` still
   overrides it.
2. **No panic on marshal failure** (`pkg/client/ginkgo_fern_reporter.go`,
   `Report()`): replace `panic(err)` with `return fmt.Errorf("client: failed to marshal test run: %w", err)`.
3. **Safe reporting helper** (new function in `pkg/client`, e.g.
   `ReportAfterSuiteSafe(projectID string, report gt.Report, opts ...ClientOption)`):
   calls `New` then `Report`, logging (`fmt.Printf`) and returning on any error
   instead of propagating it — mirrors the wrapper pattern consumers already
   hand-roll today.
4. **README**: replace the stale integration example with one that uses
   `ReportAfterSuiteSafe` and matches the real `New()` signature.

## API Changes

- New exported function: `ReportAfterSuiteSafe`.
- No changes to existing exported signatures (`New`, `Report`, `ClientOption`s).

## Data Model

No data model changes.
