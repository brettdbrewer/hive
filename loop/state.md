# Loop State

Living document. Updated by the Reflector each iteration. Read by the Scout first.

Last updated: Iteration 7, 2026-03-22.

## Current System State

Five repos, all compiling and tested:
- **eventgraph** — foundation. Postgres stores, 201 primitives, trust, authority. Complete.
- **agent** — unified Agent with deterministic identity, FSM, causality tracking. Complete.
- **work** — task store for hive agent coordination. Complete.
- **hive** — 4 agents (Strategist, Planner, Implementer, Guardian), agentic loop, budget. Complete.
- **site** — lovyou.ai on Fly.io. **Deployed and live.** Complete product:
  - Blog (43 posts, markdown → HTML)
  - Reference (cognitive grammar, graph grammar, 13 layer grammars, 201 primitives, 28 agent primitives)
  - Auth (Google OAuth — **test mode, only Matt can access**)
  - Unified graph product (3 tables: spaces/nodes/ops, 10 grammar operations, 5 lenses, HTMX, full CRUD)
  - Landing page: concrete product description with five lens cards and three-step flow
  - SEO: meta description + OG tags + Twitter cards on every page

Entry point: `cmd/hive` (CLI only). No daemon, no web dashboard for the hive itself.

Deploy method: `fly deploy --remote-only` (avoids Docker Desktop dependency).

**Fly/Neon resources can be scaled up** per user authorization.

## Lessons Learned

1. **Code is truth, not docs.** Always read code to assess current state.
2. **Accumulate knowledge.** state.md + reflections.md prevent repeated mistakes.
3. **Update state.md every iteration.** Prevents phantom gaps.
4. **Ship what you build.** Every iteration should end with changes in version control.
5. **Try alternatives before declaring blockers.** `--remote-only` worked; `--local-only` was the problem.
6. **Every Build iteration should deploy.** Scout → Builder → commit → push → deploy in one cycle.
7. **Focus on public-facing improvements.** OAuth is in test mode — build for the visitor who can't log in yet.

## Known Issues

- DATABASE_URL may not be set on Fly — visitors clicking "Open the app" could get a 503.
- No sitemap.xml or robots.txt.
- No analytics.
- AUDIT.md in hive/docs/ is stale (March 9th). Low priority.
- Boot events re-emitted every agent run (chatty but harmless).

## What the Scout Should Focus On Next

The loop is in Build mode. The site now has a clear landing page and proper SEO. Next candidates:

1. **Wire Neon DB to Fly** — if DATABASE_URL isn't set, the app returns 503. Check if it's configured. If not, set it up so the product is accessible to visitors.
2. **Sitemap.xml** — search engines need a sitemap to discover 250+ pages efficiently.
3. **App onboarding** — what does a first-time user experience? Is the anonymous fallback useful?
4. **Hive autonomy** — make the core loop self-running.

Pick one. Build it. Deploy it.
