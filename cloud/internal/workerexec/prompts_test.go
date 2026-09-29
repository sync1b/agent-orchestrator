package workerexec

import (
	"strings"
	"testing"
)

func TestWorkerPromptAllowsHelpersWithoutChangingAttemptOwnership(t *testing.T) {
	for _, hasOrchestrator := range []bool{false, true} {
		prompt := workerSystemPrompt("/skills/using-ao", hasOrchestrator)
		for _, want := range []string{
			"You may use available native subagent or task-delegation tools",
			"bounded subtasks of the assigned task or issue attempt",
			"This AO session remains the sole execution owner",
			"Helpers inherit the assigned scope, workspace, permissions, and publishing restrictions",
			"grant no additional tools, credentials, model access, or spend authority",
			"Parallel writers must own disjoint files; serialize overlapping edits",
			"never the shared project checkout or another session's workspace",
			"Helpers must not independently admit queued work, create work units or AO sessions",
			"No separate issue or AO worker is needed for an in-session helper",
			"Collect helper results or stop the helpers before declaring the task complete",
			"do not replace a required independent AO reviewer",
			"Neither the author nor its helpers may approve the author's work",
			"Helpers must not satisfy independent review or merge gates",
			"$AO_SESSION_BRANCH",
			"Do not create a different branch",
			"Only when the task asked you to review a PR",
			"$AO_REVIEW_HELP",
		} {
			if !strings.Contains(prompt, want) {
				t.Errorf("hasOrchestrator=%t: worker prompt missing %q", hasOrchestrator, want)
			}
		}
		if strings.Count(prompt, "## Native helper policy v1") != 1 {
			t.Errorf("hasOrchestrator=%t: expected exactly one native-helper policy", hasOrchestrator)
		}
		for _, forbidden := range []string{
			"Do not use your runtime's built-in subagent",
			"ao session",
			"ao status",
		} {
			if strings.Contains(prompt, forbidden) {
				t.Errorf("hasOrchestrator=%t: worker prompt includes obsolete or desktop-only instruction %q", hasOrchestrator, forbidden)
			}
		}
		if strings.Contains(prompt, "ao report") != hasOrchestrator {
			t.Errorf("hasOrchestrator=%t: report command does not match parent scope", hasOrchestrator)
		}
	}
}

func TestOrchestratorPromptAllowsDirectWorkOnlyWithAssignedIsolation(t *testing.T) {
	prompt := strings.Join(strings.Fields(orchestratorSystemPrompt("/skills/using-ao")), " ")
	for _, want := range []string{
		"coordination-only by default",
		"The shared project checkout is read-only",
		"an isolated sandbox alone does not assign you an implementation task or its ownership",
		"small, bounded implementation or verification directly only when this session has been explicitly assigned",
		"an isolated workspace and is its sole execution owner",
		"Confirm that workspace and ownership before any write",
		"use the assigned non-default branch",
		"If those conditions are absent or uncertain, spawn or redirect an AO worker",
		"Honor an existing explicit assignment without asking for the same permission again",
		"Keep PR ownership on AO workers",
		"explicit ownership handoff",
		"do not spawn normal queued work around it",
		"Helper writes must meet the same isolated-workspace and ownership conditions",
		"do not replace a required independent AO reviewer",
		"Neither the author nor its helpers may approve the author's work",
		"Helpers must not satisfy independent review or merge gates",
		"ao list",
		"ao spawn --name",
	} {
		if !strings.Contains(prompt, want) {
			t.Errorf("orchestrator prompt missing %q", want)
		}
	}
	for _, forbidden := range []string{
		"Never edit source files",
		"ask for confirmation once",
		"AO workers only",
		"ao session",
		"ao status",
		"--project",
	} {
		if strings.Contains(prompt, forbidden) {
			t.Errorf("orchestrator prompt includes obsolete or desktop-only instruction %q", forbidden)
		}
	}
}
