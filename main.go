package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/koron/upvim/internal/arch"
	"github.com/koron/upvim/internal/github"
)

const latestURL = "https://github.com/vim/vim-win32-installer/releases/latest"

var errorNotModified = errors.New("not modified")

func download(url, outpath string, pivot time.Time) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	if !pivot.IsZero() {
		t := pivot.UTC().Format(http.TimeFormat)
		req.Header.Set("If-Modified-Since", t)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		f, err := os.Create(outpath)
		if err != nil {
			return err
		}
		defer f.Close()
		if _, err := io.Copy(f, resp.Body); err != nil {
			return err
		}

	case http.StatusNotModified:
		return errorNotModified

	default:
		return fmt.Errorf("unexpected response: %s", resp.Status)
	}
	return nil
}

func extract(dir, zipName, recipeName string) error {
	prev, err := loadFileInfo(recipeName)
	if err != nil {
		log.Printf("WARN: failed to load recipe: %s", err)
		log.Println("INFO: try to extract all files")
		prev = make(fileInfoTable)
	}
	curr, err := extractZip(zipName, dir, 0, prev)
	if err != nil {
		return err
	}
	if err := curr.save(recipeName); err != nil {
		log.Printf("WARN: failed to save recipe: %s", err)
	}
	curr.clean(dir, prev)
	return nil
}

func update(c *config, srcURL string) error {
	dp, err := c.downloadPath(srcURL)
	if err != nil {
		return err
	}
	old_anchor, err := c.anchor()
	if err != nil {
		return err
	}
	if err := download(srcURL, dp, old_anchor); err != nil {
		if err == errorNotModified {
			return nil
		}
		return err
	}
	new_anchor := time.Now()
	if err := extract(c.targetDir, dp, c.recipePath()); err != nil {
		return err
	}
	if err := c.updateAnchor(new_anchor); err != nil {
		os.Remove(c.anchorPath())
		return err
	}
	if err := os.Remove(dp); err != nil {
		log.Printf("WARN: failed to remove: %s", err)
	}
	return nil
}

func determineSourceURL(c *config) (string, error) {
	r, err := github.Latest("vim", "vim-win32-installer")
	if err != nil {
		return "", err
	}
	if r.Draft || r.PreRelease {
		return "", errors.New("absence of proper latest release")
	}
	var target *github.Asset
	for _, a := range r.Assets {
		if !strings.HasPrefix(a.Name, "gvim_7.4.") {
			continue
		}
		if (c.cpu == arch.X86 && strings.HasSuffix(a.Name, "_x86.zip")) ||
			(c.cpu == arch.AMD64 && strings.HasSuffix(a.Name, "_x64.zip")) {
			target = &a
			break
		}
	}
	if target == nil {
		return "", errors.New("no assets for arch in latest release")
	}
	if target.State != "uploaded" {
		return "", fmt.Errorf("new release's status error: %s", target.State)
	}
	prev_anchor, err := c.anchor()
	if err != nil {
		return "", err
	}
	fmt.Printf("prev: %s\nupdated: %s\nbefore: %v\n", prev_anchor, target.UpdatedAt, prev_anchor.Before(target.UpdatedAt))
	if prev_anchor.After(target.UpdatedAt) {
		return "", errors.New("no updated assets in release")
	}
	return target.DownloadURL, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("USAGE: upvim {TARGET_DIR}")
		os.Exit(1)
	}
	c, err := newConfig(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	if err := c.prepare(); err != nil {
		log.Fatal(err)
	}
	srcURL, err := determineSourceURL(c)
	if err != nil {
		log.Fatal(err)
	}
	if err := update(c, srcURL); err != nil {
		log.Fatal(err)
	}
}
