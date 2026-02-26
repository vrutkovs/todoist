package todoist

import "context"

type Label struct {
	HaveID
	Color      string `json:"color"`
	IsDeleted  bool   `json:"is_deleted"`
	IsPersonal bool   `json:"is_personal"`
	ItemOrder  int    `json:"item_order"`
	Name       string `json:"name"`
}

func (label Label) AddParam() interface{} {
	return map[string]interface{}{
		"name":        label.Name,
		"is_personal": true,
	}
}

type Labels []Label

func (a Labels) Len() int           { return len(a) }
func (a Labels) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Labels) Less(i, j int) bool { return a[i].ID < a[j].ID }

func (a Labels) At(i int) IDCarrier { return a[i] }

func (a Labels) GetIDByName(name string) string {
	for _, label := range a {
		if label.Name == name {
			return label.ID
		}
	}
	return ""
}

func (a Labels) GetByName(name string) *Label {
	for i, label := range a {
		if label.Name == name {
			return &a[i]
		}
	}
	return nil
}

func (c *Client) AddLabel(ctx context.Context, label Label) error {
	commands := Commands{
		NewCommand("label_add", label.AddParam()),
	}
	return c.ExecCommands(ctx, commands)
}

// EnsureLabelsExist creates any labels in labelNames that do not yet exist in
// the store as personal labels, then returns nil. Call this before AddItem /
// UpdateItem so that all referenced labels are present on the server.
// labelNames should already be normalised (names, not IDs) via
// Store.NormalizeLabelNames before being passed here.
func (c *Client) EnsureLabelsExist(ctx context.Context, labelNames []string) error {
	var commands Commands
	for _, name := range labelNames {
		if name == "" {
			continue
		}
		// Skip if already present by name or by ID (old-cache safety net).
		if c.Store.Labels.GetByName(name) != nil || c.Store.FindLabel(name) != nil {
			continue
		}
		commands = append(commands, NewCommand("label_add", Label{Name: name}.AddParam()))
	}
	if len(commands) == 0 {
		return nil
	}
	return c.ExecCommands(ctx, commands)
}
