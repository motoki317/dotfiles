package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestCommandCreatesMissingSkillLink(t *testing.T) {
	root := t.TempDir()
	claudeSkill := filepath.Join(root, ".claude", "skills", "alpha")
	agentSkills := filepath.Join(root, ".agents", "skills")

	mustMkdirAll(t, claudeSkill)
	mustMkdirAll(t, agentSkills)
	mustWriteFile(t, filepath.Join(claudeSkill, "SKILL.md"), "---\nname: alpha\n---\n")

	command := exec.Command("go", "run", ".", "--root", root)
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("sync command failed: %v\n%s", err, output)
	}

	link := filepath.Join(agentSkills, "alpha")
	target, err := os.Readlink(link)
	if err != nil {
		t.Fatalf("read created skill link: %v", err)
	}
	if want := "../../.claude/skills/alpha"; target != want {
		t.Fatalf("skill link target = %q, want %q", target, want)
	}
}

func TestCheckReportsMissingLinkWithoutCreatingIt(t *testing.T) {
	root := t.TempDir()
	claudeSkill := filepath.Join(root, ".claude", "skills", "alpha")
	mustMkdirAll(t, claudeSkill)
	mustWriteFile(t, filepath.Join(claudeSkill, "SKILL.md"), "")

	command := exec.Command("go", "run", ".", "--check", "--root", root)
	output, err := command.CombinedOutput()
	if exitError, ok := err.(*exec.ExitError); !ok || exitError.ExitCode() != 1 {
		t.Fatalf("check exit error = %v, want status 1\n%s", err, output)
	}
	if !strings.Contains(string(output), "alpha") {
		t.Fatalf("check output = %q, want skill name alpha", output)
	}
	if _, err := os.Lstat(filepath.Join(root, ".agents", "skills", "alpha")); !os.IsNotExist(err) {
		t.Fatalf("check created missing link: %v", err)
	}
}

func TestCommandPrunesStaleManagedLinksAndPreservesExternalSkills(t *testing.T) {
	root := t.TempDir()
	claudeSkills := filepath.Join(root, ".claude", "skills")
	agentSkills := filepath.Join(root, ".agents", "skills")
	mustMkdirAll(t, claudeSkills)
	mustMkdirAll(t, filepath.Join(agentSkills, "external"))
	mustWriteFile(t, filepath.Join(agentSkills, "external", "SKILL.md"), "")
	staleLink := filepath.Join(agentSkills, "removed")
	if err := os.Symlink("../../.claude/skills/removed", staleLink); err != nil {
		t.Fatalf("create stale link: %v", err)
	}

	command := exec.Command("go", "run", ".", "--root", root)
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("sync command failed: %v\n%s", err, output)
	}

	if _, err := os.Lstat(staleLink); !os.IsNotExist(err) {
		t.Fatalf("stale managed link still exists: %v", err)
	}
	if _, err := os.Stat(filepath.Join(agentSkills, "external", "SKILL.md")); err != nil {
		t.Fatalf("external skill was changed: %v", err)
	}
}

func TestCommandPrunesLinksToIneligibleClaudeSources(t *testing.T) {
	root := t.TempDir()
	claudeSkills := filepath.Join(root, ".claude", "skills")
	agentSkills := filepath.Join(root, ".agents", "skills")
	mustMkdirAll(t, filepath.Join(claudeSkills, "missing-skill"))
	mustMkdirAll(t, filepath.Join(agentSkills, "other"))
	mustWriteFile(t, filepath.Join(agentSkills, "other", "SKILL.md"), "")
	if err := os.Symlink("../../.agents/skills/other", filepath.Join(claudeSkills, "reowned")); err != nil {
		t.Fatalf("create reverse-owned skill link: %v", err)
	}
	for _, name := range []string{"missing-skill", "reowned"} {
		if err := os.Symlink("../../.claude/skills/"+name, filepath.Join(agentSkills, name)); err != nil {
			t.Fatalf("create managed link %s: %v", name, err)
		}
	}

	command := exec.Command("go", "run", ".", "--root", root)
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("sync command failed: %v\n%s", err, output)
	}

	for _, name := range []string{"missing-skill", "reowned"} {
		if _, err := os.Lstat(filepath.Join(agentSkills, name)); !os.IsNotExist(err) {
			t.Errorf("ineligible managed link %s still exists: %v", name, err)
		}
	}
	if _, err := os.Stat(filepath.Join(agentSkills, "other", "SKILL.md")); err != nil {
		t.Fatalf("source agent skill was changed: %v", err)
	}
}

func TestCommandExcludesClaudeOnlySkillsFromCodexExposure(t *testing.T) {
	root := t.TempDir()
	claudeSkills := filepath.Join(root, ".claude", "skills")
	agentSkills := filepath.Join(root, ".agents", "skills")
	for _, name := range []string{"codex-work", "alpha"} {
		mustMkdirAll(t, filepath.Join(claudeSkills, name))
		mustWriteFile(t, filepath.Join(claudeSkills, name, "SKILL.md"), "")
	}
	mustMkdirAll(t, agentSkills)
	if err := os.Symlink("../../.claude/skills/codex-work", filepath.Join(agentSkills, "codex-work")); err != nil {
		t.Fatalf("create excluded skill link: %v", err)
	}

	command := exec.Command("go", "run", ".", "--root", root)
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("sync command failed: %v\n%s", err, output)
	}

	if _, err := os.Lstat(filepath.Join(agentSkills, "codex-work")); !os.IsNotExist(err) {
		t.Errorf("excluded skill link still exists: %v", err)
	}
	if _, err := os.Readlink(filepath.Join(agentSkills, "alpha")); err != nil {
		t.Errorf("non-excluded skill link missing: %v", err)
	}
}

func TestCommandRejectsOwnershipCollisionBeforeApplyingChanges(t *testing.T) {
	root := t.TempDir()
	claudeSkills := filepath.Join(root, ".claude", "skills")
	agentSkills := filepath.Join(root, ".agents", "skills")
	for _, name := range []string{"alpha", "beta"} {
		mustMkdirAll(t, filepath.Join(claudeSkills, name))
		mustWriteFile(t, filepath.Join(claudeSkills, name, "SKILL.md"), "")
	}
	mustMkdirAll(t, filepath.Join(agentSkills, "alpha"))
	mustWriteFile(t, filepath.Join(agentSkills, "alpha", "SKILL.md"), "")

	command := exec.Command("go", "run", ".", "--root", root)
	output, err := command.CombinedOutput()
	if exitError, ok := err.(*exec.ExitError); !ok || exitError.ExitCode() != 1 {
		t.Fatalf("collision exit error = %v, want status 1\n%s", err, output)
	}
	if !strings.Contains(string(output), "alpha") {
		t.Fatalf("collision output = %q, want skill name alpha", output)
	}
	if _, err := os.Stat(filepath.Join(agentSkills, "alpha", "SKILL.md")); err != nil {
		t.Fatalf("colliding skill was changed: %v", err)
	}
	if _, err := os.Lstat(filepath.Join(agentSkills, "beta")); !os.IsNotExist(err) {
		t.Fatalf("command partially created beta link: %v", err)
	}
}

func mustMkdirAll(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatalf("create directory %s: %v", path, err)
	}
}

func mustWriteFile(t *testing.T, path, contents string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatalf("write file %s: %v", path, err)
	}
}
