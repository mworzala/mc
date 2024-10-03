package model

import (
	"fmt"
	"time"

	"github.com/gosuri/uitable"
)

type Skin struct {
	Name string
	Date time.Time
}

func (i *Skin) String() string {
	return fmt.Sprintf("%s\t%s", i.Name, i.Date)
}

type SkinList []*Skin

func (l SkinList) String() string {
	table := uitable.New()
	table.AddRow("NAME", "DATE")
	for _, skin := range l {
		table.AddRow(skin.Name, skin.Date)
	}
	return table.String()
}
