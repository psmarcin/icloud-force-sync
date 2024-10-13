package main

import (
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"
)

func main() {
	l := slog.Default()
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	pl, err := newPlist()
	if err != nil {
		panic(err)
	}
	if err := pl.render(); err != nil {
		panic(err)
	}

	iCloudDirectory := path.Join(home, "/Library/Mobile Documents/com~apple~CloudDocs")
	_, err = os.ReadDir(iCloudDirectory)
	if os.IsNotExist(err) {
		panic("iCloud directory not found")
	}
	if err != nil {
		panic(err)
	}
	l.Info("starting", "icloud_dir", iCloudDirectory)

	tmpFilesDir := path.Join(iCloudDirectory, "iCloudForceSync")
	if err := os.MkdirAll(tmpFilesDir, 0755); err != nil {
		panic(err)
	}

	f, err := newFile(tmpFilesDir)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanup(tmpFilesDir)
		os.Exit(1)
	}()

	for {
		f, err = f.Touch()
		if err != nil {
			panic(err)
		}
	}
}

type file struct {
	dir       string
	file      *os.File
	maxWrites int
	writes    int
	sleep     time.Duration
	ttl       time.Time
}

func newFile(tmpFilesDir string) (file, error) {
	l := slog.Default()

	ttl := time.Second * time.Duration(rand.Intn(180))
	f := file{
		dir:       tmpFilesDir,
		ttl:       time.Now().Add(ttl),
		maxWrites: rand.Intn(100),
		sleep:     time.Second * time.Duration(rand.Intn(15)+10),
	}

	tmpFileName := fmt.Sprintf("%s-*.txt", time.Now().Format("2006-01-02:15:04:05"))
	tmpFile, err := os.CreateTemp(f.dir, tmpFileName)
	if err != nil {
		return f, err
	}

	f.file = tmpFile
	l.Debug("created temporary file", "path", path.Join(f.dir, tmpFile.Name()), "ttl", ttl, "max_writes", f.maxWrites, "sleep", f.sleep)

	return f, nil
}

func (f file) Touch() (file, error) {
	l := slog.Default()

	f, err := f.recreate()
	if err != nil {
		return f, err
	}
	if _, err := f.file.WriteString(time.Now().Format("2006-01-02:15:04:05\n")); err != nil {
		return f, err
	}
	f.writes++
	l.Debug("touch", "path", f.file.Name(), "sleep", f.sleep, "ttl", f.ttl, "max_writes", f.maxWrites, "writes", f.writes)
	time.Sleep(f.sleep)

	return f, nil
}

func (f file) Close() error {
	if err := f.file.Close(); err != nil {
		return err
	}

	if err := os.Remove(f.file.Name()); err != nil {
		return err
	}
	return nil
}

func (f file) recreate() (file, error) {
	if f.ttl.After(time.Now()) && f.writes <= f.maxWrites {
		return f, nil
	}

	if err := f.Close(); err != nil {
		return f, err
	}

	f, err := newFile(f.dir)
	if err != nil {
		return f, err
	}

	return f, nil
}

func cleanup(tmpFilesDir string) {
	l := slog.Default()
	l.Info("cleaning up", "dir", tmpFilesDir, "pattern", "*")
	files, err := os.ReadDir(tmpFilesDir)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if err := os.Remove(path.Join(tmpFilesDir, file.Name())); err != nil {
			panic(err)
		}
	}
	l.Info("cleaned up", "dir", tmpFilesDir)
}
