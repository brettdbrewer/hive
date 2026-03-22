# Build Report — Iteration 16

## What I planned

Apply a dark theme with warm rose accent across all site templates — layout, home, blog, reference, and graph/app views. Research-informed aesthetic: "Ember Minimalism."

## What I built

Changes across 10 files (5 templates + 5 generated) in the site repo.

### Theme system (layout.templ @theme block)
New Tailwind v4 custom colors replacing the old indigo-on-gray:
- `--color-brand: #e8a0b8` (rose, was #6366f1 indigo)
- `--color-void: #09090b` (page background, warm near-black)
- `--color-surface: #111113` (cards/panels)
- `--color-elevated: #18181b` (hover states, inputs)
- `--color-edge: #1e1e22` (borders)
- `--color-warm: #f0ede8` (primary text, warm off-white)
- `--color-warm-secondary: #c8c4bc` (body text)
- `--color-warm-muted: #78756e` (captions)
- `--color-warm-faint: #4a4844` (disabled text)

### Templates updated
1. **layout.templ** — theme colors, dark prose styles, dark header/footer
2. **home.templ** — light heading weight (300), dark cards, dark CTA buttons
3. **blog.templ** — dark blog cards, dark arc nav pills, dark post pages
4. **reference.templ** — dark reference cards, dark spec tables, dark primitive pages
5. **graph/views.templ** — shared themeBlock(), dark appLayout, dark sidebar, dark forms/inputs, dark badge colors (bg-indigo-500/15 instead of bg-indigo-100), updated helper functions (stateColorHex, priorityDotHex, stateBgClass, kindBadgeClass, opBadgeClass)

### Key design decisions
- Button text: `text-void` (dark) on rose background, not white
- Badges: semi-transparent backgrounds (`bg-indigo-500/15 text-indigo-400`) instead of light pastels
- Headings: `font-light` (300) at large sizes for sacred minimalism feel
- Prose: hardcoded hex colors match the theme (blockquote border = rose, code bg = elevated)
- Ring offset: `ring-offset-void` for selected state buttons in dark context

## Verification

- `templ generate` — success
- `go build -o /tmp/site.exe ./cmd/site/` — success
- Committed and pushed to main
- Deployed to Fly.io — both machines healthy
