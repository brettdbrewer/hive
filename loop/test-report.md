# Test Report: MCP knowledge search blackout — iteration 388

**Timestamp:** 2026-03-28

## What Was Built
No new code shipped. Builder confirmed acceptance criteria (knowledge_search returns ≥1 result
for "lesson") was already met by prior commits (`90121a9`, `3b6cd0e`).

## Scope
Focus: `cmd/mcp-knowledge` — 18 tests covering the search-blackout fix.

## Results

| Suite | Tests | Result |
|-------|-------|--------|
| `cmd/mcp-knowledge` | 18 (fresh run) | **PASS** |
| `cmd/mcp-graph` | cached | PASS |
| `cmd/post` | cached | PASS |
| `pkg/api` | cached | PASS |
| `pkg/authority` | cached | PASS |
| `pkg/hive` | cached | PASS |
| `pkg/loop` | cached | PASS |
| `pkg/resources` | cached | PASS |
| `pkg/runner` | cached | PASS |
| `pkg/workspace` | cached | PASS |

All 18 `mcp-knowledge` tests ran fresh (`-count=1`), all pass.

## Coverage — key behaviours verified

| Behaviour | Test |
|-----------|------|
| claims.md indexed when present | `TestBuildHiveLoopIncludesClaimsWhenPresent` |
| claims.md absent → not in tree | `TestBuildHiveLoopOmitsClaimsWhenAbsent` |
| search finds claims content | `TestHandleSearchFindsClaims` |
| search finds claims past 4000-char window | `TestHandleSearchFindsDeepClaims` |
| individual claim retrievable by slug ID | `TestHandleGetIndividualClaim` |
| duplicate titles get unique IDs (-2, -3) | `TestParseClaimsDuplicateTitles` |
| slug truncates at 60 chars, no trailing hyphen | `TestClaimSlugTruncation` |
| slug collapses special chars | `TestClaimSlugSpecialChars` |
| summary truncates at 120 chars | `TestClaimSummaryLongLine` |
| all-metadata body → empty summary | `TestClaimSummaryAllMetadata` |
| empty file → no panic | `TestParseClaimsEmptyFile` |
| no ## sections → no claims | `TestParseClaimsNoSections` |
| search result cap ≤10 | `TestHandleSearchResultCap` |
| empty query → error, no panic | `TestHandleSearchEmptyQuery` |
| empty id → error, no panic | `TestHandleGetEmptyID` |
| claim children visible in topics listing | `TestClaimChildrenVisibleInTopics` |
| loop children include claims.md | `TestHandleTopicsReturnsLoopChildren` |
| handleGet returns full claim content | `TestHandleGetClaims` |

## Gaps / follow-up

The residual gap noted in build.md (65/145 claims synced — Lessons 1–108 absent because
`syncClaims` only queries board, not knowledge lens) is a Scout-level gap. Tests for that
sync path would require an integration test against a live API, which is out of scope.

## Verdict
**PASS.** Acceptance criteria met. No regressions. @Critic
