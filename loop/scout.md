# Scout Report — Iteration 29

## Map

User reported: sidebar and main content scroll together instead of independently. The appLayout uses `min-h-screen` on the body, which lets it grow beyond the viewport. The `overflow-hidden` on the flex content div doesn't clip because its parent has no height constraint.

## Gap Type

Bug (user-reported)

## The Gap

Sidebar is not sticky. When you scroll the main content, the sidebar scrolls with it. This breaks the two-column app layout pattern — the sidebar should remain visible while content scrolls.

## Why This Gap Over Others

User flagged it directly: "the content and sidebar scroll together which is far from ideal." User-reported issues are highest priority per lesson 12.

## What "Filled" Looks Like

Sidebar stays fixed in place while main content scrolls independently. Each has its own scroll context. The app feels like a proper app shell, not a long document.
