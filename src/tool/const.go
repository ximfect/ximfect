package tool

import (
	_ "embed"
)

var (
	// Version is the release number
	//go:embed version
	Version string

	// Build is the build number
	//go:embed build
	Build string

	// NoDesc is a generic string shown when something has no description
	NoDesc = "No description :("
)
