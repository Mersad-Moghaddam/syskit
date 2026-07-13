package command

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Mersad-Moghaddam/syskit/internal/model"
	"github.com/Mersad-Moghaddam/syskit/internal/render"
	"github.com/Mersad-Moghaddam/syskit/internal/service"
	"github.com/spf13/cobra"
)

type ProcessService interface {
	List(service.ProcessOptions) (*model.ProcessList, error)
	Sample(time.Duration, service.ProcessOptions) (*model.ProcessList, error)
	Tree() ([]service.ProcessTreeNode, error)
}
type ProcessOptions struct {
	Format   func() string
	NoHeader func() bool
}

func NewProcessCmd(s ProcessService, o ProcessOptions) *cobra.Command {
	var filters []string
	var sort string
	var reverse bool
	var limit, pid int
	var user string
	var interval time.Duration
	cmd := &cobra.Command{Use: "process", Short: "List processes from procfs", Args: cobra.NoArgs, RunE: func(c *cobra.Command, args []string) error {
		raw := append([]string(nil), filters...)
		if pid > 0 {
			raw = append(raw, "pid="+strconv.Itoa(pid))
		}
		if user != "" {
			raw = append(raw, "user="+user)
		}
		parsed, err := service.ParseProcessFilters(raw)
		if err != nil {
			return err
		}
		options := service.ProcessOptions{Filters: parsed, Sort: sort, Reverse: reverse, Limit: limit}
		var list *model.ProcessList
		if interval > 0 {
			list, err = s.Sample(interval, options)
		} else {
			list, err = s.List(options)
		}
		if err != nil {
			return fmt.Errorf("collecting processes: %w", err)
		}
		r, err := render.New(o.Format(), render.WithNoHeader(o.NoHeader()))
		if err != nil {
			return err
		}
		if o.Format() == "table" {
			return r.Render(c.OutOrStdout(), processTable(list))
		}
		return r.Render(c.OutOrStdout(), list)
	}}
	cmd.Flags().StringSliceVar(&filters, "filter", nil, "filter with field=value (repeatable)")
	cmd.Flags().StringVar(&sort, "sort", "pid", "sort by pid, cpu, memory, or name")
	cmd.Flags().BoolVar(&reverse, "reverse", false, "reverse sort order")
	cmd.Flags().IntVar(&limit, "limit", 0, "maximum results (0 is unlimited)")
	cmd.Flags().IntVar(&pid, "pid", 0, "filter by PID")
	cmd.Flags().StringVar(&user, "user", "", "filter by user name or UID")
	cmd.Flags().DurationVar(&interval, "interval", 0, "sample process CPU usage over this interval")
	cmd.AddCommand(&cobra.Command{Use: "tree", Short: "Show processes as a parent-child tree", Args: cobra.NoArgs, RunE: func(c *cobra.Command, args []string) error {
		tree, err := s.Tree()
		if err != nil {
			return fmt.Errorf("collecting process tree: %w", err)
		}
		if o.Format() != "table" {
			r, err := render.New(o.Format())
			if err != nil {
				return err
			}
			return r.Render(c.OutOrStdout(), tree)
		}
		for _, node := range tree {
			writeTree(c, node, "")
		}
		return nil
	}})
	return cmd
}
func writeTree(c *cobra.Command, node service.ProcessTreeNode, prefix string) {
	fmt.Fprintf(c.OutOrStdout(), "%s%d %s\n", prefix, node.Process.PID, node.Process.Command)
	for _, child := range node.Children {
		writeTree(c, *child, prefix+"  ")
	}
}
func processTable(list *model.ProcessList) render.Table {
	t := render.Table{Headers: []string{"PID", "PPID", "USER", "STATE", "CPU %", "MEM %", "RSS BYTES", "START TICKS", "THREADS", "COMMAND"}}
	for _, p := range list.Processes {
		user := p.User
		if user == "" {
			user = strconv.FormatUint(p.UID, 10)
		}
		cpu, memory := "unavailable", "unavailable"
		if p.CPUPercent != nil {
			cpu = fmt.Sprintf("%.1f", *p.CPUPercent)
		}
		if p.MemoryPercent != nil {
			memory = fmt.Sprintf("%.1f", *p.MemoryPercent)
		}
		t.Rows = append(t.Rows, []string{strconv.Itoa(p.PID), strconv.Itoa(p.PPID), user, p.State, cpu, memory, strconv.FormatUint(p.ResidentBytes, 10), strconv.FormatUint(p.StartTimeTicks, 10), strconv.FormatUint(p.Threads, 10), p.Command})
	}
	return t
}
