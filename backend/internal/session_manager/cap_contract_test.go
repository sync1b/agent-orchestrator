package sessionmanager

import (
	"encoding/json"
	"testing"

	"github.com/aoagents/agent-orchestrator/backend/internal/domain"
	"github.com/aoagents/agent-orchestrator/backend/internal/ports"
)

// Keep this test on public project configuration and Spawn so the same test
// compiles against the v0.13.1 base, which did not expose the cap internals.
func TestProjectSessionCeilingRefusesBeforeCreation(t *testing.T) {
	m, st, _, _ := newManager()
	project := st.projects["mer"]
	if err := json.Unmarshal([]byte(`{"maxConcurrentSessions":1}`), &project.Config); err != nil {
		t.Fatal(err)
	}
	st.projects["mer"] = project
	spawn := ports.SpawnConfig{ProjectID: "mer", Kind: domain.KindWorker}
	if _, _, _, err := m.Spawn(ctx, spawn); err != nil {
		t.Fatalf("first spawn: %v", err)
	}
	if _, _, _, err := m.Spawn(ctx, spawn); err == nil {
		t.Fatal("second spawn succeeded above project session ceiling")
	}
	if got := len(st.sessions); got != 1 {
		t.Fatalf("refused spawn created a session: got %d, want 1", got)
	}
}
