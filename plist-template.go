package main

import (
	_ "embed"
	"log/slog"
	"os"
	"path"
	"text/template"

	"github.com/pkg/errors"
)

//go:embed templates/dev.localhost.iCloudForceSync.plist.template
var plistTemplate []byte

type plist struct {
	path string
}

// NewPlist creates a new plist instance with the given file path.
func newPlist() (plist, error) {
	p := plist{}
	p, err := p.setLaunchAgentsPath()
	if err != nil {
		return p, errors.Wrapf(err, "cannot create plist instance")
	}

	return p, nil
}

// isExisting checks if the path in the plist struct exists.
func (p plist) isExisting() bool {
	_, err := os.Stat(p.path)
	return err == nil
}

// setLaunchAgentsPath sets the path for the plist file in the user's LaunchAgents directory.
func (p plist) setLaunchAgentsPath() (plist, error) {
	l := slog.Default()
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return p, errors.Wrapf(err, "failed to get home directory")
	}
	p.path = path.Join(homeDir, "/Library/LaunchAgents/dev.localhost.iCloudForceSync.plist")
	l.Debug("plist path set", "path", p.path)
	return p, nil
}

func (p plist) currentExecutablePath() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", errors.Wrapf(err, "failed to get executable path")
	}

	return ex, nil
}

// remove deletes the plist file specified by the path in the plist struct. Returns an error if the removal fails.
func (p plist) remove() error {
	l := slog.Default().With("path", p.path)
	if err := os.Remove(p.path); err != nil {
		return errors.Wrapf(err, "failed to remove plist file: %s", p.path)
	}

	l.Debug("plist file removed")
	return nil
}

// render generates the plist file from a template, removing the old file if it exists, and saving the new one to disk.
func (p plist) render() error {
	l := slog.Default().With("path", p.path)
	t, err := template.New("plist").Parse(string(plistTemplate))
	if err != nil {
		return errors.Wrapf(err, "failed to parse template")
	}

	if p.isExisting() {
		if err := p.remove(); err != nil {
			return err
		}
	}

	f, err := os.Create(p.path)
	if err != nil {
		return errors.Wrapf(err, "failed to create plist file: %s", p.path)
	}

	binaryPath, err := p.currentExecutablePath()
	if err != nil {
		return errors.Wrapf(err, "failed to get current path")
	}
	l = l.With("binary_path", binaryPath)

	if err := t.Execute(f, binaryPath); err != nil {
		return errors.Wrapf(err, "failed to execute template and save to file: %s", p.path)
	}
	l.Info("plist file created")

	return nil
}
