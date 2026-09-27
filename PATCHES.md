# Local AO patches

## Session admission ceiling (sync1b/dev-bot-local#263)

- Base: upstream `v0.13.1` (`460d9c45f7bf5503808254e0cb02ee52bf537349`).
- Source: the admission-only change from upstream PR
  [#5649](https://github.com/Untrivial-ai/agent-orchestrator/pull/5649),
  commit `34dfe4ee9`, adapted to this base. That upstream PR is open as of
  2026-09-27; no released AO build contains the cap.
- Scope: global `AO_MAX_CONCURRENT_SESSIONS` (default 4, `0` opts out) and
  per-project `maxConcurrentSessions` (`0` inherits the global ceiling).
  Each live orchestrator occupies a slot, but an orchestrator may start at
  the ceiling to preserve recovery. Workers are refused before their seed row,
  worktree, or runtime is created. In-flight starts reserve capacity under a
  mutex; restarts count the persisted non-terminated session rows.
- Removal: once upstream merges and releases equivalent admission behavior,
  compare the released code and tests against this patch, upgrade the pin,
  run the focused tests in `ops/ao-native-cap/EVIDENCE.md`, and drop this
  patch. Do not carry both implementations.

This patch does not add scheduling or another session store. A configured
ceiling is a machine safety limit, while the orchestrator rules remain a soft
limit on planned fan-out.
