package model

import (
	"fmt"

	"github.com/gosuri/uitable"
)

type ProfileType string

const (
	ProfileTypeVanilla ProfileType = "vanilla"
	ProfileTypeFabric  ProfileType = "fabric"
)

var ProfileTypes = []ProfileType{
	ProfileType("unknown"),
	ProfileTypeVanilla,
	ProfileTypeFabric,
}

type Profile struct {
	Name      string
	Directory string

	Type    ProfileType
	Version string
}

func (p *Profile) String() string {
	return fmt.Sprintf("%s (%s %s)", p.Name, p.Type, p.Version)
}

type ProfileList []*Profile

func (l ProfileList) String() string {
	table := uitable.New()
	table.AddRow("NAME", "TYPE", "VERSION")
	for _, profile := range l {
		//todo -o wide
		table.AddRow(profile.Name, profile.Type, profile.Version)
	}
	return table.String()
}
