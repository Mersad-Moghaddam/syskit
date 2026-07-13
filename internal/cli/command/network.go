package command

import (
	"fmt"
	"strings"
	"time"

	"github.com/Mersad-Moghaddam/syskit/internal/model"
	"github.com/Mersad-Moghaddam/syskit/internal/render"
	"github.com/spf13/cobra"
)

type NetworkService interface {
	Collect() (*model.NetworkInfo, error)
	Sample(time.Duration) (*model.NetworkInfo, error)
}
type NetworkOptions struct {
	Format   func() string
	NoHeader func() bool
}

func NewNetworkCmd(s NetworkService, o NetworkOptions) *cobra.Command {
	var interval time.Duration
	cmd := &cobra.Command{Use: "network", Short: "Show network interface counters", Args: cobra.NoArgs, RunE: func(c *cobra.Command, args []string) error {
		var info *model.NetworkInfo
		var err error
		if interval > 0 {
			info, err = s.Sample(interval)
		} else {
			info, err = s.Collect()
		}
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
	cmd.Flags().DurationVar(&interval, "interval", 0, "sample interface bandwidth over this interval")
	cmd.AddCommand(newNetworkInterfacesCmd(s, o), newNetworkRoutesCmd(s, o), newNetworkDNSCmd(s, o))
	return cmd
}

func newNetworkInterfacesCmd(s NetworkService, o NetworkOptions) *cobra.Command {
	return &cobra.Command{Use: "interfaces", Short: "Show network interface counters", Args: cobra.NoArgs, RunE: func(c *cobra.Command, args []string) error {
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
		return r.Render(c.OutOrStdout(), info.Interfaces)
	}}
}

func newNetworkRoutesCmd(s NetworkService, o NetworkOptions) *cobra.Command {
	return &cobra.Command{Use: "routes", Short: "Show IPv4 routes and default gateway", Args: cobra.NoArgs, RunE: func(c *cobra.Command, args []string) error {
		info, err := s.Collect()
		if err != nil {
			return fmt.Errorf("collecting network information: %w", err)
		}
		r, err := render.New(o.Format(), render.WithNoHeader(o.NoHeader()))
		if err != nil {
			return err
		}
		if o.Format() == "table" {
			return r.Render(c.OutOrStdout(), routeTable(info.Routes))
		}
		return r.Render(c.OutOrStdout(), info.Routes)
	}}
}

func newNetworkDNSCmd(s NetworkService, o NetworkOptions) *cobra.Command {
	return &cobra.Command{Use: "dns", Short: "Show configured DNS nameservers", Args: cobra.NoArgs, RunE: func(c *cobra.Command, args []string) error {
		info, err := s.Collect()
		if err != nil {
			return fmt.Errorf("collecting network information: %w", err)
		}
		r, err := render.New(o.Format(), render.WithNoHeader(o.NoHeader()))
		if err != nil {
			return err
		}
		if o.Format() == "table" {
			return r.Render(c.OutOrStdout(), dnsTable(info.Nameservers))
		}
		return r.Render(c.OutOrStdout(), struct {
			Nameservers []string `json:"nameservers"`
		}{Nameservers: info.Nameservers})
	}}
}
func networkTable(info *model.NetworkInfo) render.Table {
	t := render.Table{Headers: []string{"IFACE", "STATE", "MTU", "MAC", "ADDRESSES", "RX BYTES", "TX BYTES", "RX B/s", "TX B/s", "RX PACKETS", "TX PACKETS", "RX ERR", "TX ERR", "RX DROP", "TX DROP"}}
	for _, n := range info.Interfaces {
		rx, tx := "unavailable", "unavailable"
		mtu := "unavailable"
		if n.MTU != nil {
			mtu = fmt.Sprint(*n.MTU)
		}
		if n.RXBytesPerSecond != nil {
			rx = fmt.Sprintf("%.0f", *n.RXBytesPerSecond)
		}
		if n.TXBytesPerSecond != nil {
			tx = fmt.Sprintf("%.0f", *n.TXBytesPerSecond)
		}
		t.Rows = append(t.Rows, []string{n.Name, n.State, mtu, n.MACAddress, strings.Join(n.Addresses, ","), fmt.Sprint(n.RXBytes), fmt.Sprint(n.TXBytes), rx, tx, fmt.Sprint(n.RXPackets), fmt.Sprint(n.TXPackets), fmt.Sprint(n.RXErrors), fmt.Sprint(n.TXErrors), fmt.Sprint(n.RXDrops), fmt.Sprint(n.TXDrops)})
	}
	return t
}

func routeTable(routes []model.Route) render.Table {
	t := render.Table{Headers: []string{"IFACE", "DESTINATION", "GATEWAY", "DEFAULT"}}
	for _, route := range routes {
		t.Rows = append(t.Rows, []string{route.Interface, route.Destination, route.Gateway, fmt.Sprint(route.Default)})
	}
	return t
}

func dnsTable(nameservers []string) render.Table {
	t := render.Table{Headers: []string{"NAMESERVER"}}
	for _, nameserver := range nameservers {
		t.Rows = append(t.Rows, []string{nameserver})
	}
	return t
}
