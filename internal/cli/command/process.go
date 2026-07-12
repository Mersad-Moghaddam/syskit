package command

import (
	"fmt"
	"strconv"

	"github.com/Mersad-Moghaddam/syskit/internal/model"
	"github.com/Mersad-Moghaddam/syskit/internal/render"
	"github.com/Mersad-Moghaddam/syskit/internal/service"
	"github.com/spf13/cobra"
)

type ProcessService interface {
	List(service.ProcessOptions) (*model.ProcessList, error)
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
		list, err := s.List(service.ProcessOptions{Filters: parsed, Sort: sort, Reverse: reverse, Limit: limit})
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
	cmd.Flags().StringVar(&user, "user", "", "filter by UID")
	return cmd
}
func processTable(list *model.ProcessList) render.Table {
	t := render.Table{Headers: []string{"PID", "PPID", "UID", "STATE", "CPU TICKS", "RSS BYTES", "THREADS", "COMMAND"}}
	for _, p := range list.Processes {
		t.Rows = append(t.Rows, []string{strconv.Itoa(p.PID), strconv.Itoa(p.PPID), strconv.FormatUint(p.UID, 10), p.State, strconv.FormatUint(p.CPUTime, 10), strconv.FormatUint(p.ResidentBytes, 10), strconv.FormatUint(p.Threads, 10), p.Command})
	}
	return t
}
