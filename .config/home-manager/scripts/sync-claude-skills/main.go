package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const managedTargetPrefix = "../../.claude/skills/"

type linkChange struct {
	path   string
	target string
}

type syncPlan struct {
	missing    []linkChange
	stale      []string
	collisions []string
}

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	flags := flag.NewFlagSet("sync-claude-skills", flag.ContinueOnError)
	root := flags.String("root", "", "home directory containing .claude and .agents")
	check := flags.Bool("check", false, "report drift without changing files")
	if err := flags.Parse(args); err != nil || flags.NArg() != 0 {
		return 2
	}

	if *root == "" {
		var err error
		*root, err = os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
	}

	plan, err := buildPlan(*root)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	for _, path := range plan.collisions {
		fmt.Fprintf(os.Stderr, "ownership collision: %s\n", path)
	}
	if len(plan.collisions) > 0 {
		return 1
	}

	if *check {
		for _, change := range plan.missing {
			fmt.Fprintf(os.Stderr, "missing link: %s -> %s\n", change.path, change.target)
		}
		for _, path := range plan.stale {
			fmt.Fprintf(os.Stderr, "stale link: %s\n", path)
		}
		if len(plan.missing) > 0 || len(plan.stale) > 0 {
			return 1
		}
		return 0
	}

	if err := applyPlan(*root, plan); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}

func buildPlan(root string) (syncPlan, error) {
	claudeSkills := filepath.Join(root, ".claude", "skills")
	agentSkills := filepath.Join(root, ".agents", "skills")
	entries, err := os.ReadDir(claudeSkills)
	if err != nil {
		return syncPlan{}, err
	}

	desired := make(map[string]linkChange)
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		skillInfo, err := os.Stat(filepath.Join(claudeSkills, entry.Name(), "SKILL.md"))
		if err != nil || skillInfo.IsDir() {
			continue
		}
		desired[entry.Name()] = linkChange{
			path:   filepath.Join(agentSkills, entry.Name()),
			target: managedTargetPrefix + entry.Name(),
		}
	}

	plan := syncPlan{}
	for _, change := range desired {
		info, err := os.Lstat(change.path)
		if os.IsNotExist(err) {
			plan.missing = append(plan.missing, change)
			continue
		}
		if err != nil {
			return syncPlan{}, err
		}
		if info.Mode()&os.ModeSymlink == 0 {
			plan.collisions = append(plan.collisions, change.path)
			continue
		}
		target, err := os.Readlink(change.path)
		if err != nil {
			return syncPlan{}, err
		}
		if target != change.target {
			plan.collisions = append(plan.collisions, change.path)
		}
	}

	agentEntries, err := os.ReadDir(agentSkills)
	if err != nil && !os.IsNotExist(err) {
		return syncPlan{}, err
	}
	for _, entry := range agentEntries {
		if _, ok := desired[entry.Name()]; ok {
			continue
		}
		path := filepath.Join(agentSkills, entry.Name())
		target, err := os.Readlink(path)
		if err == nil && strings.HasPrefix(target, managedTargetPrefix) {
			plan.stale = append(plan.stale, path)
		}
	}

	return plan, nil
}

func applyPlan(root string, plan syncPlan) error {
	agentSkills := filepath.Join(root, ".agents", "skills")
	if err := os.MkdirAll(agentSkills, 0o755); err != nil {
		return err
	}
	for _, path := range plan.stale {
		if err := os.Remove(path); err != nil {
			return err
		}
	}
	for _, change := range plan.missing {
		if err := os.Symlink(change.target, change.path); err != nil {
			return err
		}
	}
	return nil
}
