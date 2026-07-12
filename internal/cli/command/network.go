package command

import (
	"fmt"

	"github.com/Mersad-Moghaddam/syskit/internal/model"
	"github.com/Mersad-Moghaddam/syskit/internal/render"
	"github.com/spf13/cobra"
)

type NetworkService interface {
	Collect() (*model.NetworkInfo, error)
}
type NetworkOptions struct {
	Format   func() string
	NoHeader func() bool
}

func NewNetworkCmd(s NetworkService, o NetworkOptions) *cobra.Command {
	return &cobra.Command{Use: "network", Short: "Show network interface counters", Args: cobra.NoArgs, RunE: func(c *cobra.Command, args []string) error {
		info, err := s.Collect()
		if err != nil {
			return fmt.Errorf("collecting network information: %w", err)
		}
		r, err := render.New(o.Format(), render.WithNoHeader(o.NoHeader()))
		if err != nil {
			return err
		}
		if o.Format() == "table" {
			return r.Render(c.OutOrStdout(), networkTable(info))
		}
		return r.Render(c.OutOrStdout(), info)
	}}
}
func networkTable(info *model.NetworkInfo) render.Table {
	t := render.Table{Headers: []string{"IFACE", "RX BYTES", "TX BYTES", "RX PACKETS", "TX PACKETS", "RX ERR", "TX ERR", "RX DROP", "TX DROP"}}
	for _, n := range info.Interfaces {
		t.Rows = append(t.Rows, []string{n.Name, fmt.Sprint(n.RXBytes), fmt.Sprint(n.TXBytes), fmt.Sprint(n.RXPackets), fmt.Sprint(n.TXPackets), fmt.Sprint(n.RXErrors), fmt.Sprint(n.TXErrors), fmt.Sprint(n.RXDrops), fmt.Sprint(n.TXDrops)})
	}
	return t
}
