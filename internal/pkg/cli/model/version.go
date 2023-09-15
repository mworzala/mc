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
}

func (v *Version) String() string {
	if v.Tag == "dev" {
		modified := ""
		if v.State != "clean" {
			modified = "*"
		}
		return fmt.Sprintf("mc-cli dev%s", modified)
	}

	date := ""
	if t, err := time.Parse(time.RFC3339, v.Date); err == nil {
		date = fmt.Sprintf(" (%s)", t.Format("2006-01-02"))
	}

	return heredoc.Docf(`
		mc-cli version %s%s
		https://github.com/mworzala/mc/releases/tag/v%s
`, v.Tag, date, v.Tag)
}
