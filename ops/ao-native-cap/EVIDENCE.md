# Native AO session ceiling evidence

Code head: `e89e247741c8387617a1ee24708ffab9d14e62c4`.
Base: upstream `v0.13.1` (`460d9c45f7bf5503808254e0cb02ee52bf537349`).
Toolchain: Go 1.27.1. Upstream source: PR #5649, admission commit
`34dfe4ee9`; no unverified external library API was introduced.

## Discriminating base failure

The same `cap_contract_test.go` was copied into a detached base worktree:

```text
go test ./internal/session_manager -run '^TestProjectSessionCeilingRefusesBeforeCreation$' -count=1
--- FAIL: TestProjectSessionCeilingRefusesBeforeCreation
    second spawn succeeded above project session ceiling
FAIL
```

The test configures a one-session project ceiling through public JSON,
starts one worker, and asserts a second worker is refused with no second
session row. It compiles on both the base and patched code.

## Patched head

```text
go test ./internal/session_manager -run 'Test(ProjectSessionCeilingRefusesBeforeCreation|SpawnRefusesWorkerAtGlobalConcurrencyCap|SpawnOrchestratorExemptFromConcurrencyCap|SpawnEnforcesProjectConcurrencyCap|WorkerAdmissionCountsStartsNotYetDurable|WorkerAdmissionAppliesProjectCapToRestoreStarts|WorkerAdmissionCountsPersistedSessionsAfterManagerRestart)' -count=1
ok  github.com/aoagents/agent-orchestrator/backend/internal/session_manager

go test -race ./internal/session_manager -run 'Test(ProjectSessionCeilingRefusesBeforeCreation|SpawnRefusesWorkerAtGlobalConcurrencyCap|SpawnOrchestratorExemptFromConcurrencyCap|SpawnEnforcesProjectConcurrencyCap|WorkerAdmissionCountsStartsNotYetDurable|WorkerAdmissionAppliesProjectCapToRestoreStarts|WorkerAdmissionCountsPersistedSessionsAfterManagerRestart)' -count=1
ok  github.com/aoagents/agent-orchestrator/backend/internal/session_manager

go test ./internal/config ./internal/cli ./internal/daemon ./internal/service/session ./internal/session_manager -count=1
ok  5 packages
```

The seven focused tests cover below-limit spawn, refusal, orchestrator
occupancy, per-project limits, in-flight reservation, restore admission,
and restart counting from persisted rows. The race run passed. No live AO
daemon was changed; rollout/build pinning remains a coordinator action after
review and the required gates.
