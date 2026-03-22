# Scout Report — Iteration 50

## Gap: Tags display raw hex IDs instead of names

Iter 49 changed tags to store user IDs (correct). But the templates render `tag` directly — showing `36509418df854dd4a709cfee3e915a17` instead of "hive". Also, existing data (30 nodes, 30 ops) had empty `author_id`/`actor_id` columns.

## What "Filled" Looks Like

Conversation list and detail views show participant display names (resolved from IDs). Existing data has `author_id`/`actor_id` populated.
