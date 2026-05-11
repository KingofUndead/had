//go:build windows

package core

func runDaemon() error {
	logWarning("Daemon mode not supported on Windows, running in foreground")
	return nil
}
