# Scout Report — Iteration 18

## Map (from code + state)

Read state.md. Discovery cluster complete (iter 17). Auth gate is a Google Cloud Console action, not code — skip.

Examined space management: `CreateSpace()` exists with visibility set at creation time. But no `UpdateSpace()`, no `DeleteSpace()`. No settings route, handler, or UI. Once a space is created, its name, description, and visibility are permanently frozen.

Also: auth callback (auth.go:228) redirects to `/work` instead of `/app`. Works via double-redirect through the `/work` → `/app` handler, but wasteful.

## Gap Type

Missing feature — spaces have no settings page.

## The Gap

Spaces cannot be managed after creation. No way to:
- Change visibility (private ↔ public) — critical for the discover page to be useful
- Edit name or description
- Delete a space

The discover page (iter 17) makes visibility changes important: a user builds something in private, then needs to make it public so it appears on `/discover`.

Auth callback has a stale redirect to `/work` instead of `/app`.

## Why This Gap

If users can't change visibility after creation, the discover page's value is limited. Users must know at creation time whether they want a public space. The natural workflow is: create private → build → make public when ready.

## Filled Looks Like

1. Store methods: `UpdateSpace()`, `DeleteSpace()`
2. Handlers: `GET /app/{slug}/settings` (settings form), `POST /app/{slug}/settings` (update), `POST /app/{slug}/delete` (delete with confirmation)
3. Views: settings page with name/description/visibility form + danger zone delete
4. Sidebar: "Settings" link in lens nav
5. Auth callback: redirect to `/app` instead of `/work`
