package main

import (
	"context"
	"strings"

	todoist "github.com/sachaos/todoist/lib"
	"github.com/urfave/cli/v2"
)

func Modify(c *cli.Context) error {
	client := GetClient(c)

	if !c.Args().Present() {
		return CommandFailed
	}

	var err error
	item_id, err := client.CompleteItemIDByPrefix(c.Args().First())
	if err != nil {
		return err
	}
	item := client.Store.FindItem(item_id)
	if item == nil {
		return IdNotFound
	}
	item.Content = c.String("content")
	item.Priority = priorityMapping[c.Int("priority")]
	item.LabelNames = func(str string) []string {
		stringNames := strings.Split(str, ",")
		names := []string{}
		for _, stringName := range stringNames {
			names = append(names, stringName)
		}
		return names
	}(c.String("label-names"))

	item.Due = &todoist.Due{String: c.String("date")}

	projectID := c.String("project-id")
	if projectID == "" {
		projectID = client.Store.Projects.GetIDByName(c.String("project-name"))
	}

	sectionID := c.String("section-id")
	if sectionID == "" {
		sectionName := c.String("section-name")
		if sectionName != "" {
			pID := projectID
			if pID == "" {
				pID = item.ProjectID
			}
			for _, s := range client.Store.Sections {
				if s.Name == sectionName && s.ProjectID == pID {
					sectionID = s.ID
					if projectID == "" {
						projectID = s.ProjectID
					}
					break
				}
			}
		}
	} else {
		section := client.Store.FindSection(sectionID)
		if section != nil && projectID == "" {
			projectID = section.ProjectID
		}
	}

	if !c.Args().Present() {
		return CommandFailed
	}

	if err := client.UpdateItem(context.Background(), *item); err != nil {
		return err
	}

	if err := client.MoveItem(context.Background(), item, projectID, sectionID); err != nil {
		return err
	}

	return Sync(c)
}
