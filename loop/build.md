# Build Report — Iteration 199

## Test Debt Paydown

**6 new test functions** covering the Social layer sprint (iters 189-198):

| Test | What | Type |
|------|------|------|
| `TestFollows` | Follow/Unfollow/IsFollowing/CountFollowers/CountFollowing/ListFollowedIDs | Store (DB) |
| `TestReposts` | Repost/Unrepost/HasReposted/GetBulkRepostCounts/GetBulkUserReposts | Store (DB) |
| `TestQuotePost` | CreateNode with QuoteOfID, GetNode resolves quote_of_author/title/body | Store (DB) |
| `TestMessageSearch` | SearchMessages body filter, from: filter, no-match case | Store (DB) |
| `TestBulkEndorsements` | GetBulkEndorsementCounts, GetBulkUserEndorsements on posts (not users) | Store (DB) |
| `TestParseMessageSearch` | Pure function: operator parsing, 6 cases | Handler (no DB) |

**Coverage:** Follows the lesson 42 ratio (1 test iter per ~5 feature iters). Covers the 5 most critical new features. Feed algorithm tests deferred (timing-dependent scoring is hard to assert deterministically).

**Files changed:**
- `graph/store_test.go` — 5 new test functions (~200 lines)
- `graph/handlers_test.go` — `TestParseMessageSearch` (~25 lines)
