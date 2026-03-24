# Build Report — Iteration 211

## Product Map

**Output:** `hive/loop/product-map.md`

**Key findings:**

1. **~56 distinct products** across 13 layer families. Each does one thing well. All share infrastructure.

2. **14 shared infrastructure components** that every product uses: auth, identity, DMs, notifications, search, reactions, endorsements, follows, command palette, activity feed, file attachments, @mentions, keyboard shortcuts, markdown.

3. **A product is a space configuration** — which entity kinds are active, which modes are visible, what the focused view looks like. NOT a separate codebase.

4. **Navigation model:** 13-layer main menu → product family → specific product. Users drill into what they need.

5. **Cross-cutting features (DMs, search, endorsements) are built once.** Every product gets them. This is the platform advantage — each product starts 60% done.

6. **Build strategy:** First product per layer → second product per layer → open platform for hives to build more.

This iteration produced spec, not code.
