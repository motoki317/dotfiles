package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
)

type cacheFile struct {
	FetchedAt int64    `json:"fetched_at"`
	Windows   []Window `json:"windows"`
}

// usageSource describes one background-refreshed, account-global rate-limit source
// (Claude or Codex). Only these four fields differ between sources — every cache
// mechanic below is shared. bin is the presence gate: a machine without that tool
// never spawns a refresh. arg is the detached --refresh-* self-invocation. interval
// is seconds and doubles as both the cache-staleness threshold and, via the lock
// file's mtime, the between-attempts backoff, so a logged-out or failing fetch cannot
// retry on every render. Reads of these account-status endpoints are unmetered, so a
// ~per-minute refresh is safe; pair it with statusLine.refreshInterval in settings.
type usageSource struct {
	file     string
	bin      string
	arg      string
	interval int64
	fetch    func() []Window
}

func (s usageSource) current(now int64) []Window {
	cache := s.read()
	if now-cache.FetchedAt > s.interval {
		s.spawn(now)
	}
	return cache.Windows
}

// cachePath/lockPath locate the cache and its refresh-backoff lock under the user
// cache dir (~/Library/Caches on macOS, $XDG_CACHE_HOME on Linux). Empty means the
// cache dir is unavailable, which disables the cached path for that source.
func (s usageSource) cachePath() string {
	dir, err := os.UserCacheDir()
	if err != nil {
		return ""
	}
	return filepath.Join(dir, "claude-statusline", s.file)
}

func (s usageSource) lockPath() string {
	if p := s.cachePath(); p != "" {
		return p + ".lock"
	}
	return ""
}

// Cache problems act as absence so rendering stays best-effort.
func (s usageSource) read() cacheFile {
	p := s.cachePath()
	if p == "" {
		return cacheFile{}
	}
	data, err := os.ReadFile(p)
	if err != nil {
		return cacheFile{}
	}
	var c cacheFile
	if json.Unmarshal(data, &c) != nil {
		return cacheFile{}
	}
	return c
}

// write atomically replaces the cache (temp file + rename) so a concurrent render
// reads either the prior complete snapshot or the new one, never a half-written file.
func (s usageSource) write(c *cacheFile) error {
	p := s.cachePath()
	if p == "" {
		return fmt.Errorf("no user cache dir")
	}
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	tmp := p + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, p)
}

// spawn starts a detached `claude-statusline <arg>` refresher unless the source's
// tool is absent or an attempt ran within interval. The lock mtime records attempts
// (success or failure alike), so a failing fetch backs off instead of retrying on
// every render; the loose stat/write race spawns at most one harmless extra fetch.
// Setsid detaches the child so it outlives this short-lived render.
func (s usageSource) spawn(now int64) {
	lock := s.lockPath()
	if lock == "" {
		return
	}
	if st, err := os.Stat(lock); err == nil && now-st.ModTime().Unix() < s.interval {
		return
	}
	if _, err := exec.LookPath(s.bin); err != nil {
		return
	}
	self, err := os.Executable()
	if err != nil {
		return
	}
	if os.MkdirAll(filepath.Dir(lock), 0o755) != nil {
		return
	}
	if os.WriteFile(lock, []byte(strconv.FormatInt(now, 10)), 0o644) != nil {
		return
	}
	cmd := exec.Command(self, s.arg)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	if cmd.Start() == nil {
		_ = cmd.Process.Release()
	}
}

// refresh runs in the detached child: it fetches the live snapshot and, on success,
// writes it to the cache. A failed fetch leaves the previous cache untouched.
func (s usageSource) refresh() {
	if windows := s.fetch(); len(windows) > 0 {
		_ = s.write(&cacheFile{FetchedAt: time.Now().Unix(), Windows: windows})
	}
}
