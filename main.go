package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/mworzala/mc/cmd/mc"
	"github.com/mworzala/mc/internal/pkg/cli"
)

var (
	version  = "dev"
	commit   = "none"
	date     = "unknown"
	modified = false
)

func main() {
	//goland:noinspection GoBoolExpressions version is set using ldflags
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
