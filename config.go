package main

import (
	"io"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/koron/upvim/internal/arch"
)

type config struct {
	name      string
	cpu       arch.CPU
	targetDir string
	dataDir   string
	logDir    string
	tmpDir    string
	varDir    string
}

func newConfig(dir string) (*config, error) {
	exe := filepath.Join(dir, "vim.exe")
	cpu, err := arch.Exe(exe)
	if err != nil {
		return nil, err
	}
	var name string
	dataDir := filepath.Join(dir, "_upvim")
	switch cpu {
	case arch.X86:
		name = "vim74-x86"
	case arch.AMD64:
		name = "vim74-x64"
	}
	return &config{
		name:      name,
		cpu:       cpu,
		targetDir: dir,
		dataDir:   dataDir,
		logDir:    filepath.Join(dataDir, "log"),
		tmpDir:    filepath.Join(dataDir, "tmp"),
		varDir:    filepath.Join(dataDir, "var"),
	}, nil
}

func (c *config) downloadPath(targetURL string) (string, error) {
	u, err := url.Parse(targetURL)
	if err != nil {
		return "", err
	}
	return filepath.Join(c.tmpDir, filepath.Base(u.Path)), nil
}

func (c *config) recipePath() string {
	return filepath.Join(c.varDir, c.name+"-recipe.txt")
}

func (c *config) anchorPath() string {
	return filepath.Join(c.varDir, c.name+"-anchor.txt")
}

func (c *config) anchor() (time.Time, error) {
	f, err := os.Open(c.anchorPath())
	if err != nil {
		return time.Time{}, nil
	}
	defer f.Close()
	buf := make([]byte, 25)
	if _, err := io.ReadFull(f, buf); err != nil {
		return time.Time{}, nil
	}
	return time.Parse(time.RFC3339, string(buf))
}

func (c *config) updateAnchor(t time.Time) error {
	f, err := os.Create(c.anchorPath())
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := io.WriteString(f, t.Format(time.RFC3339)); err != nil {
		return err
	}
	return f.Sync()
}

func (c *config) dirs() []string {
	return []string{
		c.targetDir,
		c.dataDir,
		c.logDir,
		c.tmpDir,
		c.varDir,
	}
}

func (c *config) prepare() error {
	for _, dir := range c.dirs() {
		if err := os.MkdirAll(dir, 0777); err != nil {
			return err
		}
	}
	return nil
}
