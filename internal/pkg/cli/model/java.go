package model

import (
	"fmt"

	"github.com/gosuri/uitable"
)

type JavaInstallation struct {
	Name    string
	Path    string
	Version int
}

func (i *JavaInstallation) String() string {
	return fmt.Sprintf("%s\t%s", i.Name, i.Path)
}

type JavaInstallationList []*JavaInstallation

func (l JavaInstallationList) String() string {
	table := uitable.New()
	table.AddRow("NAME", "EXECUTABLE PATH")
	for _, install := range l {
		//todo -o wide would output the arch and version also
		table.AddRow(install.Name, install.Path)
	}
	return table.String()
}
