# Scout Report — Iteration 47

## Gap: Handler layer untested + SQL injection in Mind

Two gaps, both correctness:

1. **Handler tests** — store is tested (iter 45) but no HTTP round-trip tests. The handlers are the API surface that users and agents interact with.

2. **SQL injection** — `findAgentParticipant` uses `"{"+strings.Join(tags,",")+"}` to build a Postgres array literal. Tags with commas or braces could break or exploit the query.
