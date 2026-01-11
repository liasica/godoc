// Copyright (C) godoc. 2026-present.
//
// Created at 2026-01-10, by liasica

package godoc

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"
)

var (
	Version    = "v1.0.0"
	BuildTime  = ""
	CommitHash = ""
)

func FullVersion() string {
	appName := "godoc"

	// runtime info
	goVer := runtime.Version()
	compiler := runtime.Compiler
	goos := runtime.GOOS
	goarch := runtime.GOARCH

	// build info defaults
	from := "unknown"
	mainVer := ""
	mainSum := ""
	vcsInfo := map[string]string{}
	// depsList := []string{}

	if bi, ok := debug.ReadBuildInfo(); ok {
		if bi.GoVersion != "" {
			goVer = bi.GoVersion
		}
		if bi.Main.Path != "" {
			from = bi.Main.Path
		}
		if bi.Main.Version != "" {
			mainVer = bi.Main.Version
		}
		if bi.Main.Sum != "" {
			mainSum = bi.Main.Sum
		}
		for _, s := range bi.Settings {
			vcsInfo[s.Key] = s.Value
		}
		// for i, d := range bi.Deps {
		// 	if i >= 20 {
		// 		depsList = append(depsList, fmt.Sprintf("...and %d more", len(bi.Deps)-i))
		// 		break
		// 	}
		// 	if d.Sum != "" {
		// 		depsList = append(depsList, fmt.Sprintf("%s@%s (sum: %s)", d.Path, d.Version, d.Sum))
		// 	} else {
		// 		depsList = append(depsList, fmt.Sprintf("%s@%s", d.Path, d.Version))
		// 	}
		// }
	}

	// fallbacks
	if mainVer == "" {
		mainVer = "(unknown)"
	}
	// leave mainSum empty if unavailable; some build environments set it to "unknown"
	if mainSum == "" || mainSum == "unknown" {
		mainSum = ""
	}
	if BuildTime == "" {
		BuildTime = "unknown"
	}
	if CommitHash == "" {
		CommitHash = "unknown"
	}

	// build a multi-line, readable report
	var b strings.Builder
	_, _ = fmt.Fprintf(&b, "%s has version %s (commit: %s)\n", appName, Version, CommitHash)
	_, _ = fmt.Fprintf(&b, "  built at: %s\n", BuildTime)
	_, _ = fmt.Fprintf(&b, "  go: %s (compiler: %s), platform: %s/%s\n", goVer, compiler, goos, goarch)
	_, _ = fmt.Fprintf(&b, "  module: %s %s\n", from, mainVer)
	// only print mod sum when it's present and meaningful
	if mainSum != "" {
		_, _ = fmt.Fprintf(&b, "    mod sum: %q\n", mainSum)
	}
	// VCS-related
	if len(vcsInfo) > 0 {
		_, _ = fmt.Fprint(&b, "  vcs settings:\n")
		keys := []string{"vcs", "vcs.revision", "vcs.time", "vcs.modified"}
		for _, k := range keys {
			if v, ok := vcsInfo[k]; ok {
				_, _ = fmt.Fprintf(&b, "    %s: %s\n", k, v)
			}
		}
		// any other settings
		for k, v := range vcsInfo {
			if k == "vcs" || k == "vcs.revision" || k == "vcs.time" || k == "vcs.modified" {
				continue
			}
			_, _ = fmt.Fprintf(&b, "    %s: %s\n", k, v)
		}
	}

	// deps
	// if len(depsList) > 0 {
	// 	_, _ = fmt.Fprint(&b, "  deps (first 20):\n")
	// 	for _, d := range depsList {
	// 		_, _ = fmt.Fprintf(&b, "    %s\n", d)
	// 	}
	// } else {
	// 	_, _ = fmt.Fprint(&b, "  deps: (none or unavailable)\n")
	// }

	return b.String()
}
