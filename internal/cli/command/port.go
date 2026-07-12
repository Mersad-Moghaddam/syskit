package command

import (
	"fmt"

	"github.com/Mersad-Moghaddam/syskit/internal/model"
	"github.com/Mersad-Moghaddam/syskit/internal/render"
	"github.com/Mersad-Moghaddam/syskit/internal/service"
	"github.com/spf13/cobra"
)

type PortService interface {
	List(service.PortOptions) (*model.PortInfo, error)
}
type PortOptions struct {
	Format   func() string
	NoHeader func() bool
}

func NewPortCmd(s PortService, o PortOptions) *cobra.Command {
	var listening bool
	var protocol string
	var port int
	cmd := &cobra.Command{Use: "ports", Short: "Show TCP and UDP sockets", Args: cobra.NoArgs, RunE: func(c *cobra.Command, args []string) error {
		info, err := s.List(service.PortOptions{Listening: listening, Protocol: protocol, Port: port})
		if err != nil {
			return fmt.Errorf("collecting ports: %w", err)
		}
		r, err := render.New(o.Format(), render.WithNoHeader(o.NoHeader()))
		if err != nil {
			return err
		}
		if o.Format() == "table" {
			return r.Render(c.OutOrStdout(), portTable(info))
		}
		return r.Render(c.OutOrStdout(), info)
	}}
	cmd.Flags().BoolVar(&listening, "listening", false, "show listening sockets only")
	cmd.Flags().StringVar(&protocol, "protocol", "", "filter by tcp or udp")
	cmd.Flags().IntVar(&port, "port", 0, "filter by local port")
	return cmd
}
func portTable(info *model.PortInfo) render.Table {
	t := render.Table{Headers: []string{"PROTO", "LOCAL", "REMOTE", "STATE", "INODE"}}
	for _, s := range info.Sockets {
		t.Rows = append(t.Rows, []string{s.Protocol, fmt.Sprintf("%s:%d", s.LocalAddress, s.LocalPort), fmt.Sprintf("%s:%d", s.RemoteAddress, s.RemotePort), s.State, fmt.Sprint(s.Inode)})
	}
	return t
}
