package network

import (
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Mersad-Moghaddam/syskit/internal/platform"
)

func TestParseDev(t *testing.T) {
	interfaces, err := ParseDev([]byte("Inter-|\n face |bytes\n eth0: 10 2 3 4 0 0 0 0 20 5 6 7 0 0 0 0\n"))
	assert.NoError(t, err)
	assert.Len(t, interfaces, 1)
	assert.Equal(t, "eth0", interfaces[0].Name)
	assert.Equal(t, uint64(20), interfaces[0].TXBytes)
}

func TestCollectorEnrichesInterfacesFromSysfs(t *testing.T) {
	fsys := platform.TestFS(fstest.MapFS{
		"proc/net/dev":                 &fstest.MapFile{Data: []byte("Inter-|\n face |bytes\n eth0: 10 2 3 4 0 0 0 0 20 5 6 7 0 0 0 0\n")},
		"sys/class/net/eth0/operstate": &fstest.MapFile{Data: []byte("up\n")},
		"sys/class/net/eth0/mtu":       &fstest.MapFile{Data: []byte("1500\n")},
		"sys/class/net/eth0/address":   &fstest.MapFile{Data: []byte("02:00:00:00:00:01\n")},
	})

	info, err := NewCollector(fsys).Collect()
	require.NoError(t, err)
	require.Len(t, info.Interfaces, 1)
	assert.Equal(t, "up", info.Interfaces[0].State)
	require.NotNil(t, info.Interfaces[0].MTU)
	assert.Equal(t, uint32(1500), *info.Interfaces[0].MTU)
	assert.Equal(t, "02:00:00:00:00:01", info.Interfaces[0].MACAddress)
}

func TestCollectorAddsNetlinkAddresses(t *testing.T) {
	fsys := platform.TestFS(fstest.MapFS{
		"proc/net/dev": &fstest.MapFile{Data: []byte("Inter-|\n face |bytes\n eth0: 10 2 3 4 0 0 0 0 20 5 6 7 0 0 0 0\n")},
	})
	info, err := NewCollectorWithAddresses(fsys, addressSourceStub{addresses: []platform.InterfaceAddress{{Interface: "eth0", Address: "192.0.2.10/24"}}}).Collect()
	require.NoError(t, err)
	assert.Equal(t, []string{"192.0.2.10/24"}, info.Interfaces[0].Addresses)
}

type addressSourceStub struct {
	addresses []platform.InterfaceAddress
	err       error
}

func (s addressSourceStub) InterfaceAddresses() ([]platform.InterfaceAddress, error) {
	return s.addresses, s.err
}

func TestRoutesAndResolvers(t *testing.T) {
	routes, err := ParseRoutes([]byte("Iface Destination Gateway Flags\neth0 00000000 0100A8C0 0003\n"))
	assert.NoError(t, err)
	assert.True(t, routes[0].Default)
	assert.Equal(t, "192.168.0.1", routes[0].Gateway)
	assert.Equal(t, []string{"1.1.1.1", "8.8.8.8"}, ParseResolvConf([]byte("# comment\nnameserver 1.1.1.1\nnameserver 8.8.8.8 # public\n")))
}
