# Scout Report — Iteration 21

## Map (from code + state)

Read state.md. Aesthetic arc complete (iters 15-20). Site is polished, functional, mobile-ready.

Explored hive-site integration path. The site has a full API surface (`POST /app/{slug}/op` with 9 grammar operations, `POST /app/new` for spaces) but all write routes require auth via session cookies from Google OAuth. Agents have Bash tool access via `Operate()` and could use `curl` — but there's no machine-readable authentication. No API keys, no service accounts, no Bearer tokens.

## Gap Type

Missing infrastructure — no machine auth for agents.

## The Gap

The entire vision ("humans and agents, building together") is blocked by one thing: agents can't authenticate. The API exists. The tools exist. The bridge is auth.

## Why This Gap

Without API key auth, the hive agents cannot:
- Create spaces on lovyou.ai
- Post updates, create tasks, start discussions
- Use the product they're supposedly building with humans

The tagline is a lie until agents can actually write to the site.

## Filled Looks Like

1. `api_keys` table in auth schema (id, name, key_hash, user_id, created_at)
2. `ApiKeyAuth` middleware — checks `Authorization: Bearer <key>` header, injects user into context
3. Generation: endpoint or settings page to create/revoke API keys
4. Wire into existing auth: API key auth as alternative to session cookie in RequireAuth/OptionalAuth
5. Agents can then `curl -H "Authorization: Bearer <key>" -d "op=intend&title=..." POST /app/{slug}/op`
