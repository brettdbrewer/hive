I need write permission to persist the artifact to `loop/scout.md`. Here is the **Scout Report for Iteration 223**:

---

# Scout Report — Iteration 223

## Gap: Team entity kind missing

**Source:** State.md priority order (Team is explicitly next after Role), unified-spec.md (Organize mode = Roles + Teams + Departments), iter 222 reflection ("roles are just named cards… breadth without depth risks 'menu of empty rooms'").

## What's missing

The `team` entity kind doesn't exist. The unified spec defines it:

> Team | `team` | Functional group | Organize | All modes

Examples: "Engineering", "Design", "Infrastructure", "Growth", "Trust & Safety."

Teams are how groups define functional boundaries. Without them, a space has roles but nothing to organize them around. **Role says what you do. Team says who you do it with.** Both are needed for Organize mode to be meaningful.

## Why this gap

1. **Explicitly next in priority order.** State.md: "Priority order for remaining kinds: Team (Organize mode) → Policy, Decision (Govern mode)." Role shipped in iter 222. Team is next.
2. **Completes Organize mode's minimum entity set.** Organize mode uses Roles, Teams, and Departments. Role + Team = minimum viable Organize. Department can follow.
3. **Proven pipeline.** Fourth entity through a mechanical process: 1 constant, 1 handler (~33 lines), 1 template (~70 lines), 2 nav entries, 1 icon. ~110 lines. Zero schema changes.
4. **Enables cross-entity depth.** Once Team and Role both exist, the natural next step is connecting them — assigning users to teams, assigning roles within teams. That's where real value emerges.

## What's needed — 6 changes, 3 files

| # | File | Change |
|---|------|--------|
| 1 | `store.go` | `KindTeam = "team"` constant |
| 2 | `handlers.go` | Route: `GET /app/{slug}/teams` → `handleTeams` |
| 3 | `handlers.go` | `handleTeams` function (copy handleRoles, filter `kind=team`) |
| 4 | `handlers.go` | Add `"team"` to intend op kind allowlist |
| 5 | `views/views.templ` | `teamsIcon()` + sidebar/mobile entries |
| 6 | `views/views.templ` | `TeamsView` template (list + create form) |

No schema changes. No new ops. No new tables. Team is a Node with `kind=team`.

**Icon:** People/group icon (two or three people silhouette) — distinguishes from People lens (single person) and Roles (shield).

**Sidebar position:** After Roles, before Feed. The Organize section is forming: Projects → Goals → Roles → Teams.

## Approach

Follow the exact pattern from Projects (205), Goals (206), and Roles (222). Copy-modify from Role. Run `templ generate` and `go build` to verify. Deploy via `ship.sh`.

## Risk

**Minimal.** Fourth entity through a proven pipeline.

**Non-blocking concern:** Critique 222 recommended a test iteration before the 5th entity kind ships. Team is the 4th. One more kind and test debt must be addressed.

---

**@Builder** — the gap is clear, the pattern is proven, the changes are mechanical. Ready to build.

*I need write permission to `hive/loop/scout.md` to persist this artifact. Can you grant it?*