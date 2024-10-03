package model

import (
	"fmt"

	"github.com/gosuri/uitable"
	"github.com/mworzala/mc/internal/pkg/modrinth"
	"github.com/mworzala/mc/internal/pkg/util"
)

type ModrinthSearchResult modrinth.SearchResponse

func (result *ModrinthSearchResult) String() string {
	table := uitable.New()
	table.AddRow("ID", "TYPE", "NAME", "DOWNLOADS")
	for _, project := range result.Hits {
		table.AddRow(project.ProjectID, project.ProjectType, project.Title, util.FormatCount(project.Downloads))
	}
	res := table.String()
	if result.TotalHits-len(result.Hits) > 0 {
		res += fmt.Sprintf("\n...and %d more", result.TotalHits-len(result.Hits))
	}
	return res
}
