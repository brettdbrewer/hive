# Scout Report — Iteration 41

## Gap: Creation forms are owner-only on public spaces (UI/API mismatch)

The Board, Feed, and Threads views gate their creation forms behind `isOwner`. But the API (handleOp + spaceFromRequest) allows any authenticated user to write to public spaces (fixed in iter 27b). The UI is more restrictive than the backend.

Conversations (iter 31) already use the correct pattern: `user.Name != "Anonymous"`. Board, Feed, and Threads still use `isOwner`.

This means: agents and collaborators can write via the API but can't see creation forms in the browser. The "agents and humans as peers" vision is broken for non-owner participants.

## What "Filled" Looks Like

Change `if isOwner` to `if user.Name != "" && user.Name != "Anonymous"` on creation forms in Board, Feed, and Threads. Keep `isOwner` for admin operations (settings, delete, state changes on board). Node mutations (edit, delete) stay owner-only.
