# Build Report — Iteration 110

Space invites. New `invites` table (token, space_id, created_by). Generate invite link from Settings (owner-only). `/join/{token}` route accepts invite — joins user to space and redirects. Reuses existing invite if one exists. Invite URL shown as copyable field on Settings.

11 tables total. First growth feature — spaces can now be shared.
