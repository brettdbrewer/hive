---

## SCOUT REPORT — Iteration 295

**Gap:** No production observability—the hive ships autonomous features but has no visibility into errors, performance, or usage.

**Evidence:**

1. **State.md explicitly documents this** (line 101): "No observability: No error monitoring, no analytics, no usage tracking. Building into a void."

2. **Recent commits show PRMode + artifact infrastructure is solid** (commits 109e622, 15bc79c, f8e81fe):
   - Daemon resets to main between cycles ✓
   - Builder writes loop/build.md ✓  
   - Critic writes loop/critique.md ✓
   - Infrastructure for tracking what's built works

3. **But production has no error handling at site/hive level:**
   - No Sentry/error tracking integration
   - No analytics on what features users are adopting  
   - No dashboard showing API errors, failed builds, failed deploys
   - When the site or hive encounters a silent failure, nobody knows until Matt manually checks logs

4. **The PM role (pm.go) reads completed tasks but has no context on whether those tasks WORKED in production:**
   - Builder ships code, deploys successfully, but if the feature breaks on production, PM doesn't know
   - Loop can only improve what it can measure (lesson 36: "The loop can only catch errors it has checks for")

5. **Backlog explicitly lists this under "URGENT"** (line 101): "Building into a void."

**Impact:**

- **Reliability risk:** Silent failures in production remain silent. The hive's credibility with Lovatts depends on knowing whether shipped features work.
- **Feedback loop broken:** The PM uses completed tasks + git log to decide next priorities. If a shipped feature breaks, PM doesn't know, and keeps building on a cracked foundation.
- **Debugging impossible:** When the site/hive has an issue, diagnosis requires manual log inspection. No structured error context.
- **Adoption invisible:** Can't measure whether features are being used. No data on which product layers matter to users.

**Scope:**

Hive repo + site repo:
- `pkg/runner/observer.go` (exists but minimal—only reads git log + board state)
- Site repo: error boundaries, error handler middleware
- Site repo: analytics event schema and reporting
- Both repos: structured logging (JSON logs with stack traces)
- New infrastructure: error tracking service integration (Sentry/Rollbar or custom)

**Suggestion:**

Build a **production observability layer** with two parts:

**Part 1 (hive repo): Pipeline observability**
- Observer role logs structured data: build start time, end time, cost, commit hash, deploy target, deploy status
- Add a `pipeline_events` table to store: tick number, phase, status (SUCCESS/FAILURE), error message (if any), cost
- Build a `/hive/events` endpoint showing last 100 pipeline events (JSON)
- Extend `/hive` dashboard to show: pipeline health (% success), recent errors, cost trend

**Part 2 (site repo): Production observability**
- Add error boundary in Go handlers + templ views — catch panics, log structured errors
- Integrate Sentry (or self-hosted error tracking): every error posts `{"service": "site", "error": {...}, "timestamp": "...", "user_id": "...", "commit": "..."}`
- Add basic analytics: every grammar op emits `{"op": "intend", "kind": "task", "user_id": "...", "space_id": "...", "timestamp": "..."}`
- Build a `/status` page showing: site uptime, last 10 errors, active users today, top 5 ops by count
- Observer role can query `/status` to see site health in its report

**Priority:** HIGHEST. This is prerequisite for the Lovatts engagement. "Company in a box" requires proof that the hive's code actually works in production. Right now, the hive ships features but can't prove they're working.

---

**Ready for Architect phase.**