package model

import (
	"fmt"
	"time"

	"github.com/MakeNowJust/heredoc"
)

type Version struct {
	Tag      string
	Commit   string
	State    string
	Date     string
	Go       string
	Platform string
	Game     GameVersion
}

type GameVersion struct {
	Release  string
	Snapshot string
}

func (v *Version) String() string {
	if v.Tag == "dev" {
		modified := ""
		if v.State != "clean" {
			modified = "*"
		}
		return heredoc.Docf(`
			mc-cli dev%s
			minecraft %s (%s)`,
			modified, v.Game.Release, v.Game.Snapshot)
	}

	date := ""
	if t, err := time.Parse(time.RFC3339, v.Date); err == nil {
		date = fmt.Sprintf(" (%s)", t.Format("2006-01-02"))
	}

	return heredoc.Docf(`
		mc-cli version %s%s
		https://github.com/mworzala/mc/releases/tag/v%s
		minecraft %s (%s)`,
		v.Tag, date, v.Tag, v.Game.Release, v.Game.Snapshot)
}
