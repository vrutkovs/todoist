package main

import (
	"context"
	"fmt"
	"strings"

	todoist "github.com/sachaos/todoist/lib"
	"github.com/urfave/cli/v2"
)

var priorityMapping = map[int]int{
	1: 4,
	2: 3,
	3: 2,
	4: 1,
}

func Add(c *cli.Context) error {
	client := GetClient(c)

	item := todoist.Item{}
	if c.Args().Len() != 1 {
		return fmt.Errorf("add command requires 1 positional argument for the task title, but got %v.", c.Args().Len())
	}

	item.Content = c.Args().First()

	item.Priority = priorityMapping[c.Int("priority")]

	projectName := c.String("project-name")
	if projectName != "" {
		projectId := client.Store.Projects.GetIDByName(projectName)
		if projectId == "" {
			return fmt.Errorf("Did not find a project named '%v'", projectName)
		}
		item.ProjectID = projectId
	} else {
		item.ProjectID = c.String("project-id")
	}

	sectionName := c.String("section-name")
	item.SectionID = c.String("section-id")
	if sectionName != "" {
		var sectionID string
		for _, section := range client.Store.Sections {
			if section.Name == sectionName {
				if item.ProjectID == "" || item.ProjectID == section.ProjectID {
					sectionID = section.ID
					if item.ProjectID == "" {
						item.ProjectID = section.ProjectID
					}
					break
				}
			}
		}
		if sectionID == "" {
			return fmt.Errorf("did not find a section named '%v'", sectionName)
		}
		item.SectionID = sectionID
	}

	item.LabelNames = func(str string) []string {
		stringNames := strings.Split(str, ",")
		names := []string{}
		for _, stringName := range stringNames {
			names = append(names, stringName)
		}
		return names
	}(c.String("label-names"))

	item.LabelNames = client.Store.NormalizeLabelNames(item.LabelNames)

	item.Due = &todoist.Due{String: c.String("date")}

	item.AutoReminder = c.Bool("reminder")

	if err := client.EnsureLabelsExist(context.Background(), item.LabelNames); err != nil {
		return err
	}

	if err := client.AddItem(context.Background(), item); err != nil {
		return err
	}

	return Sync(c)
}
