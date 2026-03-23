# Scout Report — Iteration 91

## Gap: No search — the platform is discoverable but not searchable

The auth gate is open. Real users can sign up. But there's no way to search for anything. The discover page lists spaces. The market page searches tasks. There's no unified search across spaces, content, or users.

A user who heard about lovyou.ai and wants to find a specific space, post, or person has to browse manually. This is the biggest usability gap for a public platform.

## What "Filled" Looks Like

`/search?q=term` — a unified search page that finds:
- **Spaces** matching name or description
- **Content** (tasks, posts, threads) matching title or body from public spaces
- **Users** matching name

Results grouped by type. Linked to the relevant pages. Search box in the header for easy access.

## Approach

1. New store query: `Search(ctx, query, limit)` — searches spaces, nodes, and users
2. New template in views/search.templ
3. New route `GET /search` in main.go
4. Add search box to the header/nav
