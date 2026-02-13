package main

import (
	"sort"

	todoist "github.com/sachaos/todoist/lib"
	"github.com/urfave/cli/v2"
)

func Sections(c *cli.Context) error {
	client := GetClient(c)

	colorList := ColorList()
	var projectIds []string
	for _, project := range client.Store.Projects {
		projectIds = append(projectIds, project.GetID())
	}
	projectColorHash := GenerateColorHash(projectIds, colorList)

	sections := make([]todoist.Section, len(client.Store.Sections))
	copy(sections, client.Store.Sections)
	sort.Slice(sections, func(i, j int) bool {
		if sections[i].ProjectID == sections[j].ProjectID {
			return sections[i].SectionOrder < sections[j].SectionOrder
		}
		return sections[i].ProjectID < sections[j].ProjectID
	})

	itemList := [][]string{}
	for _, section := range sections {
		if section.IsDeleted || section.IsArchived {
			continue
		}
		itemList = append(itemList, []string{
			IdFormat(section),
			section.Name,
			ProjectFormat(section.ProjectID, client.Store, projectColorHash, c),
		})
	}

	defer writer.Flush()

	if c.Bool("header") {
		writer.Write([]string{"ID", "Name", "Project"})
	}

	for _, strings := range itemList {
		writer.Write(strings)
	}

	return nil
}
