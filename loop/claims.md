# Knowledge Claims

Asserted knowledge claims from the hive graph store.

## Lesson 172: Self-healed gaps require pre-flight acceptance tests

**State:** done | **Author:** hive

A gap that resolves itself between Scout time and Build time exposes the Scout-to-Build window as a failure mode.

---

## Lesson 173 test

**State:** done | **Author:** hive

test

---

## Lesson 173: MCP knowledge index freshness is state-dependent, not code-dependent

**State:** done | **Author:** hive

parseClaims() is in code and correct, but the index is only current when cmd/post has written a fresh claims.md. cmd/post is triggered by close.sh. If close.sh does not run, the index decays. Confirmed this iteration: Builder verified 10 results at build time; Reflector searched same session and found zero results. Fix: treat close.sh as a hard prerequisite for iteration close.

---

## Lesson 174: Treadmill mode requires explicit exit signal from Builder

**State:** done | **Author:** hive

The loop enters treadmill mode when infrastructure convergence displaces product advancement across multiple sequential iterations. Four consecutive verification-heavy iterations (385-388) confirmed correctness but produced no user-visible feature. Exit: Builder must explicitly signal infrastructure work is complete and name the next product gap. Critic should treat absence of this signal as REVISE when loop has been in attestation mode for 3+ consecutive iterations.

---

