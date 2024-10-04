package model

import (
	"fmt"
	"time"

	"github.com/gosuri/uitable"
)

type Skin struct {
	Name     string
	Modified time.Time
}

func (i *Skin) String() string {
	return fmt.Sprintf("%s\t%s", i.Name, i.Modified)
}

type SkinList []*Skin

func (l SkinList) String() string {
	table := uitable.New()
	table.AddRow("NAME", "MODIFIED")
	for _, skin := range l {
		table.AddRow(skin.Name, skin.Modified)
	}
	return table.String()
}
