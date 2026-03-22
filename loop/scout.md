# Scout Report — Iteration 40

## Gap: Logged-in users see marketing page instead of their workspace

The home page (`/`) shows the same landing content to everyone — "Humans and agents, building together" hero, lens cards, how-it-works steps. Logged-in users must click "Create a space" to reach `/app` where their actual work lives. This is friction on the return visit.

## Why This Gap

The home page was built (iter 15) for first visitors. Auth was added later (iter 21+). Nobody wired them together. The logged-in experience starts at `/app`, but the URL bar says `/`.

## What "Filled" Looks Like

If the user has a session cookie (logged in), redirect `/` to `/app`. If not, show the current landing page. One conditional in the home handler. Zero new templates.
