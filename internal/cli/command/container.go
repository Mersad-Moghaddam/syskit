package command

import (
	"fmt"

	"github.com/Mersad-Moghaddam/syskit/internal/model"
	"github.com/Mersad-Moghaddam/syskit/internal/render"
	"github.com/spf13/cobra"
)

type ContainerService interface {
	List() (*model.ContainerList, error)
}

type ContainerOptions struct {
	Format   func() string
	NoHeader func() bool
}

func NewContainerCmd(s ContainerService, o ContainerOptions) *cobra.Command {
	return &cobra.Command{Use: "containers", Short: "Show cgroup-associated containers", Args: cobra.NoArgs, RunE: func(c *cobra.Command, args []string) error {
		containers, err := s.List()
		if err != nil {
			return fmt.Errorf("collecting containers: %w", err)
		}
		r, err := render.New(o.Format(), render.WithNoHeader(o.NoHeader()))
		if err != nil {
			return err
		}
		if o.Format() == "table" {
			return r.Render(c.OutOrStdout(), containerTable(containers))
		}
		return r.Render(c.OutOrStdout(), containers)
	}}
}

func containerTable(containers *model.ContainerList) render.Table {
	t := render.Table{Headers: []string{"CONTAINER", "RUNTIME", "PIDS"}}
	for _, container := range containers.Containers {
		t.Rows = append(t.Rows, []string{container.ID, container.Runtime, fmt.Sprint(container.PIDs)})
	}
	return t
}
