package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_plist_currentPath(t *testing.T) {
	p, err := newPlist()
	require.NoError(t, err)
	binaryPath, err := p.currentExecutablePath()
	require.NoError(t, err)
	require.NotEmpty(t, binaryPath)
}

func Test_plist_isExisting(t *testing.T) {
	t.Run("should return true for $HOME", func(t *testing.T) {
		homeDir, err := os.UserHomeDir()
		require.NoError(t, err)
		p, err := newPlist()
		require.NoError(t, err)
		p.path = homeDir

		require.True(t, p.isExisting())
	})
}

func Test_plist_setLaunchAgentsPath(t *testing.T) {
	t.Run("should set the path to the launch agents", func(t *testing.T) {
		p, err := newPlist()
		require.NoError(t, err)
		p, err = p.setLaunchAgentsPath()
		require.NoError(t, err)
		require.NotEmpty(t, p.path)
		assert.Contains(t, p.path, "/Library/LaunchAgents")
		assert.Contains(t, p.path, "dev.localhost.iCloudForceSync.plist")
	})
}

func Test_plist_currentExecutablePath(t *testing.T) {
	t.Run("should not be empty", func(t *testing.T) {
		p, err := newPlist()
		require.NoError(t, err)
		path, err := p.currentExecutablePath()
		require.NoError(t, err)
		require.NotEmpty(t, path)
	})
}

func Test_plist_remove(t *testing.T) {
	t.Run("should remove temporary file", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "test*")
		require.NoError(t, err)
		p, err := newPlist()
		require.NoError(t, err)
		p.path = tmpFile.Name()
		require.NoError(t, p.remove())
		_, err = os.Stat(tmpFile.Name())
		assert.Error(t, err)
	})
}

func Test_plist_render(t *testing.T) {
	t.Run("should render the template", func(t *testing.T) {
		p, err := newPlist()
		require.NoError(t, err)

		tmpFile, err := os.CreateTemp("", "test*")
		require.NoError(t, err)

		defer func() {
			require.NoError(t, tmpFile.Close())
			require.NoError(t, os.Remove(tmpFile.Name()))
		}()

		p.path = tmpFile.Name()
		require.NoError(t, p.render())

		expected := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">`
		content, err := os.ReadFile(tmpFile.Name())
		require.NoError(t, err)
		assert.Contains(t, string(content), expected)
	})
}
