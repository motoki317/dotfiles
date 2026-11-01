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

type usageSource struct {
	file string
	bin  string
	arg  string
	// Keep interval below statusLine.refreshInterval (60 s in settings.json), so that
	// each idle rerun finds the cache stale and refreshes it.
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

// The lock file's mtime records every refresh attempt, so a failing fetch retries
// once per interval instead of on every render. It does not exclude a concurrent
// render, which at worst starts one extra fetch.
// Setsid moves the refresh out of the render's session and process group, so a
// signal that cancels the render does not reach it.
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

func (s usageSource) refresh() {
	if windows := s.fetch(); len(windows) > 0 {
		_ = s.write(&cacheFile{FetchedAt: time.Now().Unix(), Windows: windows})
	}
}
