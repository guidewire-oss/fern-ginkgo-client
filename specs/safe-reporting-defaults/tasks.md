# safe-reporting-defaults - Tasks

## Implementation Tasks

- [x] 1. Test + implement default HTTP timeout in `client.New()` (Requirement 1)
- [x] 2. Test + implement non-panicking marshal error in `Report()` (Requirement 2)
- [x] 3. Test + implement `ReportAfterSuiteSafe` helper (Requirement 3)
- [x] 4. Update README integration example to use the safe helper and correct `New()` signature (Requirement 3.3)

## Verification Tasks

- [x] 5. Run full test suite, confirm green
- [x] 6. `/spec:verify` against requirements.md
