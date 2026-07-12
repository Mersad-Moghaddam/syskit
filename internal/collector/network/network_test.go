package network

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseDev(t *testing.T) {
	interfaces, err := ParseDev([]byte("Inter-|\n face |bytes\n eth0: 10 2 3 4 0 0 0 0 20 5 6 7 0 0 0 0\n"))
	assert.NoError(t, err)
	assert.Len(t, interfaces, 1)
	assert.Equal(t, "eth0", interfaces[0].Name)
	assert.Equal(t, uint64(20), interfaces[0].TXBytes)
}

func TestRoutesAndResolvers(t *testing.T) {
	routes, err := ParseRoutes([]byte("Iface Destination Gateway Flags\neth0 00000000 0100A8C0 0003\n"))
	assert.NoError(t, err)
	assert.True(t, routes[0].Default)
	assert.Equal(t, "192.168.0.1", routes[0].Gateway)
	assert.Equal(t, []string{"1.1.1.1", "8.8.8.8"}, ParseResolvConf([]byte("# comment\nnameserver 1.1.1.1\nnameserver 8.8.8.8 # public\n")))
}
