package utils

import (
	_ "embed"
	"strings"
)

const (
	VersionSubcommandName = "version"
)

//go:embed VERSION
var version string

func IsDevelopment() bool {
	return version == "development"
}

func GetCurrentVersion() string {
	if version != "" {
		return version
	}

	return "unknown"
}

func GetVersionRaw() (semantic string, commit string) {
	if version == "" {
		return "unknown", ""
	}
	if IsDevelopment() {
		return version, "in-progress"
	} else {
		parts := strings.Split(
			StripMany(
				version,
				"\\",
				"(",
				")",
				"\"",
			), " ")

		if len(parts) >= 2 {
			return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		}
		return version, ""
	}
}
