package command

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Mersad-Moghaddam/syskit/internal/model"
)

func TestRouteAndDNSTables(t *testing.T) {
	routes := routeTable([]model.Route{{Interface: "eth0", Destination: "0.0.0.0", Gateway: "192.0.2.1", Default: true}})
	assert.Equal(t, []string{"IFACE", "DESTINATION", "GATEWAY", "DEFAULT"}, routes.Headers)
	assert.Equal(t, []string{"eth0", "0.0.0.0", "192.0.2.1", "true"}, routes.Rows[0])

	dns := dnsTable([]string{"1.1.1.1", "8.8.8.8"})
	assert.Equal(t, []string{"NAMESERVER"}, dns.Headers)
	assert.Equal(t, [][]string{{"1.1.1.1"}, {"8.8.8.8"}}, dns.Rows)
}
