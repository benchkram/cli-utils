package app

import "fmt"

// Application is the interface for the application
type Application interface {
	Version() string
}

// app is the implementation of the application
type app struct {
	versionInfo VersionInfo
}

// VersionInfo contains the version and commit hash
type VersionInfo struct {
	Version string
	Commit  string
}

// NewApplication creates a new application
func NewApplication(version VersionInfo) Application {
	return &app{
		versionInfo: version,
	}
}

// Version returns the version and commit hash
func (a *app) Version() string {
	return fmt.Sprintf("%s (%s)", a.versionInfo.Version, a.versionInfo.Commit)
}
