// mv_ldown.go: moves the latest download from
// ~/Downloads to the current directory
package main

import (
	"os"
	"path"
	"fmt"
	"io/ioutil"
	"sort"
	"runtime"
)

type Files []*os.FileInfo

// Sort by last modified time, from most recent to oldest
func (f Files) Len() int           { return len(f) }
func (f Files) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
func (f Files) Less(i, j int) bool { return f[i].Mtime_ns > f[j].Mtime_ns }

func errExit(msg string) {
	fmt.Printf(msg + "\n")
	os.Exit(1)
}

func main() {
	if runtime.GOOS != "darwin" {
		errExit("runs only on Mac OS X")
	}

	downloadsDir := path.Join(os.Getenv("HOME"), "Downloads")
	downloads, err := ioutil.ReadDir(downloadsDir)
	if err != nil {
		errExit("error: " + err.String())
	}
	sort.Sort(Files(downloads))

	base := ""
	for _, fi := range downloads {
		if fi.Name == ".DS_Store" {
			continue
		}
		base = fi.Name
		break
	}

	if base == "" {
		errExit("nothing in ~/Downloads")
	}

	abs := path.Join(downloadsDir, base)

	// Don't overwrite existing file of same
	// name in current directory.
	if _, err = os.Stat(base); err == nil {
		errExit("./" + base + " already exists")
	}

	os.Rename(abs, base)
	fmt.Printf("moved %s to ./%s\n", abs, base)
}
