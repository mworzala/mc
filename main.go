package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/mworzala/mc-cli/cmd/mc"
	"github.com/mworzala/mc-cli/internal/pkg/cli"
)

var (
	version  = "dev"
	commit   = "none"
	date     = "unknown"
	modified = false
)

func main() {
	// version is set using ldflags, so a complaint about constant condition here is wrong
	if version == "dev" {
		// Not built with goreleaser, so we should try to read go build info
		if buildInfo, ok := debug.ReadBuildInfo(); ok {
			for _, setting := range buildInfo.Settings {
				switch setting.Key {
				case "vcs.revision":
					commit = setting.Value
				case "vcs.time":
					date = setting.Value
				case "vcs.modified":
					modified = setting.Value != "false"
				}
			}
		}
	}

	app := cli.NewApp(cli.BuildInfo{Version: version, Commit: commit, Date: date, Modified: modified})
	rootCmd := mc.NewRootCmd(app)

	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
