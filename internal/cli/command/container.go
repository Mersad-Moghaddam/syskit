package command

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Mersad-Moghaddam/syskit/internal/model"
	"github.com/Mersad-Moghaddam/syskit/internal/render"
)

type ContainerService interface {
	List() (*model.ContainerList, error)
	Inspect(string) (*model.ContainerDetail, error)
}

type ContainerOptions struct {
	Format   func() string
	NoHeader func() bool
	Color    func() bool
}

func NewContainerCmd(s ContainerService, o ContainerOptions) *cobra.Command {
	cmd := &cobra.Command{Use: "containers", Short: "Show cgroup-associated containers", Args: cobra.NoArgs, RunE: func(c *cobra.Command, args []string) error {
		containers, err := s.List()
		if err != nil {
			return fmt.Errorf("collecting containers: %w", err)
		}
		r, err := render.New(o.Format(), render.WithNoHeader(o.NoHeader()), render.WithColor(o.Color()))
		if err != nil {
			return err
		}
		if o.Format() == "table" {
			return r.Render(c.OutOrStdout(), containerTable(containers))
		}
		return r.Render(c.OutOrStdout(), containers)
	}}
	cmd.AddCommand(&cobra.Command{Use: "inspect <id>", Short: "Show processes associated with one container", Args: cobra.ExactArgs(1), RunE: func(c *cobra.Command, args []string) error {
		detail, err := s.Inspect(args[0])
		if err != nil {
			return fmt.Errorf("inspecting container: %w", err)
		}
		r, err := render.New(o.Format(), render.WithNoHeader(o.NoHeader()), render.WithColor(o.Color()))
		if err != nil {
			return err
		}
		if o.Format() == "table" {
			table := render.Table{Headers: []string{"CONTAINER", "RUNTIME", "PID", "COMMAND"}}
			for _, process := range detail.Processes {
				table.Rows = append(table.Rows, []string{detail.ID, detail.Runtime, fmt.Sprint(process.PID), process.Command})
			}
			return r.Render(c.OutOrStdout(), table)
		}
		return r.Render(c.OutOrStdout(), detail)
	}})
	return cmd
}

func containerTable(containers *model.ContainerList) render.Table {
	t := render.Table{Headers: []string{"CONTAINER", "RUNTIME", "PIDS", "MEMORY", "CPU NS", "READ", "WRITE"}}
	for _, container := range containers.Containers {
		t.Rows = append(t.Rows, []string{container.ID, container.Runtime, fmt.Sprint(container.PIDs), containerMetric(container.Metrics, func(m *model.ContainerMetrics) *uint64 { return m.MemoryCurrentBytes }), containerMetric(container.Metrics, func(m *model.ContainerMetrics) *uint64 { return m.CPUUsageNanoseconds }), containerMetric(container.Metrics, func(m *model.ContainerMetrics) *uint64 { return m.ReadBytes }), containerMetric(container.Metrics, func(m *model.ContainerMetrics) *uint64 { return m.WrittenBytes })})
	}
	return t
}

func containerMetric(metrics *model.ContainerMetrics, field func(*model.ContainerMetrics) *uint64) string {
	if metrics == nil || field(metrics) == nil {
		return "-"
	}
	return fmt.Sprint(*field(metrics))
}
